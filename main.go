package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("Go app...")

	// use echo server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, World!")
	})
	e.GET("/echo/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Echo, World!")
	})

	// handler function #1 - returns the index.html template, with film data
	h1 := func(c echo.Context) error {
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		return tmpl.Execute(c.Response().Writer, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	h2 := func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		title := c.Request().PostFormValue("title")
		director := c.Request().PostFormValue("director")
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl := template.Must(template.ParseFiles("index.html"))
		return tmpl.ExecuteTemplate(c.Response().Writer, "film-list-element", Film{Title: title, Director: director})
	}

	halpine := func(c echo.Context) error {
		tmpl := template.Must(template.ParseFiles("halpine.html"))
		return tmpl.Execute(c.Response().Writer, nil)
	}

	// define handlers
	e.GET("/", h1)
	e.POST("/add-film/", h2)
	e.GET("/halpine", halpine)

	e.Logger.Fatal(e.Start(":1323"))

}
