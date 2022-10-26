package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/lighten/internal/validator"
)

// MovieModel wraps the sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genre,omitempty"`
	Year      int32     `json:"year,omitempty"`
}

// Insert inserts a new movie record into the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	stmt := `
		INSERT INTO movies (title, year, runtime, genres)	
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get fetches a specific movie record with the id
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `
		SELECT id, title, created_at, version, runtime, genres, year 
		FROM movies
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var movie Movie

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.CreatedAt,
		&movie.Version,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Year,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

// GetAll gets all the movie record that's matched by the query_string.
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	stmt := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version 
		FROM movies 
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC 
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection(),
	)
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	if rows.Err() != nil {
		return nil, Metadata{}, err
	}
	metadata := calcMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

// Update updates a record with the movie arg passed.
func (m MovieModel) Update(movie *Movie) error {
	stmt := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete deletes a specific movie record with the id
func (m MovieModel) Delete(id int64) error {
	stmt := `DELETE FROM movies WHERE id = $1`
	if id < 1 {
		return ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rows, err := resp.RowsAffected()
	if rows == 0 {
		return ErrRecordNotFound
	}

	return err
}

// ValidateMovie sanity-checks the movie JSON values provided.
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
