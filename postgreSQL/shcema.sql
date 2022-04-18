

CREATE TABLE movies(
    id INTEGER,
    poster_path VARCHAR(255),
    adult BOOLEAN,
    overview TEXT,
    release_date VARCHAR(255),
    original_title VARCHAR(255),
    original_language VARCHAR(255),
    title VARCHAR(255),
    backdrop_path VARCHAR(255),
    popularity FLOAT,
    vote_count INTEGER,
    video BOOLEAN,
    vote_average FLOAT,
    genre_id integer[]
);

CREATE TABLE Users(
    id INTEGER,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    faviorate_movies integer[]
);

