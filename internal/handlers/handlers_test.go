package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

type handlerTest struct {
	name               string
	path               string
	method             string
	params             []postData
	expectedStatusCode int
}

var tests = []handlerTest{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post-search", "/search-availability", "POST", []postData{
		{key: "start", value: "2023-01-01"},
		{key: "end", value: "2023-01-02"},
	}, http.StatusOK},
	{"post-search-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2023-01-01"},
		{key: "end", value: "2023-01-02"},
	}, http.StatusOK},
	{"post-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "jeremy"},
		{key: "last_name", value: "de la cruz"},
		{key: "email", value: "jeremy@example.com"},
		{key: "phone", value: "123-456-7890"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	server := httptest.NewTLSServer(routes)
	defer server.Close()

	for _, test := range tests {
		var res *http.Response
		var err error

		if test.method == "GET" {
			// todo: add a reservation to the session for /reservation-summary
			res, err = server.Client().Get(server.URL + test.path)
		} else {
			values := url.Values{}
			for _, p := range test.params {
				values.Add(p.key, p.value)
			}
			res, err = server.Client().PostForm(server.URL+test.path, values)
		}

		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, res.StatusCode)
		}
	}
}
