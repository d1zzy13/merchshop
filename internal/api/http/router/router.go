package router

import (
	"net/http"

	"merchshop/internal/api/http/auth"
	"merchshop/internal/api/http/handlers"
	"merchshop/internal/api/http/middleware"

	_ "merchshop/cmd/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func NewRouter(h *handlers.Handler, tokenManager auth.TokenManager) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/auth", h.Auth).Methods(http.MethodPost)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(tokenManager))

	api.HandleFunc("/info", h.Info).Methods(http.MethodGet)
	api.HandleFunc("/sendCoin", h.SendCoin).Methods(http.MethodPost)
	api.HandleFunc("/buy/{item}", h.Buy).Methods(http.MethodGet)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
