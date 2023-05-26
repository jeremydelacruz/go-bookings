package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jeremydelacruz/go-bookings/internal/models"
)

type handlerTest struct {
	name               string
	path               string
	method             string
	expectedStatusCode int
}

var tests = []handlerTest{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"generals", "/generals-quarters", "GET", http.StatusOK},
	{"majors", "/majors-suite", "GET", http.StatusOK},
	{"search", "/search-availability", "GET", http.StatusOK},
}

var urlEncoded = "application/x-www-form-urlencoded"

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	server := httptest.NewTLSServer(routes)
	defer server.Close()

	for _, test := range tests {
		res, err := server.Client().Get(server.URL + test.path)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, res.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	handler := http.HandlerFunc(Repo.Reservation)

	// test when reservation is in the session
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	resRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusOK)
	}

	// test when reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test when room is non-existent
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	resRecorder = httptest.NewRecorder()
	reservation.RoomID = 999
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// initialize reservation data
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-01-02")
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
	}
	badReservation := models.Reservation{
		RoomID:    999,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// initialize url-encoded form fields
	validBody := url.Values{}
	validBody.Add("first_name", "Jane")
	validBody.Add("last_name", "Doe")
	validBody.Add("email", "jane@doe.com")
	validBody.Add("phone", "1234567890")

	// initialize handler as testable HandlerFunc
	handler := http.HandlerFunc(Repo.PostReservation)

	// test existing reservation, successful form parse, valid form, successful DB insert
	validReader := strings.NewReader(validBody.Encode())
	req, _ := http.NewRequest("POST", "/make-reservation", validReader)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	session.Put(ctx, "reservation", reservation)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostReservation happy path unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusSeeOther)
	}

	// test non-existent reservation
	validReader.Seek(0, 0)
	req, _ = http.NewRequest("POST", "/make-reservation", validReader)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation non-existent reservation unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test non-existent form body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	session.Put(ctx, "reservation", reservation)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation non-existent form body unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test invalid form body
	invalidBody := url.Values{}
	invalidBody.Add("first_name", "T")
	invalidBody.Add("last_name", "Doe")
	invalidBody.Add("email", "jane@doe.com")
	invalidBody.Add("phone", "1234567890")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(invalidBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	session.Put(ctx, "reservation", reservation)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostReservation invalid form field unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusSeeOther)
	}

	// test unsuccessful reservation DB insert
	validReader.Seek(0, 0)
	req, _ = http.NewRequest("POST", "/make-reservation", validReader)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	badReservation.RoomID = 999
	session.Put(ctx, "reservation", badReservation)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation unsuccessful reservation insert unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test unsuccessful room restriction DB insert
	validReader.Seek(0, 0)
	req, _ = http.NewRequest("POST", "/make-reservation", validReader)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	badReservation.RoomID = 1000
	session.Put(ctx, "reservation", badReservation)
	resRecorder = httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation unsuccessful room restriction insert unexpected response code: got %d, expected %d", resRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// initialize url-encoded form fields
	reqBody := url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-01")
	reqBody.Add("room_id", "1")

	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	var res jsonResponse

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	err := json.Unmarshal(resRecorder.Body.Bytes(), &res)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// initialize url-encoded form fields
	reqBody := url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-01")

	handler := http.HandlerFunc(Repo.PostAvailability)

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", urlEncoded)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusSeeOther {
		t.Errorf("got status code: %d, expected: %d", resRecorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	handler := http.HandlerFunc(Repo.ReservationSummary)

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	handler := http.HandlerFunc(Repo.ChooseRoom)

	req, _ := http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/1"
	ctx := getCtx(req)
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusSeeOther {
		t.Errorf("got status code: %d, expected: %d", resRecorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	id := "id=1"
	start := "s=2050-01-01"
	end := "e=2050-01-01"
	reqURI := fmt.Sprintf("/book-room?%s&%s&%s", id, start, end)

	handler := http.HandlerFunc(Repo.BookRoom)

	req, _ := http.NewRequest("GET", reqURI, nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	resRecorder := httptest.NewRecorder()

	handler.ServeHTTP(resRecorder, req)
	if resRecorder.Code != http.StatusSeeOther {
		t.Errorf("got status code: %d, expected: %d", resRecorder.Code, http.StatusSeeOther)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
