package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func GetDataChartPermohonanMerek(c *fiber.Ctx) error {
	date := c.Params("date")
	formatDate := ""
	if date == "day" {
		formatDate = "'%Y-%m-%d'"
	} else if date == "month" {
		formatDate = "'%Y-%m'"
	} else if date == "year" {
		formatDate = "'%Y'"
	}

	var chartPermohonanMerek [][]interface{}
	rows, _ := database.DB.
		Model(&model.Brand{}).
		Select("DATE_FORMAT(created_at, " + formatDate + ") as date, count(*) as count").
		Group("date").
		Rows()
	defer rows.Close()

	for rows.Next() {
		var date string
		var count int
		rows.Scan(&date, &count)

		chartPermohonanMerek = append(chartPermohonanMerek, []interface{}{date, count})
	}

	return c.JSON(chartPermohonanMerek)
}

func GetDataChartLogin(c *fiber.Ctx) error {
	date := c.Params("date")
	formatDate := ""
	if date == "day" {
		formatDate = "'%Y-%m-%d'"
	} else if date == "month" {
		formatDate = "'%Y-%m'"
	} else if date == "year" {
		formatDate = "'%Y'"
	}

	var chartLogin [][]interface{}
	rowsLog, _ := database.DB.
		Model(&model.Log{}).
		Select("DATE_FORMAT(created_at, " + formatDate + ") as date, count(*) as count").
		Group("date").
		Rows()
	defer rowsLog.Close()

	for rowsLog.Next() {
		var date string
		var count int
		rowsLog.Scan(&date, &count)

		chartLogin = append(chartLogin, []interface{}{date, count})
	}

	return c.JSON(chartLogin)
}
