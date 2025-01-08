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
	radicalRouter.GET("/", func(c *gin.Context) {
		ctr.GetRadicals(c, server.dbConn)
	})
	radicalRouter.POST("/", func(c *gin.Context) {
		ctr.SaveRadical(c, server.dbConn)
	})
	radicalRouter.DELETE("/", func(c *gin.Context) {
		ctr.DeleteRadical(c, server.dbConn)
	})
}