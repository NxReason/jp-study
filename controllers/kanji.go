package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetKanji(ctx *gin.Context) {
	ctx.String(http.StatusOK, "kanji from controllers")
}