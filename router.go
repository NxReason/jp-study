package main

import (
	"github.com/gin-gonic/gin"
	ctr "jp.study/m/v2/controllers"
)

func InitRoutes(router *gin.Engine, server *Server) {
	router.GET("/", func(c *gin.Context) {
		ctr.Index(c)
	})
	router.GET("/radicals", ctr.RadicalsPage)

	radicalRouter := router.Group("/api/radicals")
	radicalRouter.GET("/", ctr.GetRadicals(server.dbConn))
	radicalRouter.POST("/", ctr.SaveRadical(server.dbConn))
	radicalRouter.DELETE("/", ctr.DeleteRadical(server.dbConn))
}