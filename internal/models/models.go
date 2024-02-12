package models

import "time"

type Sizes struct {
	ChrtID   int64    `json:"chrtID"`
	TechSize int      `json:"techSize"`
	Skus     []string `json:"skus"`
}

type Cards struct {
	Sizes []Sizes `json:"sizes"`
	// MediaFiles []string
	UpdateAt   string `json:"updateAt"`
	VendorCode string `json:"vendorCode"`
	Brand      string
	Object     string `json:"object"`
	NmID       int64  `json:"nmID"`
	ImtID      int64  `json:"imtID"`
	Title      string `json:"title"`
}

type Cursor struct {
	UpdatedAt time.Time `json:"updatedAt"`
	NmID      int64     `json:"nmID"`
	Total     int64     `json:"total"`
}

type Data struct {
	Cards  []Cards `json:"cards"`
	Cursor Cursor  `json:"cursor"`
}

type Out struct {
	Data  Data
	Error bool
}

type SellerPrice struct {
	Saller string
	Price  float64
}

type Product struct {
	NmID        int64
	Name        string
	Price       float64
	Barcode     string
	SellerPrice []SellerPrice
}
