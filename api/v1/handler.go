package v1

import (
	"github.com/SaidovZohid/auth-redis-verify/api/models"
	"github.com/SaidovZohid/auth-redis-verify/config"
	"github.com/SaidovZohid/auth-redis-verify/storage"
)

type handler struct {
	cfg *config.Config
	storage storage.StorageI
}

type Handler struct {
	Cfg *config.Config
	Storage *storage.StorageI
}

func New(options *Handler) *handler {
	return &handler{
		cfg: options.Cfg,
		storage: *options.Storage,
	}
}

func errResponse(err error) models.ResponseError {
	return models.ResponseError{
		Error: "Time is out register again",
	}
}