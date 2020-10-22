package main

import (
	"context"
	"crypto/tls"
	"log"
	"path/filepath"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/utils"
)

// connectClient establishes a gRPC client connected to an auth server.
func connectClient() (*auth.Client, error) {
	tlsConfig, err := setupClientTLS(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup TLS config: %v", err)
	}

	// replace 127.0.0.1:3025 (default) with your auth server address
	authServerAddr := utils.MustParseAddrList("127.0.0.1:3025")
	clientConfig := auth.ClientConfig{Addrs: authServerAddr, TLS: tlsConfig}

	return auth.NewTLSClient(clientConfig)
}

// setupClientTLS sets up client TLS authentiction between TLS client and Teleport Auth server.
func setupClientTLS(ctx context.Context) (*tls.Config, error) {
	// This function assumes you're running the api locally alongside teleport. If you are
	// running the api remotely, you'll need to provide it with the client certificates.

	// replace /var/lib/teleport (default) with your data directory filepath
	dataDir := "/var/lib/teleport"
	storage, err := auth.NewProcessStorage(ctx, filepath.Join(dataDir, teleport.ComponentProcess))
	if err != nil {
		return nil, err
	}
	defer storage.Close()

	// This uses hard coded paths to retrieve the client certificates from your teleport process
	identity, err := storage.ReadIdentity(auth.IdentityCurrent, teleport.RoleAdmin)
	if err != nil {
		return nil, err
	}

	return identity.TLSConfig(nil)
}
