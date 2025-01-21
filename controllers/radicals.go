package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"jp.study/m/v2/models"
	vm "jp.study/m/v2/viewmodels"
)


func GetRadicals(conn *pgx.Conn) gin.HandlerFunc {
	rt := models.RadicalTable{ Conn: conn }
	return func (c *gin.Context) {
		radicals := rt.All(c)
		radicalsView := vm.RadicalList(radicals)
		c.JSON(http.StatusOK, gin.H{"radicals": radicalsView })
	}
}

type NewRadicalBody struct {
	Glyph string `json:"glyph"`
	Meanings []string `json:"meanings"`
}

func SaveRadical(conn *pgx.Conn) gin.HandlerFunc {
	rt := models.RadicalTable{ Conn: conn }

	return func(c *gin.Context) {

		var radicalJSON NewRadicalBody

		if err := c.ShouldBindJSON(&radicalJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newRadical, err := rt.Save(c, radicalJSON.Glyph, radicalJSON.Meanings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error. Couldn't save new radical"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Radical saved", "radical": newRadical})
	}
}

type DeleteRadicalBody struct {
	ID int
}

func DeleteRadical(c *gin.Context, db *pgx.Conn) {
	var jsonBody DeleteRadicalBody
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `DELETE FROM radicals WHERE id = $1`
	tag, err := db.Exec(c, query, jsonBody.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tag.RowsAffected() != 1 {
		c.JSON(http.StatusOK, gin.H { "message": "No radical with this id" })
		return
	}

	c.JSON(http.StatusOK, gin.H { "message": "Radical removed successfully" })
}