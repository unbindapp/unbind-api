package main

import (
	"context"
	"fmt"
	"time"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/log"
	"golang.org/x/crypto/bcrypt"
)

type cli struct {
	repository *repository.Repository
}

func NewCLI(cfg *config.Config) *cli {
	// Load database
	dbConnInfo, err := database.GetSqlDbConn(cfg, false)
	if err != nil {
		log.Fatalf("Failed to get database connection info: %v", err)
	}
	// Initialize ent client
	db, err := database.NewEntClient(dbConnInfo)
	if err != nil {
		log.Fatalf("Failed to create ent client: %v", err)
	}
	log.Info("🦋 Running migrations...")
	if err := db.Schema.Create(context.TODO()); err != nil {
		log.Fatal("Failed to run migrations", "err", err)
	}
	repo := repository.NewRepository(db)

	return &cli{
		repository: repo,
	}
}

// List all users in the database
func (self *cli) listUsers() {
	// Query all users
	users, err := self.repository.DB.User.Query().All(context.Background())
	if err != nil {
		fmt.Printf("Error querying users: %v\n", err)
		return
	}

	// Print user information
	fmt.Println("Users:")
	fmt.Println("-------------------------------------")
	for _, u := range users {
		fmt.Printf("ID: %s\n", u.ID)
		fmt.Printf("Email: %s\n", u.Email)
		fmt.Printf("Created: %s\n", u.CreatedAt.Format(time.RFC3339))
		fmt.Println("-------------------------------------")
	}
	fmt.Printf("Total users: %d\n", len(users))
}

// Create a new user
func (self *cli) createUser(email, password string) {
	// Validate inputs
	if email == "" || password == "" {
		log.Errorf("Error: email and password are required")
		return
	}

	// Check if username already exists
	u, err := self.repository.GetUserByEmail(context.Background(), email)
	if err != nil {
		if !ent.IsNotFound(err) {
			log.Errorf("Error checking if user exists: %v", err)
			return
		}
	}
	if err == nil {
		log.Errorf("Error: Email '%s' already exists", email)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	// Create the user
	u, err = self.repository.DB.User.Create().
		SetEmail(email).
		SetPasswordHash(string(hashedPassword)).
		Save(context.Background())

	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}

	fmt.Println("User created successfully:")
	fmt.Printf("ID: %s\n", u.ID)
	fmt.Printf("Email: %s\n", u.Email)
}
