package duplicateheader

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	Source      string   `yaml:"source"`
	Destination []string `yaml:"destination"`
	Debug       bool     `yaml:"debug"`
}

// DuplicateHeader holds the necessary components of a Traefik plugin
type DuplicateHeader struct {
	next        http.Handler
	name        string
	Source      string
	Destination []string
}

var Logger = log.New(ioutil.Discard, "DEBUG: duplicateheader: ", log.Ldate|log.Ltime|log.Lshortfile)

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

	if config.Debug {
		Logger.SetOutput(os.Stdout)
	}

	return &DuplicateHeader{
		next:        next,
		name:        name,
		Source:      config.Source,
		Destination: config.Destination,
	}, nil
}

func (d *DuplicateHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	Logger.Printf("New request: %v", req)
	Logger.Printf("Source Header: %v", d.Source)

	source := rw.Header().Get(d.Source)
	if len(source) != 0 {
		for _, dest := range d.Destination {
			Logger.Printf("Set Header: %v: %v", dest, source)
			rw.Header().Set(dest, source)
		}
	}
	d.next.ServeHTTP(rw, req)
}
