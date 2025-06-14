package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Pendetot/AnimekApi/utils"
)

// Episode represents a single episode
type Episode struct {
	Title  string `json:"title"`
	Number string `json:"number"`
	URL    string `json:"url"`
}

// EpisodeList represents the episode list response
type EpisodeList struct {
	Title       string            `json:"title"`
	Episodes    []Episode         `json:"episodes"`
	Metadata    map[string]string `json:"metadata"`
	Description string            `json:"description"`
}

// extractEpisodeNumber extracts episode number from title using regex
func extractEpisodeNumber(title string) string {
	// Common patterns for episode numbers
	patterns := []string{
		`[Ee]pisode\s*(\d+)`,
		`[Ee]p\s*(\d+)`,
		`(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(title)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return "Unknown"
}

// ScrapeEpisodes scrapes the episode list from an anime series page
func ScrapeEpisodes(url string) (*EpisodeList, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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

	result := &EpisodeList{
		Metadata: make(map[string]string),
	}

	// Get anime title
	titleElement := doc.Find("h1").First()
	if titleElement.Length() > 0 {
		result.Title = strings.TrimSpace(titleElement.Text())
	} else {
		result.Title = "Title not found"
	}

	// Find episode list - check different possible containers
	var episodeSection *goquery.Selection

	// Try different possible containers for episode lists
	possibleContainers := []string{
		"div.singlelink",
		"div.hq",
		"ul.lcp_catlist",
		"div.episodes",
	}

	for _, selector := range possibleContainers {
		container := doc.Find(selector)
		if container.Length() > 0 {
			episodeSection = container
			break
		}
	}

	var episodes []Episode

	if episodeSection != nil && episodeSection.Length() > 0 {
		// Try to find episodes in list items
		episodeItems := episodeSection.Find("li")

		if episodeItems.Length() > 0 {
			episodeItems.Each(func(i int, item *goquery.Selection) {
				link := item.Find("a").First()
				if link.Length() > 0 {
					episodeTitle := strings.TrimSpace(link.Text())
					if episodeUrl, exists := link.Attr("href"); exists && episodeTitle != "" {
						episodeNumber := extractEpisodeNumber(episodeTitle)

						episodes = append(episodes, Episode{
							Title:  episodeTitle,
							Number: episodeNumber,
							URL:    episodeUrl,
						})
					}
				}
			})
		} else {
			// If no list items found, try to find direct links
			links := episodeSection.Find("a")
			links.Each(func(i int, link *goquery.Selection) {
				episodeTitle := strings.TrimSpace(link.Text())
				if episodeUrl, exists := link.Attr("href"); exists && episodeTitle != "" {
					episodeNumber := extractEpisodeNumber(episodeTitle)

					episodes = append(episodes, Episode{
						Title:  episodeTitle,
						Number: episodeNumber,
						URL:    episodeUrl,
					})
				}
			})
		}
	}

	// Sort episodes by number if possible
	sort.Slice(episodes, func(i, j int) bool {
		numA, errA := strconv.Atoi(episodes[i].Number)
		numB, errB := strconv.Atoi(episodes[j].Number)

		// If both are numbers, compare numerically
		if errA == nil && errB == nil {
			return numA < numB
		}

		// If one is "Unknown", put it at the end
		if episodes[i].Number == "Unknown" {
			return false
		}
		if episodes[j].Number == "Unknown" {
			return true
		}

		// If both are non-numeric, compare as strings
		return episodes[i].Number < episodes[j].Number
	})

	result.Episodes = episodes

	// Get anime metadata
	infoTable := doc.Find("table")
	if infoTable.Length() > 0 {
		rows := infoTable.Find("tr")
		rows.Each(func(i int, row *goquery.Selection) {
			header := row.Find("th").First()
			data := row.Find("td").First()
			if header.Length() > 0 && data.Length() > 0 {
				key := strings.TrimSpace(header.Text())
				value := strings.TrimSpace(data.Text())
				result.Metadata[key] = value
			}
		})
	}

	// Get anime description
	descElem := doc.Find("div.unduhan").First()
	if descElem.Length() > 0 {
		result.Description = strings.TrimSpace(descElem.Text())
	}

	return result, nil
}
