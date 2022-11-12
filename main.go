package main

import (
	"Golang-Scraper/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var urls = []string{
	"https://mathstat.um.ac.ir/index.php/fa/2016-12-19-08-40-28/2016-01-23-11-32-19/2019-05-21-05-38-11",
}

func main() {
	t := time.Now()
	var professors []models.Professors

	var wg sync.WaitGroup
	ch := make(chan models.Professors)

	wg.Add(len(urls))

	for _, url := range urls {
		go scrape(url, ch)
	}

	for range urls {
		go func() {
			defer wg.Done()
			professor := <-ch
			professors = append(professors, professor)
		}()
	}

	wg.Wait()

	result, err := json.Marshal(professors)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = os.WriteFile("result.json", result, 0644)
	if err != nil {
		return
	}

	file, err := os.Create("data.csv")

	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Name", "Email"})
	if err != nil {
		return
	}

	for _, professor := range professors {
		err := writer.Write([]string{
			professor.Name,
			professor.Email,
		})
		if err != nil {
			return
		}
	}

	close(ch)

	elapsed := time.Since(t).Seconds()
	fmt.Printf("Time: %.2fs\n", elapsed)
}

func scrape(url string, ch chan<- models.Professors) {
	var professor models.Professors
	c := colly.NewCollector()

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	c.OnHTML("#tab-0", func(e *colly.HTMLElement) {
		name := e.ChildText("a")
		price := e.ChildText("a")

		professor = models.Professors{
			Name:  name,
			Email: price,
		}
	})

	err := c.Visit(url)
	if err != nil {
		return
	}

	ch <- professor
}
