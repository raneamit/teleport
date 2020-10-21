package main

import (
	"context"
	"log"

	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/services"
)

// rolesCRUD performs each roles crud function as an example
func roleCRUD(ctx context.Context, client *auth.Client) {
	// create a new dev role which can log into teleport as 'teleport'
	roleName, login := "dev", "teleport"
	spec := services.RoleSpecV3{
		Allow: services.RoleConditions{
			Logins: []string{login},
		},
	}
	role, err := services.NewRole(roleName, spec)
	if err != nil {
		log.Fatalf("Failed to make new role %v", err)
	}
	err = client.UpsertRole(ctx, role)
	if err != nil {
		log.Fatalf("Failed to create role: %v", err)
	}
	log.Printf("Created Role: %v", role.GetName())

	// retrieve all roles in the cluster
	roles, err := client.GetRoles()
	if err != nil {
		log.Fatalf("Failed to retrieve roles: %v", err)
	}
	log.Println("Retrieved Roles:")
	for _, r := range roles {
		log.Printf("  %v", r.GetName())
	}

	// update the dev role to deny access to nodes labeled as production
	// you can update any attributes a role has with its setter methods
	role, err = client.GetRole(roleName)
	if err != nil {
		log.Fatal(err)
	}

	role.SetNodeLabels(services.Deny, services.Labels{"environment": []string{"production"}})
	if err = client.UpsertRole(ctx, role); err != nil {
		log.Fatalf("Failed to update role: %v", err)
	}
	log.Printf("Updated role")

	// delete the dev role we just created
	if err = client.DeleteRole(ctx, roleName); err != nil {
		log.Fatalf("Failed to delete role: %v", err)
	}
	log.Printf("Deleted role")
}
