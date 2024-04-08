package models

import (
	"net/http"
	"time"
)

type Banner struct {
	Id int `json:"id"`

	FeatureId int   `json:"feature_id"`
	TagIds    []int `json:"tags"`

	Content  []byte `json:"content"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time
	DeletedAt time.Time
}

type GetBannerInput struct {
	TagId           int  `json:"tag_id"`
	FeatureId       int  `json:"feature_id"`
	UseLastRevision bool `json:"use_last_revision,omitempty"`
}

func (mr *GetBannerInput) Bind(r *http.Request) error {
	return nil
}

type GetBannerOutput struct {
}

func (g GetBannerOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
