package main

import (
	"github.com/gin-gonic/gin"
	ctr "jp.study/m/v2/controllers"
)

func InitRoutes(router *gin.Engine, server *Server) {
	router.GET("/", func(c *gin.Context) {

	})

	kanjiRouter := router.Group("/radicals")
	kanjiRouter.GET("/", func(c *gin.Context) {
		ctr.GetRadicals(c, server.dbConn)
	})
}