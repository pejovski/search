package es

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/pejovski/search/model"
)

type MapperSuite struct {
	suite.Suite
}

func TestMapperSuite(t *testing.T) {
	suite.Run(t, new(MapperSuite))
}

func (s *MapperSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()
}

func (s *MapperSuite) TestMapHitToProduct() {
	h := Hit{
		ID: "1",
		Source: Document{
			Title: "galaxy",
			Brand: "samsung",
			Price: 50,
			Stock: 3,
		},
	}
	m := newMapper()
	p := m.mapHitToProduct(h)

	s.NotEmpty(p)
	s.Equal(h.ID, p.ID)
	s.Equal(h.Source.Title, p.Title)
	s.Equal(h.Source.Stock, p.Stock)
	s.Equal(h.Source.Brand, p.Brand)
	s.Equal(h.Source.Price, p.Price)
}

func (s *MapperSuite) TestMapProductToDocument() {
	p := &model.Product{
		ID:    "1",
		Title: "galaxy",
		Brand: "samsung",
		Price: 50,
		Stock: 3,
	}

	m := newMapper()
	d := m.mapProductToDocument(p)

	s.NotEmpty(d)
	s.Equal(p.Title, d.Title)
	s.Equal(p.Stock, d.Stock)
	s.Equal(p.Brand, d.Brand)
	s.Equal(p.Price, d.Price)
}
