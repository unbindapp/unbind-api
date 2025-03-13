package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/user"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/k8s"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
	group_service "github.com/unbindapp/unbind-api/internal/services/group"
	"golang.org/x/crypto/bcrypt"
)

type cli struct {
	repository   repositories.RepositoriesInterface
	groupService *group_service.GroupService
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
	repo := repositories.NewRepositories(db)

	kubeClient := k8s.NewKubeClient(cfg)
	rbacManager := k8s.NewRBACManager(repo, kubeClient)

	return &cli{
		repository: repo,
		groupService: group_service.NewGroupService(
			repo,
			rbacManager,
		),
	}
}

// List all users in the database
func (self *cli) listUsers() {
	// Query all users
	users, err := self.repository.Ent().User.Query().WithGroups().All(context.Background())
	if err != nil {
		fmt.Printf("Error querying users: %v\n", err)
		return
	}

	// Print user information
	fmt.Println("Users:")
	fmt.Println("-------------------------------------")
	for _, u := range users {
		groupNames := make([]string, len(u.Edges.Groups))
		for i, g := range u.Edges.Groups {
			groupNames[i] = g.Name
		}
		fmt.Printf("ID: %s\n", u.ID)
		fmt.Printf("Email: %s\n", u.Email)
		fmt.Printf("Created: %s\n", u.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Groups: %s\n", strings.Join(groupNames, ", "))
		fmt.Println("-------------------------------------")
	}
	fmt.Printf("Total users: %d\n", len(users))
}

// List all groups in the database
func (self *cli) listGroups() {
	ctx := context.Background()

	// Query all groups
	groups, err := self.repository.Ent().Group.Query().All(ctx)
	if err != nil {
		fmt.Printf("Error querying groups: %v\n", err)
		return
	}

	// Print group information
	fmt.Println("Groups:")
	fmt.Println("-------------------------------------")
	for _, g := range groups {
		fmt.Printf("ID: %s\n", g.ID)
		fmt.Printf("Name: %s\n", g.Name)
		fmt.Printf("Description: %s\n", g.Description)

		if g.TeamID != nil {
			fmt.Printf("Team ID: %s\n", *g.TeamID)
		} else {
			fmt.Printf("Scope: Global\n")
		}

		if g.K8sRoleName != "" {
			fmt.Printf("K8s Role: %s\n", g.K8sRoleName)
		}

		// Get members count
		members, err := self.repository.Ent().Group.QueryUsers(g).Count(ctx)
		if err != nil {
			fmt.Printf("Error counting members: %v\n", err)
		} else {
			fmt.Printf("Members: %d\n", members)
		}

		// Get permissions count
		perms, err := self.repository.Ent().Group.QueryPermissions(g).Count(ctx)
		if err != nil {
			fmt.Printf("Error counting permissions: %v\n", err)
		} else {
			fmt.Printf("Permissions: %d\n", perms)
		}

		fmt.Println("-------------------------------------")
	}
	fmt.Printf("Total groups: %d\n", len(groups))
}

// Create a new group
func (self *cli) createGroup(name, description string, teamID *uuid.UUID) {
	ctx := context.Background()

	// Validate inputs
	if name == "" {
		log.Errorf("Error: group name is required")
		return
	}

	// Check if group already exists
	exists, err := self.repository.Ent().Group.Query().
		Where(group.NameEQ(name)).
		Exist(ctx)
	if err != nil {
		log.Errorf("Error checking if group exists: %v", err)
		return
	}
	if exists {
		log.Errorf("Error: Group '%s' already exists", name)
		return
	}

	// Create the group
	groupBuilder := self.repository.Ent().Group.Create().
		SetName(name).
		SetDescription(description)

	// Set team ID if provided
	if teamID != nil {
		// Check if team exists
		teamExists, err := self.repository.Ent().Team.Query().
			Where(team.IDEQ(*teamID)).
			Exist(ctx)
		if err != nil {
			log.Errorf("Error checking if team exists: %v", err)
			return
		}
		if !teamExists {
			log.Errorf("Error: Team with ID '%s' does not exist", teamID.String())
			return
		}

		groupBuilder.SetTeamID(*teamID)
	}

	// Save the group
	group, err := groupBuilder.Save(ctx)
	if err != nil {
		fmt.Printf("Error creating group: %v\n", err)
		return
	}

	fmt.Println("Group created successfully:")
	fmt.Printf("ID: %s\n", group.ID)
	fmt.Printf("Name: %s\n", group.Name)
	fmt.Printf("Description: %s\n", group.Description)
	if group.TeamID != nil {
		fmt.Printf("Team ID: %s\n", *group.TeamID)
	} else {
		fmt.Printf("Scope: Global\n")
	}
}

// Add a user to a group
func (self *cli) addUserToGroup(userEmail, groupName string) {
	ctx := context.Background()

	// Get the user
	dbUser, err := self.repository.User().GetByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("Error: User '%s' not found: %v\n", userEmail, err)
		return
	}

	// Get the group
	group, err := self.repository.Ent().Group.Query().
		Where(group.NameEQ(groupName)).
		Only(ctx)
	if err != nil {
		fmt.Printf("Error: Group '%s' not found: %v\n", groupName, err)
		return
	}

	// Check if user is already in the group
	inGroup, err := self.repository.Ent().Group.QueryUsers(group).
		Where(user.IDEQ(dbUser.ID)).
		Exist(ctx)
	if err != nil {
		fmt.Printf("Error checking if user is in group: %v\n", err)
		return
	}
	if inGroup {
		fmt.Printf("User '%s' is already in group '%s'\n", userEmail, groupName)
		return
	}

	// Add user to group
	err = self.repository.Ent().Group.UpdateOne(group).
		AddUserIDs(dbUser.ID).
		Exec(ctx)
	if err != nil {
		fmt.Printf("Error adding user to group: %v\n", err)
		return
	}

	fmt.Printf("Successfully added user '%s' to group '%s'\n", userEmail, groupName)
}

// List permissions for a group
func (self *cli) listGroupPermissions(groupName string) {
	ctx := context.Background()

	// Get the group
	group, err := self.repository.Ent().Group.Query().
		Where(group.NameEQ(groupName)).
		Only(ctx)
	if err != nil {
		fmt.Printf("Error: Group '%s' not found: %v\n", groupName, err)
		return
	}

	// Get permissions
	perms, err := self.repository.Ent().Group.QueryPermissions(group).All(ctx)
	if err != nil {
		fmt.Printf("Error querying permissions: %v\n", err)
		return
	}

	// Print permissions
	fmt.Printf("Permissions for group '%s':\n", groupName)
	fmt.Println("-------------------------------------")
	for i, p := range perms {
		fmt.Printf("%d. %s %s:%s\n", i+1, p.Action, p.ResourceType, p.ResourceID)
		if p.Scope != "" {
			fmt.Printf("   Scope: %s\n", p.Scope)
		}
	}
	fmt.Println("-------------------------------------")
	fmt.Printf("Total permissions: %d\n", len(perms))
}

// Create a new user
func (self *cli) createUser(email, password string) {
	// Validate inputs
	if email == "" || password == "" {
		log.Errorf("Error: email and password are required")
		return
	}

	// Check if username already exists
	u, err := self.repository.User().GetByEmail(context.Background(), email)
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
	u, err = self.repository.Ent().User.Create().
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
