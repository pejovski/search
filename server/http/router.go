package http

import (
	"github.com/gorilla/mux"
	_ "github.com/pejovski/search/app/statik"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router struct {
	Router  *mux.Router
	Handler *Handler
}

func NewRouter(h *Handler) *Router {
	s := &Router{
		Router:  mux.NewRouter(),
		Handler: h,
	}

	s.health()
	s.swagger()
	s.routes()

	return s
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.Router.ServeHTTP(w, r)
}

func (rtr *Router) routes() {
	routerV1 := rtr.Router.PathPrefix("/v1").Subrouter()
	routerV1.Use(auth)

	routerV1.HandleFunc("/products", rtr.Handler.Products()).Methods("GET")
	routerV1.HandleFunc("/products", rtr.Handler.CreateProduct()).Methods("POST")
	routerV1.HandleFunc("/products/{id}", rtr.Handler.Product()).Methods("GET")
	routerV1.HandleFunc("/products/{id}", rtr.Handler.DeleteProduct()).Methods("DELETE")
}

func (rtr *Router) swagger() {
	// swagger handler
	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find statik", err)
	}
	sh := http.FileServer(statikFS)

	rtr.Router.Handle("/", sh).Methods("GET")
	rtr.Router.PathPrefix("/swagger").Handler(sh)
}

func (rtr *Router) health() {
	rtr.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Up"))
	}).Methods("GET")
}
