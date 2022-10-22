package data

import "time"

type Movie struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"-"`
	Version int32 `json:"version"`
	Runtime Runtime `json:"runtime,omitempty"`
	Genre []string `json:"genre,omitempty"`
	Year int32 `json:"year,omitempty"`
}

