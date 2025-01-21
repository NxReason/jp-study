package models

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Radical struct {
	ID       int
	Glyph    string
	Meanings []RadicalMeaning
}

type RadicalMeaning struct {
	ID      *int
	Meaning *string
}

type RadicalTable struct {
	Conn *pgx.Conn
}

func (rt *RadicalTable) All(c *gin.Context) map[int]Radical {

	query := `
	SELECT r.id, r.glyph, rm.id as rm_id, rm.meaning as rm
	FROM radicals r
	LEFT JOIN radical_meanings rm ON r.id = rm.radical_id`

	rows, err := rt.Conn.Query(c, query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil
	}
	defer rows.Close()

	radicalsMap := make(map[int]Radical)
	for rows.Next() {
		var radical Radical
		var rm RadicalMeaning

		err := rows.Scan(&radical.ID, &radical.Glyph, &rm.ID, &rm.Meaning)
		if err != nil {
			log.Fatalf("Scanning rows failed %v\n", err)
			return nil
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

	return radicalsMap
}

// TODO: rewrite with RETURNING all new data
func (rt *RadicalTable) Save(c *gin.Context, glyph string, meanings []string) (Radical, error) {
	// save new radical
	var id int
	query := `INSERT INTO radicals (glyph) VALUES ($1) RETURNING id`
	err := rt.Conn.QueryRow(c, query, glyph).Scan(&id)
	if err != nil { 
		return Radical{}, err
	}

	// save meanings
	query = `INSERT INTO radical_meanings (meaning, radical_id) VALUES ($1, $2)`
	for _, rm := range meanings {
		_, err = rt.Conn.Exec(c, query, rm, id)
		if err != nil {
			return Radical{}, err
		}
	}
	
	// get saved radical with its meanings
	query = `
	SELECT r.id, r.glyph, rm.id, rm.meaning 
	FROM radicals r
	LEFT JOIN radical_meanings rm
	ON r.id = rm.radical_id
	WHERE r.id = $1`
	rows, err := rt.Conn.Query(c, query, id)
	if err != nil { 
		return Radical{}, err
	}

	var newRadical Radical
	var newRM RadicalMeaning
	var newMeanings []RadicalMeaning
	for rows.Next() {
		err = rows.Scan(nil, nil, &newRM.ID, &newRM.Meaning)
		if err != nil {
			log.Fatalf("Scanning rows failed %v\n", err)
		}
		newMeanings = append(newMeanings, newRM)
	}
	rows.Scan(&newRadical.ID, &newRadical.Glyph, nil, nil)
	newRadical.Meanings = newMeanings

	return newRadical, nil
}

func (rt *RadicalTable) Delete(c *gin.Context, id int) error {
	query := `DELETE FROM radicals WHERE id = $1`
	tag, err := rt.Conn.Exec(c, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return errors.New("No radical with this id")
	}

	return nil
}
