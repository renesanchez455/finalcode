/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq" //Third party package
	"sanchezreneyongjoanne.net/test/pkg/models/postgresql"
)

func setUpDB(dsn string) (*sql.DB, error) {
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Test our connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

//Dependencies (things/variables)
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	ratings  *postgresql.RatingModel
	session  *sessions.Session
	users    *postgresql.UserModel
}

func main() {
	// Create a command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("MOVIE_DB_DSN"), "PostgreSQl DSN (Data Source Name)")
	secret := flag.String("secret", "p7Mhd+qQamgHsS*+8Tg7mNXtclui@tyz", "Secret Key")
	flag.Parse()
	// Create a logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var db, err = setUpDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Always do this before exiting
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // encrypted session cookies

	// ECDHE - Elliptic Curve Diffie-Hellman
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		ratings: &postgresql.RatingModel{
			DB: db,
		},
		session: session,
		users: &postgresql.UserModel{
			DB: db,
		},
	}

	//Create a custom web server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start our server
	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
