package server

import (
	"github.com/gin-gonic/gin"
	apimodels "github.com/jasonzbao/api-template/api/apimodels"
)

func (s *Server) addPingRoutes(router *gin.Engine) {
	router.Any("/ping", s.handlePing)
	router.Any("/public/ping", s.handlePing)
}

func (s *Server) handlePing(c *gin.Context) {
	resp := &apimodels.V1PingResponse{}
	defer resp.FormatReturn(c, resp)

	resp.Version = s.config.Version
}
