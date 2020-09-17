package lightstep

import (
	"fmt"
	"log"

	openapi "github.com/go-openapi/runtime/client"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client"
)

// Config defines the configuration options for the Lightstep provider
type Config struct {
	// The Lightstep API Key
	APIKey string

	// The LightStep Organization
	Organization string

	// The LightStep Project
	Project string
}

// Client returns a new Lightstep client
func (c *Config) Client() (*client.LightstepAPI, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("LightStep API Key is mandatory")
	}

	config := client.Config{
		AuthInfo: openapi.APIKeyAuth("Authorization", "header", c.APIKey),
	}

	client := client.New(config)

	log.Printf("[INFO] Lightstep client configured")

	return client, nil
}
