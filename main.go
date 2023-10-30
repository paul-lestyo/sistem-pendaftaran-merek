package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/router"
)

var Validate = validator.New()
var Store = session.New()

func main() {
	database.Connect()

	//metrics
	//NewHamming(),NewLevenshtein(), NewJaro(),NewJaroWinkler(), etc
	//similarity := strutil.Similarity("stackoverflow", "stackoverflw", metrics.NewLevenshtein())
	//fmt.Printf("hoho:%.2f\n", similarity) // Output: 0.75

	seedRole()

	engine := html.New("./views", ".gohtml")
	engine.Reload(true)
	engine.Debug(true)
	engine.AddFunc("increment", func(num int) int {
		return num + 1
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/", "./assets")
	router.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":8080")
}

var ResultIDRole struct {
	ID uuid.UUID
}

func seedRole() {
	err := database.DB.Table("roles").Select("id").Where("name = ?", "Pemohon").First(&ResultIDRole).Error
	if err != nil {
		database.DB.Create(&model.Role{Name: "Admin"})
		database.DB.Create(&model.Role{Name: "Pemohon"})
	}
}
