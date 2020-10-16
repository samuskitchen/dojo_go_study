CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    name VARCHAR(250) NOT NULL,
    surname VARCHAR(250) NOT NULL,
    username VARCHAR(150) NOT NULL UNIQUE,
    password varchar(256) NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY(id)
);