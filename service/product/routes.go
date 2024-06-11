package product

import (
	"ecom/internal/database"
	"ecom/models"
	"ecom/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	q *database.Queries
}

func NewHandler(q *database.Queries) *Handler {
	return &Handler{
		q: q,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/products", h.HandleProducts)
	router.Post("/admin/products", h.HandleAddProduct)
	router.Put("/admin/products", h.HandleUpdateProduct)
	router.Delete("/admin/products", h.HandleDeleteProduct)
}

func (h *Handler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	prods, err := h.q.GetProducts(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintln("Couldn't fetch products"))
		return
	}
	if len(prods) == 0 {
		utils.RespondWithJSON(w, http.StatusOK, fmt.Sprintln("Currently No products for sale"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseProductsToProducts(prods))
}

func (h *Handler) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameters struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Price       int32  `json:"price"`
		Quantity    int32  `json:"quantity"`
	}

	params := Parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}

	_, err = h.q.CreateProduct(r.Context(), database.CreateProductParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		Description: params.Description,
		Image:       params.Image,
		Price:       params.Price,
		Quantity:    params.Quantity,
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Error Adding item to products"))
		return
	}

	utils.RespondWithJSON(w, 200, fmt.Sprintln("Product added successfully"))
}

func (h *Handler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameters struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Image       string    `json:"image"`
		Price       int32     `json:"price"`
		Quantity    int32     `json:"quantity"`
	}

	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}

	prod, err := h.q.UpdateProduct(r.Context(), database.UpdateProductParams{
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		Description: params.Description,
		Image:       params.Image,
		Price:       params.Price,
		Quantity:    params.Quantity,
		ID:          params.ID,
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Error Updating product: %v", params))
		return
	}

	utils.RespondWithJSON(w, 200, prod)
}

func (h *Handler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameter struct {
		ID uuid.UUID `json:"id"`
	}
	param := Parameter{}

	err := decoder.Decode(&param)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}

	prod, err := h.q.DeleteProduct(r.Context(), param.ID)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Error deleting prod with prodID: %v", param.ID))
		return
	}

	utils.RespondWithJSON(w, 200, fmt.Sprintf("Product %v deleted", prod))
}
