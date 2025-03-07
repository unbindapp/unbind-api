package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/log"
)

func main() {
	// Define command-line flags
	startAPIFlag := flag.Bool("start-api", false, "Start the API server")
	startOauth2ApiFlag := flag.Bool("start-oauth2-api", false, "Start the OAuth2 API server")

	// User management flags
	createUserCmd := flag.NewFlagSet("create-user", flag.ExitOnError)
	email := createUserCmd.String("email", "", "Email for the new user")
	password := createUserCmd.String("password", "", "Password for the new user")

	// List users flag
	listUsersFlag := flag.Bool("list-users", false, "List all users")

	// Parse the command-line flags
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Overload()
	if err != nil {
		log.Warn("Error loading .env file:", err)
	}

	cfg := config.NewConfig()

	// Check if the -start-api flag was provided
	if *startAPIFlag {
		startAPI(cfg)
	} else if *startOauth2ApiFlag {
		startOauth2Server(cfg)
	} else if *listUsersFlag {
		cli := NewCLI(cfg)
		cli.listUsers()
	} else if len(os.Args) > 1 && os.Args[1] == "create-user" {
		cli := NewCLI(cfg)
		createUserCmd.Parse(os.Args[2:])
		cli.createUser(*email, *password)
	} else {
		fmt.Println("No command specified. Use -start-api to start the API server.")
	}
}
