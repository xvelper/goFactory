package controllers

import (
	"encoding/json"
	"gitfactory/database"
	"gitfactory/server"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type RepositoryRequest struct {
	RepoName string `json:"repo_name"`
	IsPublic bool   `json:"is_public"`
}
type DeleteRepositoryRequest struct {
	ID string `json:"id"`
}

type Commit struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
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
		UUID:     repoUUID,
		Name:     repoReq.RepoName + ".git",
		OwnerID:  user.ID,
		Path:     repoDir,
		IsPublic: repoReq.IsPublic,
	}

	result = database.DB.Create(&repo)
	if result.Error != nil {
		http.Error(w, "Error saving repository to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Repository created successfully", "uuid": repoUUID})
}

func DeleteRepository(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var deleteReq struct {
		ID uint `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&deleteReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var repo database.Repository
	result := database.DB.Where("id = ? AND owner_id = ?", deleteReq.ID, claims.Id).First(&repo)
	if result.Error != nil {
		http.Error(w, "Repository not found or unauthorized", http.StatusNotFound)
		return
	}

	err = os.RemoveAll(repo.Path)
	if err != nil {
		http.Error(w, "Error deleting repository directory", http.StatusInternalServerError)
		return
	}

	result = database.DB.Delete(&repo)
	if result.Error != nil {
		http.Error(w, "Error deleting repository from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Repository deleted successfully"})
}

func GetRepositoryCommits(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var repoReq DeleteRepositoryRequest
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

	// Получение репозитория по имени и владельцу
	var repo database.Repository
	result = database.DB.Where("id = ? AND owner_id = ?", repoReq.ID, claims.Id).First(&repo)
	if result.Error != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	// Выполнение команды git log для получения истории коммитов
	repoPath := filepath.Join(repo.Path)
	cmd := exec.Command("git", "-C", repoPath, "log", "--pretty=format:%H|%an|%ad|%s", "--date=short")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Error fetching git log", http.StatusInternalServerError)
		return
	}

	// Парсинг истории коммитов
	var commits []Commit
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}
		commit := Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		}
		commits = append(commits, commit)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commits)
}

func GetRepositoriesByUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var repos []database.Repository
	result := database.DB.Where("owner_id = ?", req.ID).Find(&repos)
	if result.Error != nil {
		http.Error(w, "Error fetching repositories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repos)
}

func GetPublicRepositories(w http.ResponseWriter, r *http.Request) {
	var repos []database.Repository
	result := database.DB.Where("is_public = ?", true).Find(&repos)
	if result.Error != nil {
		http.Error(w, "Error fetching public repositories", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repos)
}

func GetRepoLanguage(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var repo database.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userPath := filepath.Join(server.DefaultConfig.ProjectRoot, claims.Username, repo.Name+".git")
	language, err := detectLanguage(userPath)
	if err != nil {
		http.Error(w, "Failed to detect language", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"language": language})
}

// GetRepoStructure function to get the file structure of the repository
func GetRepoStructure(w http.ResponseWriter, r *http.Request) {
	claims, err := authorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var repo database.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userPath := filepath.Join(server.DefaultConfig.ProjectRoot, claims.Username, repo.Name+".git")
	structure, err := getFileStructure(userPath)
	if err != nil {
		http.Error(w, "Failed to get file structure", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(structure)
}

// detectLanguage function to determine the primary language of the repository
func detectLanguage(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fileExtensions := make(map[string]int)
	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}
		ext := filepath.Ext(file)
		fileExtensions[ext]++
	}

	var primaryLanguage string
	var maxCount int
	for ext, count := range fileExtensions {
		if count > maxCount {
			maxCount = count
			primaryLanguage = ext
		}
	}

	return primaryLanguage, nil
}

// getFileStructure function to get the file structure of the repository
func getFileStructure(repoPath string) (map[string]interface{}, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-tree", "-r", "HEAD", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	structure := make(map[string]interface{})
	for _, file := range files {
		if file == "" {
			continue
		}
		parts := strings.Split(file, "/")
		insertIntoMap(structure, parts)
	}

	return structure, nil
}

// insertIntoMap function to insert file parts into the structure map
func insertIntoMap(m map[string]interface{}, parts []string) {
	if len(parts) == 1 {
		m[parts[0]] = struct{}{}
		return
	}

	if _, ok := m[parts[0]]; !ok {
		m[parts[0]] = make(map[string]interface{})
	}

	if nextMap, ok := m[parts[0]].(map[string]interface{}); ok {
		insertIntoMap(nextMap, parts[1:])
	}
}
