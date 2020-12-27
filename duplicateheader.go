package duplicateheader

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Source      string   `yaml:"source"`
	Destination []string `yaml:"destination"`
}

// DuplicateHeader holds the necessary components of a Traefik plugin
type DuplicateHeader struct {
	next        http.Handler
	name        string
	Source      string
	Destination []string
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

	var DestCanonicalHeaderKey []string
	for _, dest := range config.Destination {
		DestCanonicalHeaderKey = append(DestCanonicalHeaderKey, http.CanonicalHeaderKey(dest))
	}

	return &DuplicateHeader{
		next:        next,
		name:        name,
		Source:      http.CanonicalHeaderKey(config.Source),
		Destination: DestCanonicalHeaderKey,
	}, nil
}

func (d *DuplicateHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if source, ok := req.Header[d.Source]; ok {
		if len(source) != 0 {
			for _, dest := range d.Destination {
				req.Header[dest] = source
			}
		}
	}
	d.next.ServeHTTP(rw, req)
}
