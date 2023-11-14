package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/router"
	"math/rand"
	"time"
)

func main() {
	database.Connect()

	//metrics
	//NewHamming(),NewLevenshtein(), NewJaro(),NewJaroWinkler(), etc
	//similarity := strutil.Similarity("stackoverflow", "stackoverflw", metrics.NewLevenshtein())
	//fmt.Printf("hoho:%.2f\n", similarity) // Output: 0.75

	seedRole()
	seedLog()

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
	var testKey = encryptcookie.GenerateKey()
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: testKey,
	}))
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

var IDBusiness struct {
	ID uuid.UUID
}

type IDUser struct {
	ID uuid.UUID
}

func seedRole() {
	err := database.DB.Table("roles").Select("id").Where("name = ?", "Pemohon").First(&ResultIDRole).Error
	if err != nil {
		admin := model.Role{Name: "Admin"}
		pemohon := model.Role{Name: "Pemohon"}
		database.DB.Create(&admin)
		database.DB.Create(&pemohon)

		hashedPassword, _ := helper.HashPassword("123")
		database.DB.Create(&model.User{
			Name:     "Paulus Lestyo A",
			Email:    "paulus.lestyo@student.uns.ac.id",
			Password: hashedPassword,
			RoleID:   admin.ID,
			IsActive: true,
		})
		database.DB.Create(&model.User{
			Name:     "Paul L A",
			Email:    "lestyo24@gmail.com",
			Password: hashedPassword,
			RoleID:   pemohon.ID,
			Business: &model.Business{},
			IsActive: true,
		})

	}

	for i := 0; i <= rand.Intn(5); i++ {
		hashedPassword, _ := helper.HashPassword("123")
		database.DB.Create(&model.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: hashedPassword,
			RoleID:   ResultIDRole.ID,
			Business: &model.Business{
				BusinessName:      faker.Name(),
				BusinessAddress:   faker.Name(),
				BusinessLogo:      "",
				OwnerName:         faker.Name(),
				UMKCertificateUrl: "",
				SignatureUrl:      "",
			},
			IsActive: true,
		})

		values := []string{"OK", "Perbaiki", "Tolak", "Menunggu"}
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(values))

		var idUser IDUser
		database.DB.Model(&model.User{}).Select("id").Order("RAND()").First(&idUser)
		database.DB.Model(&model.Business{}).Select("id").Order("RAND()").First(&IDBusiness)
		createdAt := time.Date(time.Now().Year(), time.Month(i), rand.Intn(30), 0, 0, 0, 0, time.UTC)
		database.DB.Create(&model.Brand{
			BusinessID:  IDBusiness.ID,
			BrandName:   faker.Name(),
			DescBrand:   faker.Word(),
			BrandLogo:   "",
			Status:      values[randomIndex],
			Note:        faker.Sentence(),
			CreatedByID: idUser.ID,
			UpdatedByID: idUser.ID,
			CreatedAt:   &createdAt,
		})
	}

}

func seedLog() {
	for month := 1; month <= 12; month++ {
		for i := 1; i <= rand.Intn(3); i++ {
			var idUser IDUser
			database.DB.Model(&model.User{}).Select("id").Order("RAND()").First(&idUser)
			createdAt := time.Date(time.Now().Year(), time.Month(month), rand.Intn(30), 0, 0, 0, 0, time.UTC)
			database.DB.Create(&model.Log{
				UserID:    idUser.ID,
				CreatedAt: &createdAt,
			})
		}
	}
}
