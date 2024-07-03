// Package plugindemo a demo plugin.
package plugindemo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	name     string
	template *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	splitHeaderJson(config.Headers)
	delete(config.Headers, "X-Demo")

	return &Demo{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("template").Delims("[[", "]]"),
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for key, value := range a.headers {
		tmpl, err := a.template.Parse(value)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		writer := &bytes.Buffer{}

		err = tmpl.Execute(writer, req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set(key, writer.String())
	}

	a.next.ServeHTTP(rw, req)
}
func splitHeaderJson(headers map[string]string) Demo {
	fmt.Println("Cabeceras iniciales: ", headers)
	jsonStr := headers["X-Demo"]

	// Parsea el JSON original a un mapa
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Demo{}
	}
	for key, value := range obj {
		headers["S-"+key] = fmt.Sprint(value)
	}

	return Demo{headers: headers}
}
