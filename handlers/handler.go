package handlers

import (
	"goFactory/models"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RepoRequest struct {
	Username string `json:"username"`
	Repo     string `json:"repo"`
}

func CreateUserHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		user := models.User{Username: req.Username, Password: req.Password}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
		}

		userRepoPath := filepath.Join("repos", user.Username)
		if err := os.MkdirAll(userRepoPath, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user repo directory"})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// InitRepositoryHandler обрабатывает HTTP запросы для инициализации новых репозиториев
func InitRepositoryHandler(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Repo     string `json:"repo"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON")
	}

	repoPath := filepath.Join("repos", req.Username, req.Repo+".git")

	// Создаём директорию для репозитория, если она ещё не существует
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		if err := os.MkdirAll(repoPath, 0775); err != nil {
			log.Printf("Failed to create repository directory: %s\n", repoPath)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create repository directory")
		}
	}

	// Инициализируем новый репозиторий с shared доступом
	r, err := git.PlainInit(repoPath, true)
	if err != nil {
		log.Printf("Failed to initialize repository: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to initialize repository")
	}

	// Настройка репозитория для shared доступа
	cfg, err := r.Config()
	if err == nil {
		cfg.Core.SharedRepository = gitconfig.SharedGroup
		err = r.Storer.SetConfig(cfg)
		if err != nil {
			log.Printf("Failed to configure repository: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to configure repository")
		}
	}

	return c.Status(fiber.StatusOK).SendString("Repository initialized successfully")
}

func DeleteRepositoryHandler(c *fiber.Ctx) error {
	username := c.Params("username")
	repo := c.Params("repo") + ".git" // Дополняем имя репозитория суффиксом .git для ясности
	repoPath := filepath.Join("repos", username, repo)

	// Проверяем, существует ли директория репозитория
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("Repository not found: %s\n", repoPath)
		return c.Status(fiber.StatusNotFound).SendString("Repository not found")
	}

	// Удаление директории репозитория и всех вложенных файлов
	err := os.RemoveAll(repoPath)
	if err != nil {
		log.Printf("Failed to delete repository: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete repository")
	}

	return c.SendString("Repository deleted successfully")
}

func GetRepoHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var repos []models.Repository
		if err := db.Find(&repos).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch repos"})
		}

		return c.JSON(repos)
	}
}

func GitHTTPHandler(c *fiber.Ctx) error {
	username := c.Params("username")
	repo := c.Params("repo")
	repoPath := filepath.Join("repos", username, repo)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("Repository not found: %s\n", repoPath)
		return c.Status(fiber.StatusNotFound).SendString("Repository not found")
	}

	service := c.Query("service")
	if service == "git-upload-pack" {
		service = "upload-pack"
	} else if service == "git-receive-pack" {
		service = "receive-pack"
	} else {
		log.Printf("Unsupported service: %s\n", service)
		return c.Status(fiber.StatusBadRequest).SendString("Unsupported Git service requested")
	}

	cmd := exec.Command("git", service, "--stateless-rpc", repoPath)
	cmd.Env = append(os.Environ(), "GIT_HTTP_EXPORT_ALL=1", "GIT_PROJECT_ROOT="+repoPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to create stdout pipe: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start git command: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Git command failed to start")
	}

	c.Context().SetContentType("application/x-git-" + service + "-result")
	if _, err := io.Copy(c.Response().BodyWriter(), stdout); err != nil {
		log.Printf("Failed to send data to client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send data")
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git command failed with error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Git command failed")
	}

	return nil
}
