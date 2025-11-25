package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// ----------------------------
// Scrape a single page and return content as string
// ----------------------------
func scrapePage(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("html parse error: %v", err)
	}

	// Multiple selectors to catch all content blocks
	selectors := []string{
		"article.gem-c-govspeak",
		"div.govuk-govspeak",
		"div.govspeak",
		".gem-c-govspeak",
	}

	var contentBlocks []*goquery.Selection
	for _, sel := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			contentBlocks = append(contentBlocks, s)
		})
	}

	if len(contentBlocks) == 0 {
		return "", fmt.Errorf("could not find content container for %s", url)
	}

	var builder strings.Builder

	// Loop through all content blocks
	for _, content := range contentBlocks {
		content.Find("h1, h2, h3, p, ul, ol").Each(func(i int, s *goquery.Selection) {
			tag := goquery.NodeName(s)

			switch tag {
			case "h1", "h2", "h3":
				title := clean(s.Text())
				builder.WriteString("\n" + title + "\n")
				builder.WriteString(strings.Repeat("-", len(title)) + "\n\n")

			case "p":
				txt := clean(s.Text())
				if txt != "" {
					builder.WriteString(txt + "\n\n")
				}

			case "ul", "ol":
				s.Find("li").Each(func(j int, li *goquery.Selection) {
					builder.WriteString("• " + clean(li.Text()) + "\n")
				})
				builder.WriteString("\n")
			}
		})
	}

	return builder.String(), nil
}

func clean(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
}

// ----------------------------
// Main
// ----------------------------
func main() {
	// Manual list of links
	links := []string{
		"https://www.gov.uk/guidance/the-highway-code/introduction",
		"https://www.gov.uk/guidance/the-highway-code/rules-for-pedestrians-1-to-35",
		"https://www.gov.uk/guidance/the-highway-code/rules-for-users-of-powered-wheelchairs-and-mobility-scooters-36-to-46",
		"https://www.gov.uk/guidance/the-highway-code/rules-about-animals-47-to-58",
		"https://www.gov.uk/guidance/the-highway-code/rules-for-cyclists-59-to-82",
		"https://www.gov.uk/guidance/the-highway-code/rules-for-motorcyclists-83-to-88",
		"https://www.gov.uk/guidance/the-highway-code/rules-for-drivers-and-motorcyclists-89-to-102",
		"https://www.gov.uk/guidance/the-highway-code/general-rules-techniques-and-advice-for-all-drivers-and-riders-103-to-158",
		"https://www.gov.uk/guidance/the-highway-code/using-the-road-159-to-203",
		"https://www.gov.uk/guidance/the-highway-code/road-users-requiring-extra-care-204-to-225",
		"https://www.gov.uk/guidance/the-highway-code/driving-in-adverse-weather-conditions-226-to-237",
		"https://www.gov.uk/guidance/the-highway-code/waiting-and-parking-238-to-252",
		"https://www.gov.uk/guidance/the-highway-code/motorways-253-to-274",
		"https://www.gov.uk/guidance/the-highway-code/breakdowns-and-incidents-275-to-287",
		"https://www.gov.uk/guidance/the-highway-code/road-works-level-crossings-and-tramways-288-to-307",
	}

	os.MkdirAll("pages", 0755)

	var fullBuilder strings.Builder
	fullBuilder.WriteString("THE HIGHWAY CODE — FULL TEXT\n=============================\n\n")

	for i, url := range links {
		fmt.Printf("[%d/%d] Scraping %s\n", i+1, len(links), url)
		content, err := scrapePage(url)
		if err != nil {
			fmt.Println("  ERROR:", err)
			continue
		}

		// Save individual page
		filename := "pages/" + strings.ReplaceAll(strings.TrimPrefix(url, "https://www.gov.uk/guidance/the-highway-code/"), "/", "_") + ".txt"
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Println("  ERROR writing file:", err)
			continue
		}
		fmt.Println("  Saved:", filename)

		// Append to full text
		fullBuilder.WriteString("\n===== " + strings.ReplaceAll(strings.TrimPrefix(url, "https://www.gov.uk/guidance/the-highway-code/"), "/", " ") + " =====\n\n")
		fullBuilder.WriteString(content)
	}

	// Save combined file
	if err := os.WriteFile("highway_code_full.txt", []byte(fullBuilder.String()), 0644); err != nil {
		log.Fatal("ERROR writing full text:", err)
	}

	fmt.Println("\nAll done! Individual pages in ./pages/, full text in highway_code_full.txt")
}
