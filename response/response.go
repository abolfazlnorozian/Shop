package response

import (
	"shop/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"body"`
}

type MixProductsResponse struct {
	Success bool `json:"success"`

	Message string `json:"message"`

	Body []entities.MixProducts `json:"body"`
}

type GetProductByOneField struct {
	Docs []struct {
		Id          string `json:"_id"`
		NotExist    bool   `json:"notExist"`
		Amazing     bool   `json:"amazing"`
		ProductType string `json:"productType"`
		Images      []struct {
			Id  string `json:"_id"`
			URL string `json:"url"`
		} `json:"images"`
		Name            string  `json:"name"`
		Price           int     `json:"price"`
		DiscountPercent float64 `json:"discountPercent"`
		Stock           int     `json:"stock"`
		Slug            string  `json:"slug"`
		Variations      []struct {
			Id              string `json:"_id"`
			DiscountPercent int    `json:"discountPercent"`
			Keys            []int  `json:"keys"`
			Price           int    `json:"price"`
			Quantity        int    `json:"quantity"`
		} `json:"variations"`
		SalesNumber int    `json:"salesNumber"`
		BannerUrl   string `json:"bannerUrl"`
	} `json:"docs"`

	// Total number of documents
	TotalDocs int `json:"totalDocs"`

	// Number of items per page
	Limit int `json:"limit"`

	// Total number of pages
	TotalPages int `json:"totalPages"`

	// Current page number
	Page int `json:"page"`

	// Counter for pagination
	PagingCounter int `json:"pagingCounter"`

	// Flag indicating if there is a previous page
	HasPrevPage bool `json:"hasPrevPage"`

	// Flag indicating if there is a next page
	HasNextPage bool `json:"hasNextPage"`

	// Previous page number (null if there is no previous page)
	PrevPage *int `json:"prevPage"`

	// Next page number (null if there is no next page)
	NextPage *int `json:"nextPage"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

type RegisterUsersResponse struct {
	Body struct {
		Username string `json:"username"`

		Password string `json:"password"`
	} `json:"body"`

	Success string `json:"success"`

	Message string `json:"message"`
}

type LoginResponse struct {
	Body struct {
		Token string `json:"token"`

		RefreshToken string `json:"refreshToken"`
	} `json:"body"`

	Success bool `json:"success"`

	Message string `json:"message"`
}
type CommentResponse struct {
	// in: body
	Body struct {
		// Product details including comments
		// required: true
		entities.Product
	} `json:"body"`

	// Success indicator
	// required: true
	Success bool `json:"success"`

	// Message
	// required: true
	Message string `json:"message"`
}

type Input struct {
	ProductId     primitive.ObjectID `json:"productId"`
	VariationsKey []int              `json:"variationsKey"`
	Quantity      int                `json:"quantity"`
	QuantityState string             `json:"quantityState"`
}
