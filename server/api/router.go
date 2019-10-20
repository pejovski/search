package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	_ "github.com/pejovski/search/app/statik"
	"github.com/pejovski/search/controller"
)

type Router interface {
	routes()
	swagger()
	health()
	auth(next http.Handler) http.Handler

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type router struct {
	router  *mux.Router
	handler Handler
}

func newRouter(c controller.Controller) Router {
	s := &router{
		router:  mux.NewRouter(),
		handler: newHandler(c),
	}

	s.health()
	s.swagger()
	s.routes()

	return s
}

func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.router.ServeHTTP(w, r)
}

func (rtr *router) routes() {
	routerV1 := rtr.router.PathPrefix("/v1").Subrouter()
	routerV1.Use(rtr.auth)

	routerV1.HandleFunc("/products", rtr.handler.products()).Methods("GET")
	routerV1.HandleFunc("/products", rtr.handler.createProduct()).Methods("POST")
	routerV1.HandleFunc("/products/{id}", rtr.handler.product()).Methods("GET")
	routerV1.HandleFunc("/products/{id}", rtr.handler.deleteProduct()).Methods("DELETE")
}

func (rtr *router) swagger() {
	// swagger handler
	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find statik", err)
	}
	sh := http.FileServer(statikFS)

	rtr.router.Handle("/", sh).Methods("GET")
	rtr.router.PathPrefix("/swagger").Handler(sh)
}

func (rtr *router) health() {
	rtr.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Up"))
	}).Methods("GET")
}

func (rtr *router) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))

		if apiKey != os.Getenv("APP_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
