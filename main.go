package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"gitfactory/controllers"
	"gitfactory/database"
	"gitfactory/server"

	"github.com/gorilla/handlers"
)

func init() {
	flag.StringVar(&server.DefaultConfig.AuthPassEnvVar, "auth_pass_env_var", server.DefaultConfig.AuthPassEnvVar, "set an env var to provide the basic auth pass as")
	flag.StringVar(&server.DefaultConfig.AuthUserEnvVar, "auth_user_env_var", server.DefaultConfig.AuthUserEnvVar, "set an env var to provide the basic auth user as")
	flag.StringVar(&server.DefaultConfig.DefaultEnv, "default_env", server.DefaultConfig.DefaultEnv, "set the default env")
	flag.StringVar(&server.DefaultConfig.ProjectRoot, "project_root", server.DefaultConfig.ProjectRoot, "set project root")
	flag.StringVar(&server.DefaultConfig.GitBinPath, "git_bin_path", server.DefaultConfig.GitBinPath, "set git bin path")
	flag.StringVar(&server.DefaultAddress, "server_address", server.DefaultAddress, "set server address")
	flag.StringVar(&server.DefaultConfig.RoutePrefix, "route_prefix", server.DefaultConfig.RoutePrefix, "prepend a regex prefix to each git-http-backend route")
}

func main() {
	flag.Parse()

	database.ConnectDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", controllers.Register)
	mux.HandleFunc("/api/v1/login", controllers.Login)
	mux.HandleFunc("/api/v1/welcome", controllers.Welcome)
	mux.HandleFunc("/api/v1/create_repo", controllers.CreateRepository)
	mux.HandleFunc("/api/v1/delete_repo", controllers.DeleteRepository)
	mux.HandleFunc("/api/v1/get_commits", controllers.GetRepositoryCommits)
	mux.HandleFunc("/api/v1/public_repos", controllers.GetPublicRepositories)
	mux.HandleFunc("/api/v1/user_repos", controllers.GetUserRepositories)
	mux.HandleFunc("/api/v1/user_details_jwt", controllers.GetUserDetailsJWT)
	mux.HandleFunc("/api/v1/user_details", controllers.GetUserDetails)
	mux.Handle("/", controllers.BasicAuth(http.HandlerFunc(server.Handler())))

	// Настройка CORS
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	credentialsOk := handlers.AllowCredentials()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(mux)))
}
