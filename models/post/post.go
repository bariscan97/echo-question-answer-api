package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostModel struct {
	Title    string `json:"title" validate:"required,min=10,max=255"`
	Content  string `json:"content" validate:"required,min=10,max=255"`
	Category string `json:"category" validate:"omitempty,min=1,max=255"`
}

type UpdatePostModel struct {
	Title   string `json:"title" validate:"omitempty,min=10,max=255"`
	Content string `json:"content" validate:"omitempty,min=10,max=255"`
}

type FetchPostModel struct {
	Id          uuid.UUID
	Username    string
	Parent_id   *uuid.UUID
	Profile_img *string
	Like_count  int
	Title       string
	Content     string
	Created_at  time.Time
	Updated_at  time.Time
}
