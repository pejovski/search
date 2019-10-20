package controller

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/pejovski/search/gen/mock"
	"github.com/pejovski/search/model"
)

type Suite struct {
	suite.Suite
	repository *mock.MockRepository
}

func (s *Suite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.repository = mock.NewMockRepository(ctrl)
}

func TestRepositoryApiSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetProduct() {
	id := "123"

	p := &model.Product{
		Id:    id,
		Title: "Galaxy S10",
		Brand: "Samsung",
		Price: 500,
		Stock: 5,
	}

	s.repository.EXPECT().Product(gomock.Eq(id)).Return(p, nil).Times(1)

	c := New(s.repository)

	res, err := c.GetProduct(id)

	s.NotEmpty(res)
	s.Equal(nil, err)
	s.Equal(p.Title, res.Title)
	s.Equal(p.Brand, res.Brand)
	s.Equal(p.Price, res.Price)
	s.Equal(p.Stock, res.Stock)
}

func (s *Suite) TestGetProductError() {
	id := "123"

	errRet := fmt.Errorf("failed to get product")

	s.repository.EXPECT().Product(gomock.Eq(id)).Return(nil, errRet).Times(1)

	c := New(s.repository)

	res, err := c.GetProduct(id)

	s.Nil(res)
	s.Equal(errRet, err)
}
