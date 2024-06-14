package order

import (
	"database/sql"
	"ecom/internal/database"
	"ecom/service/auth"
	"ecom/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	db *sql.DB
	q  *database.Queries
}

func NewHandler(db *sql.DB, q *database.Queries) *Handler {
	return &Handler{
		db: db,
		q:  q,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/order", auth.MiddlewareAuth(h.HandleGetOrders, h.q))   //returns a list of orders
	router.Post("/order", auth.MiddlewareAuth(h.HandleOrderItems, h.q)) //order items in cart

}

func (h *Handler) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	orderItems, err := h.q.GetOrders(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Error Fetching items from Cart"))
		return
	}

	utils.RespondWithJSON(w, 200, orderItems)
}

func (h *Handler) HandleOrderItems(w http.ResponseWriter, r *http.Request) {
	//initiate transaction
	tx, err := h.db.Begin()
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Couldn't process transaction"))
		return
	}
	defer tx.Rollback()

	userID := auth.GetUserIDFromContext(r.Context())

	qtx := h.q.WithTx(tx)
	//get cart items
	cartItems, err := qtx.GetCart(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Error fetching cart"))
		return
	}

	for _, item := range cartItems {
		itemID := item.ID

		//check the availablity and update listing
		prod, err := h.q.GetProductByID(r.Context(), item.ProdID)
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintln("Couldn't process transaction"))
			return
		}

		//check quantity
		if prod.Quantity < item.Quantity {
			utils.RespondWithError(w, 500, fmt.Sprintf("Only %d units of product %v available", prod.Quantity, prod.ID))
			return
		}

		//update product listing
		_, err = h.q.UpdateProduct(r.Context(), database.UpdateProductParams{
			UpdatedAt:   time.Now().UTC(),
			Name:        prod.Name,
			Description: prod.Description,
			Image:       prod.Image,
			Price:       prod.Price,
			Quantity:    prod.Quantity - item.Quantity,
			ID:          prod.ID,
		})
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Only %d units of product %v available", prod.Quantity, prod.ID))
			return
		}

		//delete item from cart
		err = h.q.DeleteCartItem(r.Context(), itemID)
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintln("Error Deleting item from cart"))
			return
		}

		//order item
		_, err = h.q.OrderItem(r.Context(), database.OrderItemParams{
			ID:        uuid.New(),
			UserID:    item.UserID,
			ProdID:    item.ProdID,
			CreatedAt: time.Now().UTC(),
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Error Ordering item: %v", itemID))
			return
		}

	}

	tx.Commit()
	utils.RespondWithJSON(w, 200, fmt.Sprintln("Order Placed successfully"))
}
