package models

import "time"

type PaginationParams struct {
	Cursor time.Time `query:"cursor" required:"false"`
}

type PaginationResponseMetadata struct {
	HasNext        bool       `json:"has_next"`
	PreviousCursor *time.Time `json:"previous,omitempty"`
	NextCursor     *time.Time `json:"next,omitempty"`
}
