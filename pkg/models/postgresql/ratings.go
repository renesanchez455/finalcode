/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

package postgresql

import (
	"database/sql"
	"errors"

	"sanchezreneyongjoanne.net/test/pkg/models"
)

type RatingModel struct {
	DB *sql.DB
}

func (m *RatingModel) Insert(movie, director, release_date, rating, review string) (int, error) {
	var id int

	s := `
		INSERT INTO movierating(movie_name, director_name, release_date, movie_rating, movie_review)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := m.DB.QueryRow(s, movie, director, release_date, rating, review).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *RatingModel) Read() ([]*models.Rating, error) {
	s := `
		SELECT movie_name, director_name, release_date, movie_rating, movie_review
		FROM movierating
		LIMIT 20
	`
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	// cleanup before we leave Read()
	defer rows.Close()

	ratings := []*models.Rating{}

	for rows.Next() {
		r := &models.Rating{}
		err = rows.Scan(&r.Movie_name, &r.Director_name, &r.Release_date, &r.Movie_rating, &r.Movie_review)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (m *RatingModel) Get(id int) (*models.Rating, error) {
	s := `
		SELECT movie_name, director_name, release_date, movie_rating, movie_review
		FROM movierating
		WHERE id = $1
	`
	r := &models.Rating{}
	err := m.DB.QueryRow(s, id).Scan(&r.Movie_name, &r.Director_name, &r.Release_date, &r.Movie_rating, &r.Movie_review)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return r, nil
}
