package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/jasonzbao/api-template/api/apimodels"
	"github.com/jasonzbao/api-template/api/utils/ginutils"
)

func (s *Server) addKvsRoutes(router *gin.Engine) {
	router.GET("/kvs", s.handleGetKvs)
	router.POST("/kvs", s.handlePostKvs)
}

func (s *Server) handleGetKvs(c *gin.Context) {
	resp := &apimodels.V1KvsResponse{}
	defer resp.FormatReturn(c, resp)

	req := &apimodels.V1KvsGetRequest{}
	if err := ginutils.ShouldBindWith(c, req, binding.Query); err != nil {
		return
	}

	fmt.Println("req.Key", req.Key)

	value, err := s.dbClient.Get(req.Key)
	if err != nil {
		resp.Error = err
		return
	}

	resp.Value = value
}

func (s *Server) handlePostKvs(c *gin.Context) {
	resp := &apimodels.V1KvsResponse{}
	defer resp.FormatReturn(c, resp)

	req := &apimodels.V1KVSPostRequest{}
	if err := ginutils.ShouldBindWith(c, req, binding.JSON); err != nil {
		return
	}

	err := s.dbClient.Set(req.Key, req.Value)
	if err != nil {
		resp.Error = err
		return
	}

	resp.Value = req.Value
}
