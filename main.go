package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	env, err := Env()
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", env["PG_USER"], env["PG_PASS"], env["PG_DB"])
	port, ok := env["PORT"]
	if !ok {
		port = "8080"
	}
	
	server := Server{}
	server.DbConnect(connString)
	server.SetRouter()
	server.Start(port)
}

type Server struct {
	dbConn *pgx.Conn
	router *gin.Engine
}

func (s *Server) DbConnect(connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	s.dbConn = conn
}

func (s *Server) SetRouter() {
	router := gin.Default()

	router.LoadHTMLGlob("./views/*")
	// static
	router.Static("/public", "./public")

	InitRoutes(router, s)

	s.router = router
}

func (s *Server) Start(port string) {
	s.router.Run(":" + port)
}