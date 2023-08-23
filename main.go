package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Film struct {
	Title    string `form:"title" json:"title" binding:"required"`
	Director string `form:"director" json:"director" binding:"required"`
}

func main() {
	fmt.Println("Go app...")

	// gin router
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong from gin",
		})
	})

	// handler function #1 - returns the index.html template, with film data
	h1 := func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		tmpl.Execute(c.Writer, films)
	}

	// handler function #2 - returns the template block with the newly added film, as an HTMX response
	h2 := func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		var film Film
		if err := c.Bind(&film); err != nil {
			// TODO: this should not echo the error to the client
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		slog.Info("add-film", "film", film)

		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl := template.Must(template.ParseFiles("index.html"))
		// c.Render(http.StatusOK, func(w http.ResponseWriter) error {
		tmpl.ExecuteTemplate(c.Writer, "film-list-element", film)
		// })
		// 	c.Render(http.StatusOK, func(w http.ResponseWriter) error {
		// 		return tmpl.Execute(w, film)
		// return tmpl.ExecuteTemplate(w, "film-list-element", film)
		// 	})

		c.Render(http.StatusOK, func(w http.ResponseWriter) error {})
	}

	halpine := func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("halpine.html"))
		tmpl.Execute(c.Writer, nil)
	}

	// define handlers
	r.GET("/", h1)
	r.POST("/add-film/", h2)
	r.GET("/halpine/", halpine)

	r.Run() // listen and serve on
}
