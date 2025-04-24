package link

import (
	"link-manager/pkg/middleware"
	"link-manager/pkg/request"
	"link-manager/pkg/response"
	"strconv"

	//"link-manager/pkg/token"
	"net/http"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepo *LinkRepository
}

type LinkHandler struct {
	LinkRepo *LinkRepository
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		body, err := request.HandleBody[CreateRequest](&w, req)
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
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		hash := req.PathValue("hash")
		link, err := handler.LinkRepo.GetByHash(hash)
		if err != nil {
			response.Json(w, err, http.StatusNotFound) //404
			return
		}
		http.Redirect(w, req, link.Url, http.StatusTemporaryRedirect) //307
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		body, err := request.HandleBody[UpdateRequest](&w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
			return
		}

		id := req.PathValue("id")
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
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		id := req.PathValue("id")
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

func NewLinkHendler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepo: deps.LinkRepo,
	}

	router.Handle("POST /link", middleware.IsAuthed(handler.Create()))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update()))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.Read())
}
