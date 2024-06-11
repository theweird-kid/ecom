package cart

import (
	"ecom/internal/database"
	"ecom/service/auth"
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
	router.Get("/cart", auth.MiddlewareAuth(h.HandleGetCart, h.q))
	router.Post("/cart", auth.MiddlewareAuth(h.HandleAddToCart, h.q))
	router.Put("/cart", auth.MiddlewareAuth(h.HandleUpdateCart, h.q))
	router.Delete("/cart", auth.MiddlewareAuth(h.HandleDeleteCart, h.q))
}

func (h *Handler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	cartItems, err := h.q.GetCart(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Error Fetching items from Cart"))
		return
	}

	utils.RespondWithJSON(w, 200, cartItems)

}

func (h *Handler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameters struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int32     `json:"quantity"`
	}

	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}

	prod, err := h.q.GetProductByID(r.Context(), params.ProductID)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Product Doesn't exist"))
		return
	}

	if prod.Quantity < params.Quantity {
		utils.RespondWithJSON(w, 200, fmt.Sprintf("Only %d units of Product %v available. ", prod.Quantity, prod.ID))
		return
	}

	_, err = h.q.AddToCart(r.Context(), database.AddToCartParams{
		ID:         uuid.New(),
		UserID:     auth.GetUserIDFromContext(r.Context()),
		ProdID:     params.ProductID,
		UpdatedAt:  time.Now().UTC(),
		TotalPrice: prod.Price * params.Quantity,
		Quantity:   params.Quantity,
	})

	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("The Product already exists in cart"))
		return
	}

	utils.RespondWithJSON(w, 200, fmt.Sprintf("product %v successfully added to cart of %v", prod.ID, auth.GetUserIDFromContext(r.Context())))

}

func (h *Handler) HandleUpdateCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameters struct {
		CartItemID uuid.UUID `json:"cart_item_id"`
		ProductID  uuid.UUID `json:"product_id"`
		Quantity   int32     `json:"quantity"`
	}

	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}
	prod, err := h.q.GetProductByID(r.Context(), params.ProductID)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Invalid Item"))
		return
	}

	if prod.Quantity < params.Quantity {
		utils.RespondWithJSON(w, 200, fmt.Sprintf("Only %d units of Product %v available. ", prod.Quantity, prod.ID))
		return
	}

	newItem, err := h.q.UpdateCart(r.Context(), database.UpdateCartParams{
		Quantity:   params.Quantity,
		UpdatedAt:  time.Now().UTC(),
		ID:         params.CartItemID,
		TotalPrice: params.Quantity * prod.Price,
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Unable to update cart"))
		return
	}

	utils.RespondWithJSON(w, 200, fmt.Sprintf("Item updated to: %v", newItem))
}

func (h *Handler) HandleDeleteCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type Parameters struct {
		CartItemID uuid.UUID `json:"cart_item_id"`
	}

	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Error parsing JSON"))
		return
	}

	err = h.q.DeleteCartItem(r.Context(), params.CartItemID)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintln("Failed to Delete Item"))
		return
	}

	utils.RespondWithJSON(w, 200, fmt.Sprintln("Item Deleted Successfully from cart"))
}

/*
user_id: 08dec99a-3843-4443-9067-2ee466be560a
prod_id: e83aeb52-d9df-4a0a-80f6-7e8df9acf1f0
q:4
c:1000
*/
