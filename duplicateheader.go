package duplicateheader

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Source      string   `json:"source,omitempty"`
	Destination []string `json:"destination,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates and returns a plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Source) == 0 {
		return nil, fmt.Errorf("source can't be empty")
	}
	if len(config.Destination) == 0 {
		return nil, fmt.Errorf("destination can't be empty")
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		source := req.Header.Get(config.Source)
		if source != "" {
			for _, dest := range config.Destination {
				req.Header.Set(dest, source)
			}
		}
		next.ServeHTTP(rw, req)
	}), nil
}
