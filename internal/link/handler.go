package link

import (
	"fmt"
	"link-manager/configs"
	"link-manager/pkg/middleware"
	"link-manager/pkg/request"
	"link-manager/pkg/response"
	"log"
	"strconv"

	//"link-manager/pkg/token"
	"net/http"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepo *LinkRepository
	Config   *configs.Config
}

type LinkHandler struct {
	LinkRepo *LinkRepository
}

func (handler *LinkHandler) Create() http.HandlerFunc {
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

func (handler *LinkHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		hash := r.PathValue("hash")
		link, err := handler.LinkRepo.GetByHash(hash)
		if err != nil {
			response.Json(w, err, http.StatusNotFound) //404
			return
		}
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect) //307
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
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

func (handler *LinkHandler) Delete() http.HandlerFunc {
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

func (handler *LinkHandler) GetList() http.HandlerFunc {
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

		resp := GetAllResponse{
			Links: handler.LinkRepo.GetAll(page, pageSize),
			Count: handler.LinkRepo.GetCount(),
		}

		response.Json(w, resp, http.StatusOK) // 200
	}
}

func NewLinkHendler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepo: deps.LinkRepo,
	}

	router.Handle("POST /link", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.Read())
	router.Handle("GET /link", middleware.IsAuthed(handler.GetList(), deps.Config))
}
