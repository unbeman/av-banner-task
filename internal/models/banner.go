package models

import (
	"fmt"
	"net/http"
	"time"
)

type Banner struct {
	Id int `json:"id"`

	FeatureId int   `json:"feature_id"`
	TagIds    []int `json:"tag_ids"`

	Content  string `json:"content"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

func (b Banner) Bind(r *http.Request) error {
	uniqueTags := make(map[int]struct{})
	for _, tag := range b.TagIds {
		if _, ok := uniqueTags[tag]; ok {
			return fmt.Errorf("tag_ids are not unique")
		}
		uniqueTags[tag] = struct{}{}
	}
	return nil
}

type Banners []*Banner

func (b Banners) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetBannerInput struct {
	TagId           int  `json:"tag_id"`
	FeatureId       int  `json:"feature_id"`
	UseLastRevision bool `json:"use_last_revision,omitempty"`
}

func (gbi *GetBannerInput) Bind(r *http.Request) error {
	return nil
}

type GetBannerOutput string

func (gbo GetBannerOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetBannersInput struct {
	TagId     *int `json:"tag_id,omitempty"`
	FeatureId *int `json:"feature_id,omitempty"`
	Limit     *int `json:"limit,omitempty"`
	Offset    *int `json:"offset,omitempty"`
}

func (gbi GetBannersInput) Bind(r *http.Request) error {
	return nil
}

type CreateBannerOutput struct {
	BannerId int `json:"banner_id"`
}

func (cbo CreateBannerOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UpdateBannerInput struct {
	Id int `json:"-"`

	FeatureId *int   `json:"feature_id,omitempty"`
	TagIds    *[]int `json:"tag_ids,omitempty"`

	Content  *string `json:"content,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

func (ubi UpdateBannerInput) Bind(r *http.Request) error {
	return nil
}
