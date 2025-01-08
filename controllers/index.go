package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", nil)
}