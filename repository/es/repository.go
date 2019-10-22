package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/search/model"
	"github.com/pejovski/search/pkg/scope"
	repo "github.com/pejovski/search/repository"
)

const index = "products"

type repository struct {
	client *elasticsearch.Client
	mapper Mapper
}

func NewRepository(es *elasticsearch.Client) repo.Repository {
	return repository{
		client: es,
		mapper: newMapper(),
	}
}

func (r repository) Product(id string) (*model.Product, error) {
	var h Hit

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

	return r.mapper.mapHitToProduct(h), nil
}

func (r repository) Create(p *model.Product) (string, error) {
	d := r.mapper.mapProductToDocument(p)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(d); err != nil {
		logrus.Errorf("Failed to decode body for product %s", p.ID)
		return "", err
	}

	res, err := r.client.Create(index, p.ID, &buf)
	if err != nil {
		logrus.Errorf("Failed to create product %s", p.ID)
		return "", err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", p.ID, res.StatusCode, res.String())
		return "", fmt.Errorf("response error")
	}

	return p.ID, nil
}

func (r repository) Delete(id string) error {
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

func (r repository) Products(s *scope.Scope) ([]*model.Product, int, error) {

	query := r.mapper.mapScopeToQuery(s)

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
		if res.StatusCode == http.StatusNotFound {
			return nil, 0, nil
		}
		logrus.Errorf("Error in the response. Status code: %d. Response: %s", res.StatusCode, res.String())
		return nil, 0, fmt.Errorf("response error")
	}

	var result *Result

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.Errorf("Failed to decode result")
		return nil, 0, err
	}

	var products []*model.Product

	for _, hit := range result.Hits.Hits {
		products = append(products, r.mapper.mapHitToProduct(hit))
	}

	return products, result.Hits.Total.Value, nil
}
