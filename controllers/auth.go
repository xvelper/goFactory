package controllers

import (
	"encoding/json"
	"fmt"
	"gitfactory/database"
	"gitfactory/server"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := database.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := database.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		IsAdmin:  req.IsAdmin,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	userDir := filepath.Join(server.DefaultConfig.ProjectRoot, user.Username)
	err = os.MkdirAll(userDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating user directory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user database.User
	result := database.DB.Where("username = ?", creds.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !database.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "No token provided"})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome %s!", claims.Username)})
}

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var dbUser database.User
		result := database.DB.Where("username = ?", user).First(&dbUser)
		if result.Error != nil || !database.CheckPasswordHash(pass, dbUser.Password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Извлечение имени пользователя и репозитория из URL
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid repository path", http.StatusBadRequest)
			return
		}
		urlUsername, urlRepoName := parts[0], parts[1]

		// Проверяем, является ли репозиторий публичным
		var repo database.Repository
		result = database.DB.Where("name = ? AND owner_id = ?", urlRepoName, dbUser.ID).First(&repo)
		if result.Error != nil {
			// Если репозиторий не найден, проверяем публичные репозитории
			result = database.DB.Where("name = ? AND is_public = ?", urlRepoName, true).First(&repo)
			if result.Error != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Проверяем, соответствует ли имя пользователя в URL имени аутентифицированного пользователя или является ли пользователь админом или репозиторий публичный
		if dbUser.Username != urlUsername && !dbUser.IsAdmin && !repo.IsPublic {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("Hw?CF5Z=D#rl;3djrxKDAnE35)HOD")

func authorizeRequest(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return nil, err
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return nil, err
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}

	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, err
	}

	return claims, nil
}
