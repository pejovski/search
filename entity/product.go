package entity

type Product struct {
	Id    string  `json:"id"`
	Title string  `json:"title"`
	Brand string  `json:"brand"`
	Price float32 `json:"price"`
	Stock int     `json:"stock"`
}
