package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Pendetot/AnimekApi/utils"
)

// AnimeScheduleItem represents a single anime in the schedule
type AnimeScheduleItem struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	BroadcastTime string `json:"broadcast_time"`
	ImageURL      string `json:"image_url"`
}

// DaySchedule represents anime schedule for a specific day
type DaySchedule struct {
	Day       string              `json:"day"`
	AnimeList []AnimeScheduleItem `json:"anime_list"`
}

// AnoboyJadwalItem represents an item from the anoboy jadwal table
type AnoboyJadwalItem struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Day      string `json:"day"`
	Time     string `json:"time"`
	ImageURL string `json:"image_url"`
}

// ScheduleResponse represents the complete schedule response
type ScheduleResponse struct {
	ScheduleByDay []DaySchedule      `json:"schedule_by_day"`
	AnoboyJadwal  []AnoboyJadwalItem `json:"anoboy_jadwal"`
}

// scrapeAnimeImage scrapes image from an anime page
func scrapeAnimeImage(url string) (string, error) {
	client := &http.Client{
		Timeout: 20 * time.Second, // Double timeout for image scraping
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", utils.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Try multiple selectors in order of specificity
	selectors := []string{
		".column-three-fourth amp-img",
		"amp-img[width="640"][height="360"]",
		".column-three-fourth img",
		"amp-img[width][height]",
		"amp-img",
		"img",
	}

	// Try each selector
	for _, selector := range selectors {
		element := doc.Find(selector).First()
		if element.Length() > 0 {
			if src, exists := element.Attr("src"); exists {
				// Make sure the URL is complete
				if strings.HasPrefix(src, "http") {
					return src, nil
				} else if strings.HasPrefix(src, "/") {
					// If it's a relative URL, make it absolute
					return "https://ww1.anoboy.app" + src, nil
				}
			}
		}
	}

	// If we couldn't find an image with the selectors, try the meta tags
	ogImage := doc.Find("meta[property="og:image"]")
	if ogImage.Length() > 0 {
		if content, exists := ogImage.Attr("content"); exists {
			return content, nil
		}
	}

	// If we still can't find an image, look for URLs in the HTML that appear to be images
	htmlStr, _ := doc.Html()
	imgRegex := regexp.MustCompile(`https?://[^"']+\.(jpg|jpeg|png|gif|webp)`)
	matches := imgRegex.FindAllString(htmlStr, -1)
	
	if len(matches) > 0 {
		// Find the first URL that contains the word "01as" which is common in the images
		for _, match := range matches {
			if strings.Contains(match, "01as") {
				return match, nil
			}
		}
		
		// Otherwise just return the first one that looks like a content image (not a UI element)
		for _, match := range matches {
			if !strings.Contains(match, "newlogo") &&
				!strings.Contains(match, "discord") &&
				!strings.Contains(match, "qq288") &&
				!strings.Contains(match, "icon") {
				return match, nil
			}
		}
	}

	// If we really can't find anything, return default image
	return "https://ww1.anoboy.app/wp-content/uploads/2019/02/cropped-512x512-192x192.png", nil
}

// extractScheduleData extracts links and anime data from schedule page
func extractScheduleData(doc *goquery.Document) []DaySchedule {
	var days []DaySchedule

	// Find all day sections (h1 in div.unduhan)
	dayDivs := doc.Find("div.unduhan").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find("h1").Length() > 0
	})

	dayDivs.Each(func(i int, dayDiv *goquery.Selection) {
		dayName := strings.TrimSpace(dayDiv.Find("h1").Text())

		// Find all anime entries (li elements) under this day
		var animeEntries []AnimeScheduleItem
		animeItems := dayDiv.Find("ul.lcp_catlist li")

		animeItems.Each(func(j int, item *goquery.Selection) {
			animeLink := item.Find("a").First()
			if animeLink.Length() > 0 {
				title := strings.TrimSpace(animeLink.Text())
				if url, exists := animeLink.Attr("href"); exists {
					// Extract broadcast time from the text after the link
					fullText := strings.TrimSpace(item.Text())
					timeRegex := regexp.MustCompile(`([^,]+,\s*\d+:\d+)$`)
					matches := timeRegex.FindStringSubmatch(fullText)
					
					var broadcastTime string
					if len(matches) > 1 {
						broadcastTime = strings.TrimSpace(matches[1])
					}

					animeEntries = append(animeEntries, AnimeScheduleItem{
						Title:         title,
						URL:           url,
						BroadcastTime: broadcastTime,
						ImageURL:      "", // Will be filled later
					})
				}
			}
		})

		if len(animeEntries) > 0 {
			days = append(days, DaySchedule{
				Day:       dayName,
				AnimeList: animeEntries,
			})
		}
	})

	return days
}

// ScrapeSchedule scrapes the anime broadcast schedule from jadwal page
func ScrapeSchedule() (*ScheduleResponse, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// URL to scrape (schedule page)
	url := "https://ww1.anoboy.app/2015/05/anime-subtitle-indonesia-ini-adalah-arsip-file-kami/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", utils.GetUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract schedule data
	scheduleData := extractScheduleData(doc)

	// Fetch images for all anime with improved error handling and retry logic
	for i := range scheduleData {
		for j := range scheduleData[i].AnimeList {
			anime := &scheduleData[i].AnimeList[j]
			if anime.URL != "" {
				imageUrl, err := scrapeAnimeImage(anime.URL)
				if err != nil {
					// Use default image on error
					anime.ImageURL = "https://ww1.anoboy.app/wp-content/uploads/2019/02/cropped-512x512-192x192.png"
				} else {
					anime.ImageURL = imageUrl
				}
				
				// Add a small delay to avoid overwhelming the server
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// Also scrape the anoBoy internal schedule table if available
	var anoboyJadwal []AnoboyJadwalItem
	jadwalTable := doc.Find("table").First()

	if jadwalTable.Length() > 0 {
		rows := jadwalTable.Find("tr")

		// Skip header row
		rows.Slice(1, goquery.ToEnd).Each(func(i int, row *goquery.Selection) {
			cols := row.Find("td")
			if cols.Length() >= 3 {
				titleElem := cols.Eq(0).Find("a").First()
				var title, url string
				
				if titleElem.Length() > 0 {
					title = strings.TrimSpace(titleElem.Text())
					if href, exists := titleElem.Attr("href"); exists {
						url = href
					}
				} else {
					title = strings.TrimSpace(cols.Eq(0).Text())
				}

				day := ""
				if cols.Length() > 1 {
					day = strings.TrimSpace(cols.Eq(1).Text())
				}

				time := ""
				if cols.Length() > 2 {
					time = strings.TrimSpace(cols.Eq(2).Text())
				}

				anoboyJadwal = append(anoboyJadwal, AnoboyJadwalItem{
					Title:    title,
					URL:      url,
					Day:      day,
					Time:     time,
					ImageURL: "", // Will be filled later if we decide to scrape images for these too
				})
			}
		})
	}

	return &ScheduleResponse{
		ScheduleByDay: scheduleData,
		AnoboyJadwal:  anoboyJadwal,
	}, nil
}
