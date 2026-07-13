package dto

import (
	"time"

	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
)

type CustomerDTO struct {
	Name  		string `json:"name"`
	Email 		string `json:"email"`
	Phone 		string `json:"phone"`
	Password 	string `json:"password"`
}

type CustomerResponseDTO struct {
	ID       			uuid.UUID `json:"id"`
	Name      		string 		`json:"name"`
	Email     		string 		`json:"email"`
	Phone     		string 		`json:"phone"`
	PasswordHash 	string 		`json:"passwordHash"`
	CreatedAt 		time.Time `json:"createdAt"`
	UpdatedAt 		time.Time `json:"updatedAt"`
}

func NewCustomerResponseDTO(customer model.Customer) CustomerResponseDTO {
	return CustomerResponseDTO{
		ID:        		customer.ID,
		Name:      		customer.Name,
		Email:     		customer.Email,
		Phone:     		customer.Phone,
		PasswordHash: customer.PasswordHash,
		CreatedAt: 		customer.CreatedAt,
		UpdatedAt: 		customer.UpdatedAt,
	}
}