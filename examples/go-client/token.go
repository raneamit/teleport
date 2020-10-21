package main

import (
	"context"
	"log"
	"time"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/lib/auth"
)

// tokenCRUD performs each token crud function as an example
func tokenCRUD(ctx context.Context, client *auth.Client) {
	// generate a cluster join token for adding another proxy to a cluster.
	tokenName := "mytoken"
	token, err := client.GenerateToken(ctx, auth.GenerateTokenRequest{
		Token: tokenName,
		Roles: teleport.Roles{teleport.RoleProxy},
		TTL:   time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	log.Printf("Generated token: %v\n", token)

	// generate a random cluster join token for adding a node to a cluster
	randToken, err := client.GenerateToken(ctx, auth.GenerateTokenRequest{
		Roles: teleport.Roles{teleport.RoleNode},
		TTL:   time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	log.Printf("Generated random token: %v\n", randToken)

	// retrieve all active cluster join tokens
	tokens, err := client.GetTokens()
	if err != nil {
		log.Fatalf("Failed to get tokens: %v", err)
	}
	log.Println("Retrieved tokens:")
	for _, t := range tokens {
		log.Printf("  %v", t.GetName())
	}

	// update a token
	provToken, err := client.GetToken(tokenName)
	if err != nil {
		log.Fatal(err)
	}

	provToken.SetRoles(teleport.Roles{teleport.RoleProxy})
	err = client.UpsertToken(provToken)
	if err != nil {
		log.Fatal(err)
	}

	// delete the cluster tokens we just created
	if err = client.DeleteToken(token); err != nil {
		log.Fatalf("Failed to delete token: %v", err)
	}
	if err = client.DeleteToken(randToken); err != nil {
		log.Fatalf("Failed to delete randToken: %v", err)
	}
	log.Printf("Deleted generated tokens\n")
}
