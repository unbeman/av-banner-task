package models

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Banner struct {
	Id int `json:"banner_id"`

	FeatureId int   `json:"feature_id"`
	TagIds    []int `json:"tag_ids"`

	Content  string `json:"content"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type Banners []*Banner

func (b Banners) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetBannerInput struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
}

func (i *GetBannerInput) FromURI(r *http.Request) error {
	tagParam := r.URL.Query().Get("tag_id")
	tagId, err := strconv.Atoi(tagParam)
	if err != nil {
		return err
	}
	i.TagId = tagId

	featureParam := r.URL.Query().Get("feature_id")
	featureId, err := strconv.Atoi(featureParam)
	if err != nil {
		return err
	}
	i.FeatureId = featureId
	return nil
}

type GetBannerOutput string

func (gbo GetBannerOutput) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetBannersInput struct {
	TagId     *int
	FeatureId *int
	Limit     *int
	Offset    *int
}

func (i *GetBannersInput) FromURI(r *http.Request) error {
	tagParam := r.URL.Query().Get("tag_id")
	if tagParam != "" {
		tagId, err := strconv.Atoi(tagParam)
		if err != nil {
			return err
		}
		i.TagId = &tagId
	}

	featureParam := r.URL.Query().Get("feature_id")
	if featureParam != "" {
		featureId, err := strconv.Atoi(featureParam)
		if err != nil {
			return err
		}
		i.FeatureId = &featureId
	}

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return err
		}
		i.Limit = &limit
	}

	offsetParam := r.URL.Query().Get("offset")
	if limitParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			return err
		}
		i.Offset = &offset
	}
	return nil
}

type CreateBannerInput struct {
	Id        int   `json:"-"`
	FeatureId int   `json:"feature_id"`
	TagIds    []int `json:"tag_ids"`

	Content  string `json:"content"`
	IsActive bool   `json:"is_active"`
}

func (b CreateBannerInput) Bind(r *http.Request) error {
	uniqueTags := make(map[int]struct{})
	for _, tag := range b.TagIds {
		if _, ok := uniqueTags[tag]; ok {
			return fmt.Errorf("tag_ids are not unique")
		}
		uniqueTags[tag] = struct{}{}
	}
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

func (b UpdateBannerInput) Bind(r *http.Request) error {
	uniqueTags := make(map[int]struct{})
	for _, tag := range *b.TagIds {
		if _, ok := uniqueTags[tag]; ok {
			return fmt.Errorf("tag_ids are not unique")
		}
		uniqueTags[tag] = struct{}{}
	}
	return nil
}
