package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jeremydelacruz/go-bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Error("return type is not *chi.Mux")
	}
}
