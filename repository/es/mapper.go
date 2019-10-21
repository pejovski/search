package es

import (
	"github.com/pejovski/search/model"
	"github.com/pejovski/search/pkg/scope"
)

type Mapper interface {
	mapHitToProduct(h Hit) *model.Product
	mapProductToDocument(p *model.Product) Document
	mapScopeToQuery(s *scope.Scope) map[string]interface{}
}

type mapper struct {
}

func newMapper() Mapper {
	return mapper{}
}

func (m mapper) mapHitToProduct(h Hit) *model.Product {
	s := h.Source
	return &model.Product{Id: h.Id, Title: s.Title, Brand: s.Brand, Price: s.Price, Stock: s.Stock}
}

func (m mapper) mapProductToDocument(p *model.Product) Document {
	return Document{Title: p.Title, Brand: p.Brand, Price: p.Price, Stock: p.Stock}
}

func (m mapper) mapScopeToQuery(s *scope.Scope) map[string]interface{} {
	must := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":  s.SearchQuery,
			"fields": scope.SearchFields,
		},
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": mapScopeFilterToFilter(s),
			},
		},
	}

	return query
}

func mapScopeFilterToFilter(s *scope.Scope) []interface{} {
	var price map[string]interface{}
	var filter []interface{}

	for _, f := range s.Filters {
		switch f.Key {
		case scope.FilterKeyBrand:
			filter = append(filter, map[string]interface{}{
				"match": map[string]interface{}{
					f.Key: map[string]interface{}{
						"query": f.Value,
					},
				},
			},
			)
		case scope.FilterKeyPrice:
			if price == nil {
				price = map[string]interface{}{
					f.Operator: f.Value,
				}
				continue
			}
			price[f.Operator] = f.Value
		}
	}

	if price != nil {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": price,
			},
		})
	}
	return filter
}
