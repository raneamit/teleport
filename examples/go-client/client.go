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

func connectClient() (*auth.Client, error) {
	// set up client TLS authentiction between TLS client and Teleport Auth server.
	tlsConfig, err := setupClientTLS(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup TLS config: %v", err)
	}

	// replace 127.0.0.1:3025 (default) with your auth server address
	authServerAddr := []utils.NetAddr{*utils.MustParseAddr("127.0.0.1:3025")}
	clientConfig := auth.ClientConfig{Addrs: authServerAddr, TLS: tlsConfig}

	return auth.NewTLSClient(clientConfig)
}

// This function assumes you're running the api locally alongside teleport. If you are
// running the api remotely, you'll need to provide it with the client certificates.
func setupClientTLS(ctx context.Context) (*tls.Config, error) {
	// replace /var/lib/teleport (default) with your data directory filepath
	dataDir := "/home/bjoerger/gravitational/joerger/teleport/local/bin/main/data"
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
