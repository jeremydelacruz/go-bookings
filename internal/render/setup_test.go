package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jeremydelacruz/go-bookings/internal/config"
	"github.com/jeremydelacruz/go-bookings/internal/models"
)

var testApp config.AppConfig
var session *scs.SessionManager

type testWriter struct{}

func (tw *testWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *testWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (tw *testWriter) WriteHeader(statusCode int) {}

func TestMain(m *testing.M) {
	testApp.InProduction = false

	// configure loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// Register this type to use in the session
	gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}
