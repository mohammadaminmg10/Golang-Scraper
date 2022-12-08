package test

import (
	"Golang-Scraper/models"
	"fmt"
	"github.com/gocolly/colly"
	"testing"
)

func TestScrap(t *testing.T) {
	var professor models.Professors
	c := colly.NewCollector()
	c.OnHTML("#tab-0", func(e *colly.HTMLElement) {
		name := e.ChildText("a")
		email := e.ChildText("a")

		professor = models.Professors{
			Name:  name,
			Email: email,
		}
	})
	url := "https://mathstat.um.ac.ir/index.php/fa/2016-12-19-08-40-28/2016-01-23-11-32-19/2019-05-21-05-38-11"
	err := c.Visit(url)
	if err != nil {
		return
	}
	ch := make(chan models.Professors)
	fmt.Println(ch, professor)
}
