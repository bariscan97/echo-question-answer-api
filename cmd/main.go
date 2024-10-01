package main

import (
	"log"
	"os"
	"articles-api/database"
	"articles-api/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	
	err := godotenv.Load("../.env")
	
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	e := echo.New()
	
	db := database.GetPool()
	
	r := routes.NewRouters(e, db)
	
	r.InitRouter()
	
	log.Fatal(r.RoutersRun(os.Getenv("PORT")))
}