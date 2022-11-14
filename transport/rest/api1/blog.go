package api1

import (
	"net/http"

	"github.com/aryanugroho/blogapp/app"
	"github.com/aryanugroho/blogapp/internal/helpers/errors"
	"github.com/aryanugroho/blogapp/internal/helpers/validation"
	"github.com/aryanugroho/blogapp/internal/httpapi/server"
)

func (api *API) InitAPI() {
	api.BaseRoutes.Post = api.BaseRoutes.ApiRoot.PathPrefix("/post").Subrouter()
	api.BaseRoutes.Post.Handle("/", http.HandlerFunc(api.createPost)).Methods("POST")
}

func (api *API) createPost(w http.ResponseWriter, r *http.Request) {
	dto := app.PostPayload{}
	err := server.ReadJSON(r, &dto)
	if err != nil {
		server.WriteResponseError(w, r, errors.BadRequest.New(err.Error()))
		return
	}

	err = validation.Validate.Struct(dto)
	if err != nil {
		server.WriteResponseError(w, r, errors.BadRequest.New(err.Error()))
		return
	}

	result, err := api.app.CreatePost(r.Context(), &dto)
	if err != nil {
		server.WriteResponseError(w, r, err)
		return
	}

	server.WriteJSON(w, r, http.StatusOK, result)
}
