/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE users (
  id serial PRIMARY KEY,
  full_name VARCHAR ( 60 ) NOT NULL,
  phone_number VARCHAR ( 13 ) UNIQUE NOT NULL,
  password VARCHAR ( 64 ) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_phone_number ON users(phone_number);

CREATE TABLE login_histories (
  id serial PRIMARY KEY,
  users_id serial NOT NULL,
  counter int NOT NULL DEFAULT 0
);

CREATE INDEX idx_login_histories_users_id ON login_histories(users_id);