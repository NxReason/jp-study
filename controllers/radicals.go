package controllers

import (
	"log"
	"net/http"

	"golang.org/x/exp/maps"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"jp.study/m/v2/models"
)

func GetRadicals(c *gin.Context, db *pgx.Conn) {
	query := `
	SELECT r.id, r.glyph, rm.id as rm_id, rm.meaning as rm
	FROM radicals r
	LEFT JOIN radical_meanings rm ON r.id = rm.radical_id`

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

		if rm.ID != nil {
			radical.Meanings = append(radical.Meanings, rm)
		}
		if _, exists := radicalsMap[radical.ID]; !exists {
			radicalsMap[radical.ID] = radical
		} else {
			oldRadical := radicalsMap[radical.ID]
			radical.Meanings = append(oldRadical.Meanings, rm)
			radicalsMap[radical.ID] = radical
		}
	}

	c.JSON(http.StatusOK, gin.H{"radicals": maps.Values(radicalsMap)})
}

type NewRadicalBody struct {
	Glyph string `json:"glyph"`
	Meanings []string `json:"meanings"`
}

func SaveRadical(c *gin.Context, db *pgx.Conn) {
	var radicalJSON NewRadicalBody

	if err := c.ShouldBindJSON(&radicalJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save new radical
	var id int
	query := `INSERT INTO radicals (glyph) VALUES ($1) RETURNING id`
	err := db.QueryRow(c, query, radicalJSON.Glyph).Scan(&id)
	if err != nil { 
		c.JSON(http.StatusBadRequest, gin.H {"error": err.Error()})
		return
	}

	// save meanings
	query = `INSERT INTO radical_meanings (meaning, radical_id) VALUES ($1, $2)`
	for _, rm := range radicalJSON.Meanings {
		_, err = db.Exec(c, query, rm, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H {"error": err.Error()})
			return
		}
	}
	

	// get saved radical with its meanings
	query = `
	SELECT r.id, r.glyph, rm.id, rm.meaning 
	FROM radicals r
	LEFT JOIN radical_meanings rm
	ON r.id = rm.radical_id
	WHERE r.id = $1`
	rows, err := db.Query(c, query, id)
	if err != nil { return } // TODO: bad request

	var newRadical models.Radical
	var newRM models.RadicalMeaning
	var meanings []models.RadicalMeaning
	for rows.Next() {
		err = rows.Scan(nil, nil, &newRM.ID, &newRM.Meaning)
		if err != nil {
			log.Fatalf("Scanning rows failed %v\n", err)
		}
		meanings = append(meanings, newRM)
	}
	rows.Scan(&newRadical.ID, &newRadical.Glyph, nil, nil)
	newRadical.Meanings = meanings

	c.JSON(http.StatusOK, gin.H{"message": "Radical saved", "radical": newRadical})
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