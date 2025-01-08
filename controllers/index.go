package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RadicalsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "radicals.html", nil)
}