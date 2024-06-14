package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ecom/config"
	"ecom/internal/database"
	"ecom/models"
	"ecom/service/auth"
	"ecom/utils"

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
	router.Post("/login", h.HandleLogin)
	router.Get("/login", auth.MiddlewareAuth(h.HandleTest, h.q))
	router.Post("/register", h.HandleRegister)
}

func (h *Handler) HandleTest(w http.ResponseWriter, r *http.Request) {
	id := auth.GetUserIDFromContext(r.Context())
	user, err := h.q.GetUserByID(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintln("Internal server problem"))
		return
	}
	utils.RespondWithJSON(w, 200, models.DatabaseUserToUser(user))
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Wrong Credentials: %v", err))
		return
	}

	//get user by email
	u, err := h.q.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintln("Invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(params.Password)) {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintln("Invalid email or password"))
		return
	}

	token, err := auth.CreateJWT(config.Envs.Secret, u.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintln("Error Generating authentication"))
		return
	}

	//create token cookie
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		MaxAge:   3600 * 4,
		Path:     "/",  // Make cookie accessible to all paths
		HttpOnly: true, // Recommended to help prevent XSS
		Secure:   true,
		SameSite: http.SameSiteNoneMode, // Use SameSiteNoneMode if your frontend and backend are served from different origins
	}

	http.SetCookie(w, cookie)
	utils.RespondWithJSON(w, http.StatusOK, u)

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 300, fmt.Sprintf("Error Parsing JSON: %v", err))
		return
	}

	//check if a user with the given email already exists
	u, _ := h.q.GetUserByEmail(r.Context(), params.Email)
	fmt.Println(u)
	if u.Email == params.Email {
		utils.RespondWithError(w, 500, fmt.Sprintf("User with email %v already exists", params.Email))
		return
	}
	pass, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("%v", err))
	}

	user, err := h.q.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  pass,
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Couldn't create user: %v", err))

		return
	}
	utils.RespondWithJSON(w, 200, models.DatabaseUserToUser(user))
}
