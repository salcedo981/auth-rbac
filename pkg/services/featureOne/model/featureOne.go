package mdlFeatureOne

import "github.com/google/uuid"

type (
	SampleModel struct {
		Id        int    `json:"id"`
		Code      string `json:"code"`
		Name      string `json:"name"`
		EncodedBy string `json:"encodedBy"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	Items struct {
		ID		uuid.UUID `json:"id"`
		Image	string    `json:"link"`
		ProductName string`json:"product_name"`
		Category string 	`json:"category"`
		Price float64		`json:"price"`
		Quantity int 		`json:"quantity"`

	}

	ItemBody struct {
		Link	string    `json:"link"`
		ProductName string`json:"product_name"`
		Category string 	`json:"category"`
		Price float64		`json:"price"`
		Quantity int 		`json:"quantity"`
	}


	CategoryParams struct{
		Category string	`json:"category"`
	}


)

