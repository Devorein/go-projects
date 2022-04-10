package lead

import (
	"go-crm-fiber/database"

	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   int    `json:"phone"`
}

func GetLeads(ctx *fiber.Ctx) {
	db := database.DBConn
	var leads []Lead
	db.Find(&leads)
	ctx.JSON(leads)
}

func GetLead(ctx *fiber.Ctx) {
	id := ctx.Params("id")
	db := database.DBConn
	var lead Lead
	db.Find(&lead, id)
	ctx.JSON(lead)
}

func CreateLead(ctx *fiber.Ctx) {
	db := database.DBConn
	lead := new(Lead)
	if err := ctx.BodyParser(lead); err != nil {
		ctx.Status(503).Send(err)
		return
	}

	db.Create(&lead)
	ctx.JSON(lead)
}

func DeleteLead(ctx *fiber.Ctx) {
	db := database.DBConn
	id := ctx.Params("id")
	var lead Lead
	db.First(&lead, id)
	if lead.Name == "" {
		ctx.Status(400).Send("No lead found with ID")
		return
	}
	db.Delete(&lead)
	ctx.Send("Lead successfully deleted")
}
