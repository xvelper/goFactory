package controllers

import (
	"encoding/json"
	"gitfactory/database"
	"gitfactory/server"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type RepositoryRequest struct {
	RepoName string `json:"repo_name"`
}

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var repoReq RepositoryRequest
	err = json.NewDecoder(r.Body).Decode(&repoReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Получение пользователя из базы данных
	var user database.User
	result := database.DB.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	repoUUID := uuid.New().String()
	repoDir := filepath.Join(server.DefaultConfig.ProjectRoot, user.Username, repoReq.RepoName+".git")
	err = os.MkdirAll(repoDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating repository directory", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("git", "init", "--bare", repoDir)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Error initializing git repository", http.StatusInternalServerError)
		return
	}

	repo := database.Repository{
		UUID:    repoUUID,
		Name:    repoReq.RepoName + ".git",
		OwnerID: user.ID,
		Path:    repoDir,
	}

	result = database.DB.Create(&repo)
	if result.Error != nil {
		http.Error(w, "Error saving repository to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Repository created successfully", "uuid": repoUUID})
}

func authorizeRequest(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, json.NewEncoder(w).Encode(map[string]string{"message": "No token provided"})
		}
		return nil, json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
		}
		return nil, json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
	}

	if !token.Valid {
		return nil, json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
	}

	return claims, nil
}
