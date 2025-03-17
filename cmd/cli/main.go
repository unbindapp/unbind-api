package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

func main() {
	// User management flags
	createUserCmd := flag.NewFlagSet("create-user", flag.ExitOnError)
	email := createUserCmd.String("email", "", "Email for the new user")
	password := createUserCmd.String("password", "", "Password for the new user")

	// Group management flags
	createGroupCmd := flag.NewFlagSet("create-group", flag.ExitOnError)
	groupName := createGroupCmd.String("name", "", "Name of the group")
	groupDescription := createGroupCmd.String("description", "", "Description of the group")
	teamID := createGroupCmd.String("team-id", "", "ID of the team to associate the group with (Optional)")

	// Add user to group
	addUserToGroupCmd := flag.NewFlagSet("add-user-to-group", flag.ExitOnError)
	addToGroupUserEmail := addUserToGroupCmd.String("email", "", "ID of the user to add to the group")
	addToGroupName := addUserToGroupCmd.String("group-name", "", "ID of the group to add the user to")

	// Grant permisisons to group
	// e.g.: /app/server grant-group-permissions -group-name=admin -action=admin -resource-type=team -resource-id="*"
	grantGroupPermissionsCmd := flag.NewFlagSet("grant-group-permissions", flag.ExitOnError)
	grantGroupName := grantGroupPermissionsCmd.String("group-name", "", "Name of the group to grant permissions to")
	grantAction := grantGroupPermissionsCmd.String("action", "", "Action to grant permissions for")
	grantResourceType := grantGroupPermissionsCmd.String("resource-type", "", "Resource type to grant permissions for")
	grantResourceID := grantGroupPermissionsCmd.String("resource-id", "", "Resource ID to grant permissions for")
	grantScope := grantGroupPermissionsCmd.String("scope", "", "Scope to grant permissions for")

	// List users flag
	listUsersFlag := flag.Bool("list-users", false, "List all users")

	// List groups flag
	listGroupsFlag := flag.Bool("list-groups", false, "List all groups")

	// List group permissions
	listGroupPermissionsCmd := flag.NewFlagSet("list-group-permissions", flag.ExitOnError)
	listGroupPermissionsGroupName := listGroupPermissionsCmd.String("group-name", "", "Name of the group to list permissions for")

	// Create team
	createTeamCmd := flag.NewFlagSet("create-team", flag.ExitOnError)
	teamName := createTeamCmd.String("name", "", "Name of the team")
	displayname := createTeamCmd.String("displayname", "", "Display name of the team")

	// Sync permissions with Kubernetes
	syncGroupToK8sFlag := flag.Bool("sync-group-to-k8s", false, "Sync group permissions with Kubernetes")

	// Parse the command-line flags
	flag.Parse()

	// Load environment variables from .env file
	err := godotenv.Overload()
	if err != nil {
		log.Warn("Error loading .env file:", err)
	}

	cfg := config.NewConfig()

	// Processs flags
	if *listUsersFlag {
		cli := NewCLI(cfg)
		cli.listUsers()
	} else if *listGroupsFlag {
		cli := NewCLI(cfg)
		cli.listGroups()
	} else if len(os.Args) > 1 && os.Args[1] == "create-user" {
		cli := NewCLI(cfg)
		createUserCmd.Parse(os.Args[2:])
		cli.createUser(*email, *password)
	} else if len(os.Args) > 1 && os.Args[1] == "create-group" {
		cli := NewCLI(cfg)
		createGroupCmd.Parse(os.Args[2:])
		if teamID != nil {
			cli.createGroup(*groupName, *groupDescription, nil)
		} else {
			parsedTeamID, err := uuid.Parse(*teamID)
			if err != nil {
				log.Fatalf("Failed to parse team ID: %v", err)
			}
			cli.createGroup(*groupName, *groupDescription, &parsedTeamID)
		}
	} else if len(os.Args) > 1 && os.Args[1] == "add-user-to-group" {
		cli := NewCLI(cfg)
		addUserToGroupCmd.Parse(os.Args[2:])
		cli.addUserToGroup(*addToGroupUserEmail, *addToGroupName)
	} else if len(os.Args) > 1 && os.Args[1] == "list-group-permissions" {
		cli := NewCLI(cfg)
		listGroupPermissionsCmd.Parse(os.Args[2:])
		cli.listGroupPermissions(*listGroupPermissionsGroupName)
	} else if len(os.Args) > 1 && os.Args[1] == "grant-group-permissions" {
		cli := NewCLI(cfg)
		grantGroupPermissionsCmd.Parse(os.Args[2:])
		cli.grantPermission(*grantGroupName, *grantAction, *grantResourceType, *grantResourceID, *grantScope)
	} else if *syncGroupToK8sFlag {
		cli := NewCLI(cfg)
		cli.syncPermissionsWithK8S()
	} else if len(os.Args) > 1 && os.Args[1] == "create-team" {
		cli := NewCLI(cfg)
		createTeamCmd.Parse(os.Args[2:])
		cli.createTeam(*teamName, *displayname)
	} else {
		fmt.Println("No command specified. Use -start-api to start the API server.")
	}
}
