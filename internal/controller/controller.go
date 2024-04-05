package controller

import "github.com/unbeman/av-banner-task/internal/storage"

type Controller struct {
	database storage.Database
}

func NewController(db storage.Database) (*Controller, error) {
	ctrl := &Controller{database: db}
	return ctrl, nil
}
