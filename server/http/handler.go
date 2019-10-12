package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pejovski/search/domain"
	"github.com/pejovski/search/entity"
	"github.com/pejovski/search/scope"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	controller domain.SearchController
}

func NewHandler(c domain.SearchController) *Handler {
	return &Handler{
		controller: c,
	}
}

func (h Handler) Products() http.HandlerFunc {

	type Result struct {
		Offset int        `json:"offset"`
		Limit  int        `json:"limit"`
		Count  int        `json:"count"`
		Total  int        `json:"total"`
		Items  []*Product `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		s := scope.New(
			scope.NewSearch(r),
			scope.NewPagination(r),
			scope.NewFilters(r),
			scope.NewSorting(r),
		)

		if s.SearchQuery == "" {
			logrus.Warnln("Search query empty")
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		dps, total, err := h.controller.GetProducts(s)
		if err != nil {
			logrus.Errorf("Failed to get products. Error: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		res := Result{
			Offset: s.Pagination.Offset,
			Limit:  s.Pagination.Limit,
			Count:  len(dps),
			Total:  total,
			Items:  mapDomainProductsToProducts(dps),
		}

		h.respond(w, r, res, http.StatusOK)
	}
}

func (h Handler) Product() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		p, err := h.controller.GetProduct(id)
		if err != nil {
			logrus.Errorf("Failed to get product with id %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if p == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		h.respond(w, r, mapDomainProductToProduct(p), http.StatusOK)
	}
}

func (h Handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p *entity.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			logrus.Warnln("Failed to decode request body")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		id, err := h.controller.CreateProduct(p)
		if err != nil {
			logrus.Errorf("Failed to create product for id %s. Error: %s", p.Id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/products/%s", id))
		h.respond(w, r, nil, http.StatusCreated)
	}
}

func (h Handler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		if err := h.controller.DeleteProduct(id); err != nil {
			logrus.Errorf("Failed to delete product %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusNoContent)
	}
}

func (h Handler) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.Errorf("Failed to encode data. Error: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
