package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"jp.study/m/v2/models"
)

func GetRadicals(c *gin.Context, db *pgx.Conn) {
	query := `
	SELECT r.id, r.glyph, rm.id as rm_id, rm.meaning as rm
	FROM radicals r
	INNER JOIN radical_meanings rm ON r.id = rm.radical_id`

	rows, err := db.Query(c, query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	radicalsMap := make(map[int]models.Radical)
	for rows.Next() {
		var radical models.Radical
		var rm models.RadicalMeaning

		err := rows.Scan(&radical.ID, &radical.Glyph, &rm.ID, &rm.Meaning)
		if err != nil {
			log.Fatalf("Scanning rows failed %v\n", err)
		}

		// TODO: rewrite: update meanings slice and reassign value once
		radical.Meanings = append(radical.Meanings, rm)
		if _, exists := radicalsMap[radical.ID]; !exists {
			radicalsMap[radical.ID] = radical
		} else {
			oldRadical := radicalsMap[radical.ID]
			radical.Meanings = append(oldRadical.Meanings, rm)
			radicalsMap[radical.ID] = radical
		}
	}

	c.JSON(http.StatusOK, gin.H{"radicals": radicalsMap})
}