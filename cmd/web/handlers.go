/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/justinas/nosurf"
	"sanchezreneyongjoanne.net/test/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	re, err := app.ratings.Read()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// an instance of templeData
	data := &templateData{
		Ratings:         re,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	// Display the ratings using a template
	ts, err := template.ParseFiles("./ui/html/showmratings.page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *application) createMovieForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/movies.page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	movie_name := r.PostForm.Get("movie_name")
	director_name := r.PostForm.Get("director_name")
	release_date := r.PostForm.Get("release_date")
	movie_rating := r.PostForm.Get("movie_rating")
	movie_review := r.PostForm.Get("movie_review")

	errors := make(map[string]string)

	if strings.TrimSpace(movie_name) == "" {
		errors["movie_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(movie_name) > 100 {
		errors["movie_name"] = "This field is too long (maximum is 100 characters)"
	}
	if strings.TrimSpace(director_name) == "" {
		errors["director_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(director_name) > 75 {
		errors["director_name"] = "This field is too long (maximum is 75 characters)"
	}
	if strings.TrimSpace(release_date) == "" {
		errors["release_date"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(release_date) != 10 {
		errors["release_date"] = "This field must be 10 characters long"
	}
	if strings.TrimSpace(movie_rating) == "" {
		errors["movie_rating"] = "This field cannot be left blank"
	}
	if strings.TrimSpace(movie_review) == "" {
		errors["movie_review"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(movie_review) > 300 {
		errors["movie_review"] = "This field is too long (maximum is 300 characters)"
	}

	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/html/movies.page.tmpl")
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})
		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}
		return
	}
	// Insert a Rating
	id, err := app.ratings.Insert(movie_name, director_name, release_date, movie_rating, movie_review)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// set some session data
	app.session.Put(r, "flash", "Rating successfully added!")
	http.Redirect(w, r, fmt.Sprintf("/movie/%d", id), http.StatusSeeOther)
}

func (app *application) showRating(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	re, err := app.ratings.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Get/Check for the flash message
	flash := app.session.PopString(r, "flash")

	// an instance of templeData
	data := &templateData{
		Rating:          re,
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}

	// Display the Rating using a template
	ts, err := template.ParseFiles("./ui/html/mrating_page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/signup.page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	// check the web form fields to validity
	errors_user := make(map[string]string)

	if strings.TrimSpace(name) == "" {
		errors_user["name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(name) > 50 {
		errors_user["name"] = "This field is too long(maximum is 50 characters)"
	} else if utf8.RuneCountInString(name) < 5 {
		errors_user["name"] = "This field is too short(minimum is 5 characters)"
	}
	if strings.TrimSpace(email) == "" {
		errors_user["email"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(email) > 50 {
		errors_user["email"] = "This field is too long(maximum is 50 characters)"
	}

	// Check to see if the email is valid
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		errors_user["email"] = "Invalid email"
	}

	// Check the length of the password
	if strings.TrimSpace(password) == "" {
		errors_user["password"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(password) > 50 {
		errors_user["password"] = "This field is too long(maximum is 50 characters)"
	} else if utf8.RuneCountInString(password) < 10 {
		errors_user["password"] = "This field is too short(minimum is 10 characters)"
	}

	// check if there are any errors in the map
	if len(errors_user) > 0 {
		ts, err := template.ParseFiles("./ui/html/signup.page.tmpl")
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm:  errors_user,
			FormData:        r.PostForm,
			IsAuthenticated: app.isAuthenticated(r),
			CSRFToken:       nosurf.Token(r),
		})

		if err != nil {
			app.serverError(w, err)
			return
		}
		return
	}

	// Add the user to the database
	err = app.users.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			errors_user["email"] = "Email already in use"
			// rerender the signup form
			ts, err := template.ParseFiles("./ui/html/signup.page.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &templateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
				CSRFToken:       nosurf.Token(r),
			})

			if err != nil {
				app.serverError(w, err)
				return
			}
			return
		}
	}
	app.session.Put(r, "flash", "Signup was successful.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/login.page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, &templateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	})

	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	// check the web form fields to validity
	errors_user := make(map[string]string)
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			errors_user["default"] = "Email or Password is incorrect"
			// rerender the login form
			ts, err := template.ParseFiles("./ui/html/login.page.tmpl")
			if err != nil {
				app.serverError(w, err)
				return
			}
			err = ts.Execute(w, &templateData{
				ErrorsFromForm:  errors_user,
				FormData:        r.PostForm,
				IsAuthenticated: app.isAuthenticated(r),
				CSRFToken:       nosurf.Token(r),
			})

			if err != nil {
				app.serverError(w, err)
				return
			}
			return
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/movie/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
