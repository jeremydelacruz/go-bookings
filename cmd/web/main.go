package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jeremydelacruz/go-bookings/internal/config"
	"github.com/jeremydelacruz/go-bookings/internal/driver"
	"github.com/jeremydelacruz/go-bookings/internal/handlers"
	"github.com/jeremydelacruz/go-bookings/internal/helpers"
	"github.com/jeremydelacruz/go-bookings/internal/models"
	"github.com/jeremydelacruz/go-bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the application entrypoint
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Printf("starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// run configures the app config, session, and handlers
func run() (*driver.DB, error) {
	// change this to true in prod
	app.InProduction = false

	// configure loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// register this type to use in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=jdelacruz password=")
	if err != nil {
		return nil, fmt.Errorf("run: failed connecting to database: %w", err)
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("run: failed creating template cache: %w", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
