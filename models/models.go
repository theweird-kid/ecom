package models

import (
	"ecom/internal/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       int32     `json:"price"`
	Quantity    int32     `json:"quantity"`
}

type Cart struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ProdID     uuid.UUID `json:"prod_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	TotalPrice int32     `json:"total_price"`
	Quantity   int32     `json:"quantity"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
	}
}

func DatabaseProductToProduct(dbProduct database.Product) Product {
	return Product{
		ID:          dbProduct.ID,
		CreatedAt:   dbProduct.CreatedAt,
		UpdatedAt:   dbProduct.UpdatedAt,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Image:       dbProduct.Image,
		Price:       dbProduct.Price,
		Quantity:    dbProduct.Quantity,
	}
}

func DatabaseProductsToProducts(dbProducts []database.Product) []Product {
	products := []Product{}
	for _, product := range dbProducts {
		products = append(products, DatabaseProductToProduct(product))
	}
	return products
}
