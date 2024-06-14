package controllers

import (
	"encoding/json"
	"gitfactory/database"
	"log"
	"net/http"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserProfileResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserRequest struct {
	ID uint `json:"id"`
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// Парсим тело запроса
	var userReq UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ищем пользователя в базе данных
	var user database.User
	if err := database.DB.First(&user, userReq.ID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Формируем и отправляем ответ
	response := UserResponse{Username: user.Username}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetUserDetailsJWT(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("GetUserDetails: Authorized user %s", claims.Username)

	var user database.User
	result := database.DB.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func GetUserRepositories(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	log.Printf("Я вызываюсь")
	log.Printf("GetUserRepositories: Authorized user %s", claims.Username)

	var user database.User
	result := database.DB.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var repositories []database.Repository
	result = database.DB.Where("owner_id = ?", user.ID).Find(&repositories)
	if result.Error != nil {
		http.Error(w, "Error retrieving repositories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repositories)
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var user database.User
	result := database.DB.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userProfile := UserProfileRequest{
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

// UpdateUserProfile обработчик для обновления профиля пользователя
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var profileReq UserProfileResponse
	err = json.NewDecoder(r.Body).Decode(&profileReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user database.User
	result := database.DB.Where("username = ?", claims.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Firstname = profileReq.FirstName
	user.Lastname = profileReq.LastName
	user.Email = profileReq.Email

	if profileReq.Password != "" {
		user.Password = profileReq.Password
	}

	database.DB.Save(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
