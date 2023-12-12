package api

import "parsingWB/internal/models"

type storage interface {
	InsertData(product models.Product) error
}
