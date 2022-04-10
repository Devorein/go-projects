package main

import (
	"fmt"
	"go-crm-fiber/database"
	"go-crm-fiber/lead"

	"github.com/gofiber/fiber"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.CreateLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error
	database.Connect()
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Connected to database")
	database.DBConn.AutoMigrate(&lead.Lead{})
}

func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(3000)
	defer database.DBConn.Close()
}
