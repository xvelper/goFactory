package main

import (
	"flag"

	"gitfactory/controllers"
	"gitfactory/database"
	"log"
	"net/http"

	"gitfactory/server"
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

	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/welcome", controllers.Welcome)
	http.HandleFunc("/create_repo", controllers.CreateRepository)

	http.Handle("/", controllers.BasicAuth(http.HandlerFunc(server.Handler())))

	if err := http.ListenAndServe(server.DefaultAddress, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
