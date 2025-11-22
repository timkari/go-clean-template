// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

// Comment - represents a comment entity.
type Comment struct {
	ID         int64     `json:"id"          example:"1"`
	EntityType string    `json:"entity_type" example:"order"`
	EntityID   int64     `json:"entity_id"   example:"123"`
	UserID     int64     `json:"user_id"     example:"456"`
	UserName   string    `json:"user_name"   example:"John Doe"`
	Content    string    `json:"content"     example:"This is a comment"`
	CreatedAt  time.Time `json:"created_at"  example:"2024-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at"  example:"2024-01-01T00:00:00Z"`
}

// CommentList - represents a paginated list of comments.
type CommentList struct {
	Comments []Comment `json:"comments"`
	Total    int64     `json:"total"     example:"100"`
	Page     int       `json:"page"      example:"1"`
	PageSize int       `json:"page_size" example:"10"`
}

// SortDirection - represents sorting direction.
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// CommentFilter - represents filters for comments query.
type CommentFilter struct {
	EntityType string        `json:"entity_type"`
	EntityID   int64         `json:"entity_id"`
	Page       int           `json:"page"        example:"1"`
	PageSize   int           `json:"page_size"   example:"10"`
	SortBy     string        `json:"sort_by"     example:"created_at"`
	SortOrder  SortDirection `json:"sort_order"  example:"desc"`
}
