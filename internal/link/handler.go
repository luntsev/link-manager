package link

import (
	"fmt"
	"link-manager/configs"
	"link-manager/pkg/event"
	"link-manager/pkg/middleware"
	"link-manager/pkg/request"
	"link-manager/pkg/response"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepo *LinkRepository
	EventBus *event.EventBus
	Config   *configs.Config
}

type LinkHandler struct {
	LinkRepo *LinkRepository
	EventBus *event.EventBus
}

func NewLinkHendler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepo: deps.LinkRepo,
		EventBus: deps.EventBus,
	}

	router.Handle("POST /link", middleware.IsAuthed(handler.create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.delete(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.read())
	router.Handle("GET /link", middleware.IsAuthed(handler.getList(), deps.Config))
}

func (handler *LinkHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := request.HandleBody[CreateRequest](&w, r)
		if err != nil {
			return
		}

		link := NewLink(body.Url, body.Name)
		n := 9
		for _, err := handler.LinkRepo.GetByHash(link.Hash); err == nil; n++ {
			link.GenHash(n / 3)
		}

		resp, err := handler.LinkRepo.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}
		response.Json(w, resp, http.StatusCreated) //201
	}
}

func (handler *LinkHandler) read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		hash := r.PathValue("hash")
		link, err := handler.LinkRepo.GetByHash(hash)
		if err != nil {
			response.Json(w, err, http.StatusNotFound) //404
			return
		}
		go handler.EventBus.Pubish(event.Event{
			Type:  event.EventLinkGet,
			Event: link.ID,
		})

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect) //307
	}
}

func (handler *LinkHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			log.Println("wrong email in JWT token")
			return
		}

		body, err := request.HandleBody[UpdateRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		id := r.PathValue("id")
		idRow, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		link, err := handler.LinkRepo.Update(&Link{
			Model: gorm.Model{ID: uint(idRow)},
			Name:  body.Name,
			Url:   body.Url,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		response.Json(w, link, http.StatusOK)
	}
}

func (handler *LinkHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			log.Println("wrong email in JWT token")
			return
		}

		id := r.PathValue("id")
		idRow, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		link, err := handler.LinkRepo.GetById(uint(idRow))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		err = handler.LinkRepo.Delete(uint(idRow))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //500
			return
		}
		response.Json(w, link, http.StatusOK) //200
	}
}

func (handler *LinkHandler) getList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			log.Println("wrong email in JWT token")
			return
		}
		fmt.Println(email)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil {
			pageSize = 100
		}

		links, count := handler.LinkRepo.GetAll(page, pageSize)
		resp := GetAllResponse{
			Links: links,
			Count: count,
		}

		response.Json(w, resp, http.StatusOK) // 200
	}
}
