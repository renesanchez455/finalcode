/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

DROP TABLE IF EXISTS movierating;

CREATE TABLE movierating (
    id serial PRIMARY KEY,
    movie_name text NOT NULL,
    director_name text NOT NULL,
    release_date date NOT NULL,
    movie_rating int NOT NULL,
    movie_review text NOT NULL
);
