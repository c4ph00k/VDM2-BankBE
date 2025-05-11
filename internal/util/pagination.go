package util

import (
	"strconv"

	"github.com/pkg/errors"
)

// PaginationParams stores pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// NewPaginationParams creates new pagination parameters with defaults
func NewPaginationParams(pageStr, limitStr string) (*PaginationParams, error) {
	params := &PaginationParams{
		Page:  1,
		Limit: 10,
	}

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid page parameter")
		}
		if page < 1 {
			return nil, errors.New("page parameter must be greater than 0")
		}
		params.Page = page
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid limit parameter")
		}
		if limit < 1 {
			return nil, errors.New("limit parameter must be greater than 0")
		}
		if limit > 100 {
			limit = 100 // Cap at 100 to prevent abuse
		}
		params.Limit = limit
	}

	return params, nil
}

// Offset calculates the offset for database queries
func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int `json:"current_page"`
		TotalPages  int `json:"total_pages"`
		TotalItems  int `json:"total_items"`
		PerPage     int `json:"per_page"`
	} `json:"pagination"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, params *PaginationParams, totalItems int) *PaginatedResponse {
	resp := &PaginatedResponse{
		Data: data,
	}

	resp.Pagination.CurrentPage = params.Page
	resp.Pagination.PerPage = params.Limit
	resp.Pagination.TotalItems = totalItems

	// Calculate total pages
	totalPages := totalItems / params.Limit
	if totalItems%params.Limit != 0 {
		totalPages++
	}
	resp.Pagination.TotalPages = totalPages

	return resp
}