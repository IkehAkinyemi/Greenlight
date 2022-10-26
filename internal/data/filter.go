package data

import (
	"math"
	"strings"

	"github.com/lighten/internal/validator"
)

// Filter contains the parsed query_string
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// sortColumn checks that the client-provided Sort field matches one of
// the entries in the safelist and if it does, retrieves it.
func (f Filters) sortColumn() string {
	for _, safeVal := range f.SortSafelist {
		if f.Sort == safeVal {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	// failsafe to help stop a SQL injection attack occurring.
	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection intructs a descending or ascending 
// order for 'GET /v1/movies?<query_string>'
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters validates the query_string for abnormalities
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// Provides extra info about the filtered, sorted and paginated 
// info returned on 'GET /v1/movies?<query_string>'
type Metadata struct {
	CurrentPage int `json:"current_page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	FirstPage int `json:"first_page,omitempty"`
	LastPage int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calcMetadata calculates and return pagination info
func calcMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage: page,
		PageSize: pageSize,
		FirstPage: 1,
		LastPage: int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}