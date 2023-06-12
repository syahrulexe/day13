package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"syahrul/connection"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Id           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Description  string
	postingTime  string
	Html         bool
	Css          bool
	Javascript   bool
	Java         bool
	Technologies []string
	Image        string
}

var dataProject = []Project{
	// {
	// 	ProjectName: "ex 1",
	// 	StartDate:   ,
	// 	EndDate:     ,
	// 	Duration:    "3 Months",
	// 	Html:        false,
	// 	Css:         true,
	// 	Javascript:  true,
	// 	Java:        false,
	// 	postingTime: "23",
	// 	Description: "lorem lorem",
	// },
	// {
	// 	ProjectName: "ex 2",
	// 	StartDate:   "11-01",
	// 	EndDate:     "22-01",
	// 	Duration:    "3 Months",
	// 	Html:        true,
	// 	Css:         true,
	// 	Javascript:  false,
	// 	Java:        true,
	// 	postingTime: "23",
	// 	Description: "blablabalala",
	// },
}

func main() {
	connection.DatabaseConnect()
	e := echo.New()

	e.Static("/public", "public")
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/detailproject/:id", detailproject)
	e.GET("/myproject", myProject)
	e.GET("/testimoni", testimoni)
	e.POST("/add-Project", addProject)
	e.POST("/deleteProject/:id", deleteProject)
	e.POST("/edit-project/:id", ressEditProject)
	e.GET("/edit-project/:id", editProject)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date,duration, description, html, css, javascript, java FROM tb_projects")

	var ress []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Html, &each.Css, &each.Javascript, &each.Java)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}
		ress = append(ress, each)
	}

	projects := map[string]interface{}{
		"Projects": ress,
	}

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func testimoni(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimoni.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func myProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/myproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var DetailProject = Project{}

	err := connection.Conn.QueryRow(context.Background(),
		"SELECT id, name, start_date, end_date,duration, description, technologies, html, css, javascript, java FROM tb_projects WHERE id=$1", id).Scan(
		&DetailProject.Id, &DetailProject.ProjectName, &DetailProject.StartDate, &DetailProject.EndDate, &DetailProject.Duration, &DetailProject.Description, &DetailProject.Technologies, &DetailProject.Html, &DetailProject.Css, &DetailProject.Javascript, &DetailProject.Java)

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	data := map[string]interface{}{
		"Project": DetailProject,
	}

	var tmpl, errTemplate = template.ParseFiles("views/detailproject.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	projectName := c.FormValue("input-name")
	startDate := c.FormValue("input-start")
	endDate := c.FormValue("input-end")
	description := c.FormValue("input-description")
	html := c.FormValue("input-check-html")
	css := c.FormValue("input-check-css")
	javascript := c.FormValue("input-check-javascript")
	java := c.FormValue("input-check-java")
	// konversi value cekbox, string to boolean
	htmlValue := html != ""
	cssValue := css != ""
	javascriptValue := javascript != ""
	javaValue := java != ""
	// parsing string to time.Time
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects (name, start_date, end_date, description, duration, html, css, javascript, java) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		projectName, start, end, description, getDuration(startDate, endDate), htmlValue, cssValue, javascriptValue, javaValue)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, duration, html, css, javascript, java FROM tb_projects WHERE id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Duration, &ProjectDetail.Html, &ProjectDetail.Css, &ProjectDetail.Javascript, &ProjectDetail.Java)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTemplate = template.ParseFiles("views/edit-project.html")
	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func ressEditProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Id :", id)

	projectName := c.FormValue("input-name")
	startDate := c.FormValue("input-start")
	endDate := c.FormValue("input-end")
	description := c.FormValue("input-description")
	html := c.FormValue("input-check-html")
	css := c.FormValue("input-check-css")
	javascript := c.FormValue("input-check-javascript")
	java := c.FormValue("input-check-java")
	// postingTime := time.Now()

	// konversi cekbox string to boolean
	htmlValue := html != ""
	cssValue := css != ""
	javascriptValue := javascript != ""
	javaValue := java != ""
	// parsing string to time.Time
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	_, err := connection.Conn.Exec(
		context.Background(), "UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, duration=$5, html=$6, css=$7, javascript=$8, java=$9 WHERE id=$10",
		projectName, start, end, description, getDuration(startDate, endDate), htmlValue, cssValue, javascriptValue, javaValue, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	fmt.Println("edit :", projectName)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func getDuration(startDate, endDate string) string {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	durationTime := int(endTime.Sub(startTime).Hours())
	durationDays := durationTime / 24
	durationWeeks := durationDays / 7
	durationMonths := durationWeeks / 4
	durationYears := durationMonths / 12

	var duration string

	if durationYears > 1 {
		duration = strconv.Itoa(durationYears) + " years"
	} else if durationYears == 1 {
		duration = strconv.Itoa(durationYears) + " year"
	} else if durationMonths > 1 {
		duration = strconv.Itoa(durationMonths) + " months"
	} else if durationMonths == 1 {
		duration = strconv.Itoa(durationMonths) + " month"
	} else if durationWeeks > 1 {
		duration = strconv.Itoa(durationWeeks) + " weeks"
	} else if durationWeeks == 1 {
		duration = strconv.Itoa(durationWeeks) + " week"
	} else if durationDays > 1 {
		duration = strconv.Itoa(durationDays) + " days"
	} else {
		duration = strconv.Itoa(durationDays) + " day"
	}

	return duration
}
