package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Pendetot/AnimekApi/utils"
)

// AnimeItem represents a single anime item
type AnimeItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// AnimeGroup represents a group of anime organized by letter
type AnimeGroup struct {
	Letter    string      `json:"letter"`
	AnimeList []AnimeItem `json:"anime_list"`
}

// ScrapeAnimeList scrapes the overall anime list from the website
func ScrapeAnimeList() ([]AnimeGroup, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// URL to scrape
	url := "https://ww1.anoboy.app/anime-list/"

	// Create request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", utils.GetUserAgent())

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var animeGroups []AnimeGroup

	// Find anime list structure - typically organized by first letter
	letterGroups := doc.Find("div[class*="letter-group"], ul.lcp_catlist")

	// If the standard structure isn't found, try to find all list items
	if letterGroups.Length() == 0 {
		var animeList []AnimeItem
		
		allListItems := doc.Find("li")
		allListItems.Each(func(i int, item *goquery.Selection) {
			link := item.Find("a").First()
			if link.Length() > 0 {
				title := strings.TrimSpace(link.Text())
				if url, exists := link.Attr("href"); exists && title != "" {
					animeList = append(animeList, AnimeItem{
						Title: title,
						URL:   url,
					})
				}
			}
		})

		if len(animeList) > 0 {
			animeGroups = append(animeGroups, AnimeGroup{
				Letter:    "All",
				AnimeList: animeList,
			})
		}
	} else {
		// Process each letter group
		letterGroups.Each(func(i int, group *goquery.Selection) {
			// Try to find the letter/heading
			heading := group.Find("h2, h3").First()
			letter := "Unknown"
			if heading.Length() > 0 {
				letter = strings.TrimSpace(heading.Text())
			}

			// Find all anime links in this group
			links := group.Find("a")
			var animeList []AnimeItem

			links.Each(func(j int, link *goquery.Selection) {
				title := strings.TrimSpace(link.Text())
				if url, exists := link.Attr("href"); exists && title != "" {
					animeList = append(animeList, AnimeItem{
						Title: title,
						URL:   url,
					})
				}
			})

			if len(animeList) > 0 {
				animeGroups = append(animeGroups, AnimeGroup{
					Letter:    letter,
					AnimeList: animeList,
				})
			}
		})
	}

	return animeGroups, nil
}
