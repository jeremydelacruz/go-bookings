package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mHandler mockHandler
	h := NoSurf(&mHandler)
	switch h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error("return type is not http.Handler")
	}
}

func TestSessionLoad(t *testing.T) {
	var mHandler mockHandler
	h := SessionLoad(&mHandler)
	switch h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error("return type is not http.Handler")
	}
}
