package main

import (
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

// PostDetail represents the post's details
type PostDetail struct {
	Title       string
	Categoria   string
	Descripcion string
	LinkPost    string
}

// Post represents a blog's post
type Post struct {
	URLImage   string
	DateString string
	Detalle    PostDetail
}
type PostDAO struct {
	Server string
	Database string
}

var db *mgo.Database
const (
	category = "Entretenimiento"
	collection = "post"
)

func main() {
	dao := PostDAO {
		Server: "localhost:27017",
		Database: "posts-escuelita",
	}

	dao.connect()
	c := colly.NewCollector(
		colly.AllowedDomains("www.venusgo.com", "venusgo.com"),
		colly.Debugger(&debug.LogDebugger{}),
		colly.CacheDir("./venus_cache"),
	)

	detailsCollector := c.Clone()

	c.OnHTML("h4.titulo a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/blog") {
			return
		}

		detailsCollector.Visit("http://www.venusgo.com" + link)
	})

	detailsCollector.OnHTML("div.main-container", func(e *colly.HTMLElement) {
		title := e.DOM.Find("h1.page-header").Text()
		if title == "" {
			fmt.Println("No title found", e.Request.URL)
			return
		}

		description, err := e.DOM.Find("div.blog-descripcion").Html()
		if err != nil {
			fmt.Println("No Descripcion found", e.Request.URL)
			return
		}
		post := Post{
			Detalle: PostDetail{
				Title:       title,
				LinkPost:    e.Request.URL.String(),
				Descripcion: description,
				Categoria:   category,
			},
		}

		err = dao.insert(post)
		if err != nil {
			log.Fatal(err)
		}
	})

	c.Visit(os.Args[1])
	
}


func (p *PostDAO) connect() {
	session, err := mgo.Dial(p.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(p.Database)
}

func (p *PostDAO) insert(post Post) error{
	err := db.C(collection).Insert(post)
	return err
}