package data

import (
	"database/sql"
	"time"

	"github.com/Greenlight/internal/validator"
	"github.com/lib/pq"
)

//MovieModel wraps the sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

type Movie struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"-"`
	Version int32 `json:"version"`
	Runtime Runtime `json:"runtime,omitempty"`
	Genres []string `json:"genre,omitempty"`
	Year int32 `json:"year,omitempty"`
}

// Insert inserts a new movie record into the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	stmt := `
		INSERT INTO movies (title, year, runtime, genres)	
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get fetches a specific movie record with the id
func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// Update updates a record with the movie arg passed.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Delete deletes a specific movie record with the id
func (m MovieModel) Delete(id int64) error {
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater or equal to 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate genre")
}
