/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

package models

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound     = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// A struct to hold a rating
type Rating struct {
	Id            int
	Movie_name    string
	Director_name string
	Release_date  time.Time
	Movie_rating  int
	Movie_review  string
}

// A struct to hold a user
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
