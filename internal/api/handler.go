package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/olekturbo/mysterious/internal/service"
)

type Handler struct {
	idService      *service.ID
	cacheService   *service.Cache
	userService    *service.User
	tokenService   *service.Token
	predictService *service.Predict
	cookieManager  *CookieManager
}

func NewHandler(
	idService *service.ID,
	cacheService *service.Cache,
	userService *service.User,
	tokenService *service.Token,
	predictService *service.Predict,
	cookieManager *CookieManager,
) *Handler {
	return &Handler{
		idService:      idService,
		cacheService:   cacheService,
		userService:    userService,
		tokenService:   tokenService,
		predictService: predictService,
		cookieManager:  cookieManager,
	}
}

// Home godoc
// @Summary Home endpoint
// @Description Checks cookie and rate limits
// @Tags general
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Failure 429 {string} string "Too many requests"
// @Router / [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	id, err := h.cookieManager.Read(r, idCookie)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.cookieManager.Write(w, idCookie, h.idService.GenerateID())
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = h.cacheService.Limit(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	_, _ = w.Write([]byte("OK"))
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// Register godoc
// @Summary Register new user
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce plain
// @Param registerRequest body RegisterRequest true "Register payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Router /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.Create(service.CreateParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte("OK"))
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Log in user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse "Auth token"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := h.userService.Match(service.MatchParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.tokenService.Generate(match.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
	})
}

// PrivateHome godoc
// @Summary Private home endpoint
// @Description Accessible only with valid JWT
// @Tags private
// @Security ApiKeyAuth
// @Success 200 {string} string "Welcome <user_id>"
// @Failure 401 {string} string "Unauthorized"
// @Router /private [get]
func (h *Handler) PrivateHome(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["sub"].(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, _ = w.Write([]byte("Welcome user " + userID))
}

type PredictRequest struct {
	Text string `json:"text" validate:"required"`
}

type PredictResponse struct {
	Result string `json:"result"`
}

// Predict godoc
// @Summary      Predict result from input text
// @Description  Accepts a JSON payload with text input and returns a prediction result
// @Tags         prediction
// @Accept       json
// @Produce      json
// @Param        request body     PredictRequest true "Input text for prediction"
// @Success      200     {object} PredictResponse
// @Failure      400     {string} string "Invalid input"
// @Failure      500     {string} string "Internal server error"
// @Router       /predict [post]
func (h *Handler) Predict(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	predict, err := h.predictService.Predict(r.Context(), req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(PredictResponse{
		Result: predict,
	})
}
