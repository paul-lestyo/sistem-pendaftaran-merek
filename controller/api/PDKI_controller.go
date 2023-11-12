package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
)

func SearchDataPDKI(c *fiber.Ctx) error {
	search := c.Params("search")
	return c.JSON(helper.GetDataSearchPDKI(search))
}
