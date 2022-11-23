package api

import (
	v1 "github.com/SaidovZohid/auth-redis-verify/api/v1"
	"github.com/SaidovZohid/auth-redis-verify/config"
	"github.com/SaidovZohid/auth-redis-verify/storage"
	"github.com/gin-gonic/gin"
)

type RouteOption struct {
	Cfg *config.Config
	Storage storage.StorageI
}

func New(opt *RouteOption) *gin.Engine {
	router := gin.Default()

	handler := v1.New(&v1.Handler{
		Cfg: opt.Cfg,
		Storage: &opt.Storage,
	})
	apiV1 := router.Group("/user")
	{
		apiV1.POST("/register", handler.Register)
		apiV1.POST("/veirfy", handler.Verify)
	}
	return router
}