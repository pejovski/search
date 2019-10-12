package repository

type Document struct {
	Title string  `json:"title"`
	Brand string  `json:"brand"`
	Price float32 `json:"price"`
	Stock int     `json:"stock"`
}

type Update struct {
	Doc *Document `json:"doc"`
}

type Hit struct {
	Id     string   `json:"_id"`
	Source Document `json:"_source"`
}

type Total struct {
	Value int `json:"value"`
}

type Hits struct {
	Hits  []Hit `json:"hits"`
	Total Total `json:"total"`
}

type Result struct {
	Hits Hits `json:"hits"`
}
