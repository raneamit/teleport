package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/services"
)

// rolesCRUD performs each roles crud function as an example
func roleCRUD(ctx context.Context, client *auth.Client) error {
	// create a new dev role which can log into teleport as 'teleport'
	roleSpec := services.RoleSpecV3{
		Allow: services.RoleConditions{
			Logins: []string{"teleport"},
		},
	}

	role, err := services.NewRole("dev", roleSpec)
	if err != nil {
		return fmt.Errorf("Failed to make new role %v", err)
	}

	err = client.UpsertRole(ctx, role)
	if err != nil {
		return fmt.Errorf("Failed to create role: %v", err)
	}
	log.Printf("Created Role: %v", role.GetName())

	// retrieve all roles in the cluster
	roles, err := client.GetRoles()
	if err != nil {
		return fmt.Errorf("Failed to retrieve roles: %v", err)
	}
	log.Println("Retrieved Roles:")
	for _, r := range roles {
		log.Printf("  %v", r.GetName())
	}

	// updates the dev role to deny access to nodes labeled as production
	role, err = client.GetRole("dev")
	if err != nil {
		return fmt.Errorf("Failed to retrieve role for updating: %v", err)
	}

	role.SetNodeLabels(services.Deny, services.Labels{"environment": []string{"production"}})
	if err = client.UpsertRole(ctx, role); err != nil {
		return fmt.Errorf("Failed to update role: %v", err)
	}
	log.Printf("Updated role")

	// delete the dev role we just created
	if err = client.DeleteRole(ctx, "dev"); err != nil {
		return fmt.Errorf("Failed to delete role: %v", err)
	}
	log.Printf("Deleted role")

	return nil
}
