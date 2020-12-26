package duplicateheader

import (
	"context"
	"fmt"
	"net"
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
		if len(source) != 0 && net.ParseIP(source[0]) != nil {
			for _, dest := range d.Destination {
				fmt.Printf("set %s as %s\n", source[0], dest)
				req.Header.Set(dest, source[0])
			}
		}
	}
	d.next.ServeHTTP(rw, req)
}
