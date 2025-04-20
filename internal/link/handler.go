package link

import (
	"link-manager/pkg/request"
	"link-manager/pkg/response"
	//"link-manager/pkg/token"
	"net/http"
)

type LinkHandlerDeps struct {
	LinkRepo *LinkRepository
}

type LinkHandler struct {
	LinkRepo *LinkRepository
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		body, err := request.HandleBody[CreateRequest](&w, req)
		if err != nil {
			return
		}

		link := NewLink(body.Url)
		for existLink, _ := handler.LinkRepo.GetByHash(link.Hash); existLink != nil; {
			link.GenHash()
		}

		resp, err := handler.LinkRepo.Create(link)
		if err != nil {
			response.Json(w, err, 500)
		}
		response.Json(w, resp, 200)
	}
}

func (handler *LinkHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hash := req.PathValue("hash")
		link, err := handler.LinkRepo.GetByHash(hash)
		if err != nil {
			response.Json(w, err, 404)
		}
		http.Redirect(w, req, link.Url, 307)
	}
}

/*
func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		updateReq, err := request.HandleBody[UpdateRequest](&w, req)
		if err != nil {
			return
		}

		resp := CreateResponse{
			Url:  updateReq.Url,
			Hash: token.GenToken(6),
		}

		response.Json(w, resp, 200)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")

		delReq, err := request.HandleBody[DelRequest](&w, req)
		if err != nil {
			return
		}

		resp := CreateResponse{
			Url:  "",
			Hash: delReq.Hash,
		}

		response.Json(w, resp, 200)
	}
}*/

func NewLinkHendler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepo: deps.LinkRepo,
	}

	router.HandleFunc("POST /link", handler.Create())
	//router.HandleFunc("PATCH /link/{id}", handler.Update())
	//router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.Read())
}
