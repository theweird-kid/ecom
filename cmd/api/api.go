package api

import (
	"database/sql"
	"log"
	"net/http"

	"ecom/service/cart"
	"ecom/service/order"
	"ecom/service/product"
	"ecom/service/user"

	"ecom/internal/database"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
	q    *database.Queries
}

func NewAPIServer(addr string, db *sql.DB, q *database.Queries) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		q:    q,
	}
}

func (s *APIServer) Run() error {

	//router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)

	//subrouter
	subrouter := chi.NewRouter()
	router.Mount("/api/v1", subrouter)

	//user service
	userHandler := user.NewHandler(s.q)
	userHandler.RegisterRoutes(subrouter)

	//product service
	productHandler := product.NewHandler(s.q)
	productHandler.RegisterRoutes(subrouter)

	//cart service
	cartHandler := cart.NewHandler(s.q)
	cartHandler.RegisterRoutes(subrouter)

	//order service
	orderHandler := order.NewHandler(s.db, s.q)
	orderHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
