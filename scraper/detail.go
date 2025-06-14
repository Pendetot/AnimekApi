package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Pendetot/AnimekApi/utils"
)

// AnimeDetail represents anime episode details
type AnimeDetail struct {
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	ImageURL       string           `json:"image_url"`
	StreamingLinks []StreamingLink  `json:"streaming_links"`
	DownloadLinks  []DownloadLink   `json:"download_links"`
	Metadata       map[string]string `json:"metadata"`
}

// StreamingLink represents a streaming source
type StreamingLink struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Quality string `json:"quality"`
}

// DownloadLink represents a download source
type DownloadLink struct {
	Server  string `json:"server"`
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

// ScrapeAnimeDetail scrapes details from an anime episode page
func ScrapeAnimeDetail(url string) (*AnimeDetail, error) {
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

	detail := &AnimeDetail{
		Metadata: make(map[string]string),
	}

	// Get anime title
	titleElement := doc.Find("h1").First()
	if titleElement.Length() > 0 {
		detail.Title = strings.TrimSpace(titleElement.Text())
	} else {
		detail.Title = "Title not found"
	}

	// Get anime description
	descElement := doc.Find("div.contentdeks").First()
	if descElement.Length() > 0 {
		detail.Description = strings.TrimSpace(descElement.Text())
	} else {
		detail.Description = "Description not found"
	}

	// Get anime image
	imgElement := doc.Find("amp-img.gambar").First()
	if imgElement.Length() > 0 {
		if src, exists := imgElement.Attr("src"); exists {
			detail.ImageURL = src
		}
	}
	if detail.ImageURL == "" {
		detail.ImageURL = "Image not found"
	}

	// Get streaming sources
	detail.StreamingLinks = []StreamingLink{}

	// Get main player iframe source
	playerIframe := doc.Find("iframe#mediaplayer").First()
	if playerIframe.Length() > 0 {
		if src, exists := playerIframe.Attr("src"); exists {
			detail.StreamingLinks = append(detail.StreamingLinks, StreamingLink{
				Name:    "Main Player",
				URL:     src,
				Quality: "Default",
			})
		}
	}

	// Get alternative streaming sources
	altPlayers := doc.Find("a#allmiror")
	altPlayers.Each(func(i int, s *goquery.Selection) {
		if dataVideo, exists := s.Attr("data-video"); exists {
			playerText := strings.TrimSpace(s.Text())
			if playerText == "" {
				playerText = "Alternative Player"
			}

			quality := "Unknown"
			// Simple quality detection - can be improved
			if strings.Contains(strings.ToLower(playerText), "720") {
				quality = "720P"
			} else if strings.Contains(strings.ToLower(playerText), "480") {
				quality = "480P"
			} else if strings.Contains(strings.ToLower(playerText), "360") {
				quality = "360P"
			}

			detail.StreamingLinks = append(detail.StreamingLinks, StreamingLink{
				Name:    playerText,
				URL:     dataVideo,
				Quality: quality,
			})
		}
	})

	// Get download links
	detail.DownloadLinks = []DownloadLink{}
	downloadSection := doc.Find("div.download")
	if downloadSection.Length() > 0 {
		linkSpans := downloadSection.Find("span.ud")
		linkSpans.Each(func(i int, span *goquery.Selection) {
			serverName := span.Find("span.udj").First()
			server := "Unknown Server"
			if serverName.Length() > 0 {
				server = strings.TrimSpace(serverName.Text())
			}

			qualityLinks := span.Find("a.udl")
			qualityLinks.Each(func(j int, link *goquery.Selection) {
				if href, exists := link.Attr("href"); exists && href != "none" {
					if style, hasStyle := link.Attr("style"); !hasStyle || !strings.Contains(style, "display:none") {
						detail.DownloadLinks = append(detail.DownloadLinks, DownloadLink{
							Server:  server,
							Quality: strings.TrimSpace(link.Text()),
							URL:     href,
						})
					}
				}
			})
		})
	}

	// Get metadata
	infoTable := doc.Find("div.contenttable")
	if infoTable.Length() > 0 {
		rows := infoTable.Find("tr")
		rows.Each(func(i int, row *goquery.Selection) {
			header := row.Find("th").First()
			data := row.Find("td").First()
			if header.Length() > 0 && data.Length() > 0 {
				key := strings.TrimSpace(header.Text())
				value := strings.TrimSpace(data.Text())
				detail.Metadata[key] = value
			}
		})
	}

	return detail, nil
}
