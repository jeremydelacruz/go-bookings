package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	form := New(url.Values{})

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	form := New(url.Values{})

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when it has all required fields")
	}
}

func TestForm_Has(t *testing.T) {
	form := New(url.Values{})

	hasField := form.Has("x")
	if hasField {
		t.Error("got true for Has(x) in empty form")
	}

	postedData := url.Values{}
	postedData.Add("x", "x")
	form = New(postedData)

	hasField = form.Has("x")
	if !hasField {
		t.Error("got false for Has(x) when form should contain x")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("x", "x")
	form := New(postedData)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows valid for field value length < minimum length")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData = url.Values{}
	form = New(postedData)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows valid min length for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("x", "xyz")
	form = New(postedData)

	form.MinLength("x", 3)
	if !form.Valid() {
		t.Error("form shows invalid for field value length >= minimum length")
	}

	isError = form.Errors.Get("x")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_Email(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "xyz")
	form := New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid email for invalid email format")
	}

	postedData = url.Values{}
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "xyz@example.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid email for valid email format")
	}
}
