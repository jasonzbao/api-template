package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	apiErrors "github.com/jasonzbao/api-template/api/errors"
	"github.com/jasonzbao/api-template/api/response"
)

func (s *Server) NewRouter() *gin.Engine {
	if s.config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.AllowedOrigins,
		AllowMethods:     []string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	router.Use(response.Middleware)

	{
		s.addPingRoutes(router)
		s.addKvsRoutes(router)
	}

	router.NoRoute(func(c *gin.Context) {
		resp := &response.V1ResponseBase{}
		resp.Error = apiErrors.ErrorNotFound
		resp.FormatReturn(c, resp)
	})

	return router
}
