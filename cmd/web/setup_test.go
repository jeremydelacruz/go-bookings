package main

import (
	"net/http"
	"os"
	"testing"
)

type mockHandler struct{}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
