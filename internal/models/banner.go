package models

import "time"

type Banner struct {
	Id int64 `json:"id"`

	FeatureId int64   `json:"feature_id"`
	TagIds    []int64 `json:"tags"`

	Content  []byte `json:"content"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time
	DeletedAt time.Time
}
