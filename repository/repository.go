package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pejovski/search/entity"
	"github.com/pejovski/search/scope"
	"github.com/sirupsen/logrus"
	"net/http"
)

const index = "products"

type ESProductRepository struct {
	client *elasticsearch.Client
}

func NewESProductRepository(es *elasticsearch.Client) *ESProductRepository {
	return &ESProductRepository{client: es}
}

func (r ESProductRepository) Product(id string) (*entity.Product, error) {
	var h *Hit

	res, err := r.client.Get(index, id)
	if err != nil {
		logrus.Errorf("Failed to get product %s", id)
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return nil, fmt.Errorf("response error")
	}

	if err := json.NewDecoder(res.Body).Decode(&h); err != nil {
		logrus.Errorf("Failed to decode body for product %s", id)
		return nil, err
	}
	defer res.Body.Close()

	return mapHitToProduct(h), nil
}

func (r ESProductRepository) Create(p *entity.Product) (string, error) {
	d := mapProductToDocument(p)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(d); err != nil {
		logrus.Errorf("Failed to decode body for product %s", p.Id)
		return "", err
	}

	res, err := r.client.Create(index, p.Id, &buf)
	if err != nil {
		logrus.Errorf("Failed to create product %s", p.Id)
		return "", err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", p.Id, res.StatusCode, res.String())
		return "", fmt.Errorf("response error")
	}

	return p.Id, nil
}

func (r ESProductRepository) Delete(id string) error {
	res, err := r.client.Delete(index, id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s", id)
		return err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return fmt.Errorf("response error")
	}

	return nil
}

func (r ESProductRepository) Products(s *scope.Scope) ([]*entity.Product, int, error) {

	query := mapScopeToQuery(s)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logrus.Errorf("Failed to encode query %v", query)
		return nil, 0, err
	}

	// Perform the search request.
	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
		r.client.Search.WithSize(s.Pagination.Limit),
		r.client.Search.WithFrom(s.Pagination.Offset),
		r.client.Search.WithSort(fmt.Sprintf("%s:%s", s.Sorting.Field, s.Sorting.Order)),
	)
	if err != nil {
		logrus.Errorf("Failed to get response for query")
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logrus.Errorf("Error in the response. Status code: %d. Response: %s", res.StatusCode, res.String())
		return nil, 0, fmt.Errorf("response error")
	}

	var result *Result

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.Errorf("Failed to decode result")
		return nil, 0, err
	}

	var products []*entity.Product

	for _, hit := range result.Hits.Hits {
		products = append(products, mapHitToProduct(&hit))
	}

	return products, result.Hits.Total.Value, nil
}
