/*
CMPS3162 - Final Project
Joanne Yong & Rene Sanchez
2019120152  & 2018118383
*/

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,       
    activated bool NOT NULL DEFAULT TRUE
);