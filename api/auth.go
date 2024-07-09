package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/GermanPachec0/app-go/domain"
	jwt "github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	logReq := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		s.errorResponse(w, r, 500, err)
		return
	}
	user, err := s.userRepo.GetByEmail(ctx, logReq.Email)
	if err != nil {
		s.errorResponse(w, r, 404, err)
		return
	}
	if !user.ValidatePassword(logReq.Password) {
		s.errorResponse(w, r, 401, err)
		return
	}

	token, err := s.createJWT(&user)
	if err != nil {
		s.errorResponse(w, r, 500, err)
	}

	res := LoginResponse{
		Email: user.Email,
		Type:  user.Type,
		Token: token,
	}

	WriteJson(w, 200, res)
}

func (s *APIServer) withJwtAuth(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("missing Authorization header"))
			return
		}

		token, err := s.validateJwt(tokenString)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid token claims type"))
			return
		}

		emailClaim, ok := claims["email"].(string)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid or missing email claim"))
			return
		}

		_, err = s.userRepo.GetByEmail(ctx, emailClaim)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		handleFunc(w, r)
	}
}

func (s *APIServer) isAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("missing Authorization header"))
			return
		}

		token, err := s.validateJwt(tokenString)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid token claims type"))
			return
		}
		userType, ok := claims["role"].(string)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("invalid or missing user type"))
			return
		}

		if userType != "admin" {
			s.errorResponse(w, r, http.StatusUnauthorized, errors.New("Invalid user"))
			return
		}

		next(w, r)
	}
}

func (s *APIServer) createJWT(user *domain.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":   15000,
		"memeber_uid": user.Uuid,
		"email":       user.Email,
		"role":        user.Type,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

func (s *APIServer) validateJwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}
