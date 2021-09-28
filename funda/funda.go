package funda

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"time"
)

type Listing struct {
	Title            string
	Address          string
	LivingArea       string
	Bedrooms         string
	Price            string
	ListedSince      string
	Status           string
	ConstructionYear string
	Facilities       string
	EnergyLabel      string
	Insulation       string
	URL              string
}

var err error
var currentURL string

func GetListings(domain string, city string) []Listing {
	c := colly.NewCollector(
		colly.AllowedDomains(domain, "www."+domain),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"),
		colly.CacheDir("./listings_cache"),
	)
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*"+domain+".*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("Error while limiting scraping rules:", err)
	}

	// Create another collector to scrape listing details
	detailCollector := c.Clone()

	var listings []Listing
	c.OnHTML("div.search-result__header-title-col", func(e *colly.HTMLElement) {
		e.DOM.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			link, _ := s.Attr("href")
			currentURL = "https://www." + domain + link
			err = detailCollector.Visit(currentURL)
			if err != nil {
				fmt.Println(err)
			}
			return false
		})
	})
	c.OnHTML("a[rel=next]", func(e *colly.HTMLElement) {
		link := "https://www." + domain + e.Attr("href")
		err = c.Visit(link)
		if err != nil {
			fmt.Println("Error while trying to visit", link, err)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	if err != nil {
		fmt.Println(err)
	}
	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.DOM.Find("span.object-header__title").Text())
		fullAddress := strings.TrimSpace(e.DOM.Find("span.object-header__subtitle.fd-color-dark-3").Text())
		field := strings.Split(fullAddress, "\n")
		address := strings.TrimSpace(field[0])
		var city string
		if len(field) > 1 {
			city = strings.TrimSpace(field[1])
		}
		var livingArea, bedrooms, price, listedSince, status, constructionYear, facilities, energyLabel, insulation string
		e.DOM.Find("dt").Each(func(i int, s *goquery.Selection) {
			dd := s.NextFiltered("dd")
			if s.Text() == "Living area" {
				livingArea = strings.TrimSpace(dd.Text())
			}
			if s.Text() == "Number of rooms" {
				bedrooms = strings.TrimSpace(dd.Text())
			}
			if s.Text() == "Asking price" {
				price = strings.Replace(strings.TrimSpace(dd.Text()), " kosten koper", "", -1)
			}
			if strings.Contains(s.Text(), "Listed since") {
				listedSince = strings.TrimSpace(dd.Text())
			}
			if strings.Contains(s.Text(), "Status") {
				status = strings.TrimSpace(dd.Text())
			}
			if strings.Contains(s.Text(), "Year of construction") {
				constructionYear = strings.TrimSpace(dd.Text())
			}
			if strings.Contains(s.Text(), "Construction period") {
				constructionYear = strings.TrimSpace(dd.Text())
			}
			if strings.Contains(s.Text(), "Facilities") {
				if facilities != "" {
					facilities = facilities + ", " + strings.TrimSpace(dd.Text())
				} else {
					facilities = strings.TrimSpace(dd.Text())
				}
			}
			if strings.Contains(s.Text(), "Energy label") {
				energyLabel = strings.Split(strings.TrimSpace(dd.Text()), "\n")[0]
			}
			if strings.Contains(s.Text(), "Insulation") {
				insulation = strings.Split(strings.TrimSpace(dd.Text()), "\n")[0]
			}
		})
		listing := Listing{
			Title:            title,
			Address:          address + " " + city,
			LivingArea:       livingArea,
			Bedrooms:         bedrooms,
			Price:            price,
			ListedSince:      listedSince,
			Status:           status,
			ConstructionYear: constructionYear,
			Facilities:       facilities,
			EnergyLabel:      energyLabel,
			Insulation:       insulation,
			URL:              currentURL,
		}
		listings = append(listings, listing)
	})
	err = c.Visit("https://www." + domain + "/en/koop/" + city + "/")
	return listings
}
