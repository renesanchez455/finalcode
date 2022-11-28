/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

package main

import (
	"net/url"

	"sanchezreneyongjoanne.net/test/pkg/models"
)

type templateData struct {
	CSRFToken       string
	Ratings         []*models.Rating
	Rating          *models.Rating
	ErrorsFromForm  map[string]string
	Flash           string
	FormData        url.Values
	IsAuthenticated bool
}
