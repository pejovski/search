package repository

import (
	"testing"
)

func TestMapHitToProduct(t *testing.T) {
	hit := &Hit{
		Id: "111",
		Source: Document{
			Title: "Galaxy",
			Brand: "Samsung",
			Price: 800,
			Stock: 3,
		},
	}

	p := mapHitToProduct(hit)

	if p.Id != hit.Id {
		t.Error("Expected ids to be equal")
	}
}
