package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoUrl = "mongodb://localhost:27017/" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

type ApiResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employees", func(ctx *fiber.Ctx) error {
		var employees []Employee

		query := bson.D{{}}

		cursor, err := mg.Db.Collection("employees").Find(ctx.Context(), query)
		if err != nil {
			return ctx.Status(500).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		if err := cursor.All(ctx.Context(), &employees); err != nil {
			return ctx.Status(500).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		return ctx.JSON(ApiResponse{
			Error:   false,
			Message: "",
			Data:    employees,
		})
	})

	app.Get("/employee/:id", func(ctx *fiber.Ctx) error {
		var employee Employee

		employeeId, err := primitive.ObjectIDFromHex(ctx.Params("id"))

		if err != nil {
			return ctx.Status(400).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		singleResult := mg.Db.Collection("employees").FindOne(ctx.Context(), bson.D{{
			Key:   "_id",
			Value: employeeId,
		}})
		if singleResult.Err() != nil {
			return ctx.Status(500).JSON(ApiResponse{
				Error:   true,
				Message: singleResult.Err().Error(),
				Data:    nil,
			})
		}

		singleResult.Decode(&employee)

		return ctx.JSON(ApiResponse{
			Error:   true,
			Message: "",
			Data:    employee,
		})
	})

	app.Post("/employee", func(ctx *fiber.Ctx) error {
		var employee Employee

		err := ctx.BodyParser(&employee)

		if err != nil {
			return ctx.Status(400).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
			})
		}

		employee.ID = ""

		insertionResult, err := mg.Db.Collection("employees").InsertOne(ctx.Context(), &employee)
		if err != nil {
			return ctx.Status(500).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
			})
		}

		employee.ID = insertionResult.InsertedID.(primitive.ObjectID).String()
		return ctx.JSON(ApiResponse{
			Error:   false,
			Message: "",
			Data:    employee,
		})
	})

	app.Put("/employee/:id", func(ctx *fiber.Ctx) error {
		var employee Employee

		employeeCollection := mg.Db.Collection("employees")
		if err := ctx.BodyParser(&employee); err != nil {
			return ctx.Status(400).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
			})
		}
		idParam := ctx.Params("id")
		hexIdParam, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return ctx.Status(400).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
			})
		}

		findQuery := bson.D{{Key: "_id", Value: hexIdParam}}
		updateQuery := bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "name",
				Value: employee.Name,
			}, {
				Key:   "age",
				Value: employee.Age,
			}, {
				Key:   "salary",
				Value: employee.Salary,
			}},
		}}

		singleResult := employeeCollection.FindOneAndUpdate(ctx.Context(), findQuery, updateQuery)
		err = singleResult.Err()

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return ctx.JSON(ApiResponse{
					Error:   true,
					Message: err.Error(),
					Data:    nil,
				})
			}
			return ctx.JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		employee.ID = idParam

		return ctx.JSON(ApiResponse{
			Error:   false,
			Message: "",
			Data:    employee,
		})
	})

	app.Delete("/employee/:id", func(ctx *fiber.Ctx) error {
		idParam := ctx.Params("id")
		hexIdParam, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return ctx.JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		_, err = mg.Db.Collection("employees").DeleteOne(ctx.Context(), bson.D{{
			Key:   "_id",
			Value: hexIdParam,
		}})

		if err != nil {
			return ctx.Status(400).JSON(ApiResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
		}

		return ctx.JSON(ApiResponse{
			Error:   false,
			Message: "",
			Data:    true,
		})
	})

	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err.Error())
	}
}
