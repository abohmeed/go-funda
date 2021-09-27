package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
)

type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	Rating      string
}

func main() {
	fName := "houses.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("funda.nl","www.funda.nl"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./coursera_cache"),
	)

	// Create another collector to scrape course details
	//detailCollector := c.Clone()

	//courses := make([]Course, 0, 200)

	// On every <a> element which has "href" attribute call callback
	// https://www.funda.nl/en/koop/amersfoort/huis-42447075-babylon-14/?navigateSource=resultlist
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//link := e.Attr("href")
		//if strings.HasPrefix(link,"https://www.funda.nl/en/koop/") {
		//	e.Request.Visit(link)
		//	fmt.Println("Visiting",link)
		//} else {
		//	fmt.Println(link)
		//}
		//fmt.Println(link)
	})
	c.OnRequest(func(f *colly.Request){
		f.Headers.Set("Referrer","https://www.funda.nl/")
		f.Headers.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	})
	c.OnResponse(func(r *colly.Response){
		fmt.Println(string(r.Body))
	})


	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	err = c.Visit("https://www.funda.nl/en/koop/amersfoort/beschikbaar/0-400000/100+woonopp/woonhuis/5+kamers/bouwperiode-1991-2000/bouwperiode-2001-2010/bouwperiode-2011-2020/bouwperiode-na-2020/energielabel-a/energielabel-b/")
	if err != nil {
		fmt.Println(err)
	}
}

// On every <a> element with collection-product-card class call callback
//	c.OnHTML(`a.collection-product-card`, func(e *colly.HTMLElement) {
//		// Activate detailCollector if the link contains "coursera.org/learn"
//		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
//		if strings.Index(courseURL, "coursera.org/learn") != -1 {
//			detailCollector.Visit(courseURL)
//		}
//	})
//
//	// Extract details of the course
//	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
//		log.Println("Course found", e.Request.URL)
//		title := e.ChildText(".banner-title")
//		if title == "" {
//			log.Println("No title found", e.Request.URL)
//		}
//		course := Course{
//			Title:       title,
//			URL:         e.Request.URL.String(),
//			Description: e.ChildText("div.content"),
//			Creator:     e.ChildText("li.banner-instructor-info > a > div > div > span"),
//			Rating:      e.ChildText("span.number-rating"),
//		}
//		// Iterate over div components and add details to course
//		e.ForEach(".AboutCourse .ProductGlance > div", func(_ int, el *colly.HTMLElement) {
//			svgTitle := strings.Split(el.ChildText("div:nth-child(1) svg title"), " ")
//			lastWord := svgTitle[len(svgTitle)-1]
//			switch lastWord {
//			// svg Title: Available Langauges
//			case "languages":
//				course.Language = el.ChildText("div:nth-child(2) > div:nth-child(1)")
//			// svg Title: Mixed/Beginner/Intermediate/Advanced Level
//			case "Level":
//				course.Level = el.ChildText("div:nth-child(2) > div:nth-child(1)")
//			// svg Title: Hours to complete
//			case "complete":
//				course.Commitment = el.ChildText("div:nth-child(2) > div:nth-child(1)")
//			}
//		})
//		courses = append(courses, course)
//	})
//
//	// Start scraping on http://coursera.com/browse
//	c.Visit("https://coursera.org/browse")
//
//	enc := json.NewEncoder(file)
//	enc.SetIndent("", "  ")
//
//	// Dump json to the standard output
//	enc.Encode(courses)
//}
