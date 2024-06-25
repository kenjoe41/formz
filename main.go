package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/html"
)

func main() {
	filePath := flag.String("file", "", "Path to file containing URLs")
	followRedirects := flag.Bool("follow-redirects", true, "Whether to follow HTTP redirects")
	flag.Parse()

	var urls []string
	var err error

	if *filePath != "" {
		urls, err = getURLsFromFile(*filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			os.Exit(1)
		}
	} else {
		urls = getURLsFromStdin()
	}

	client := retryablehttp.NewClient()
	client.Logger = nil
	if !*followRedirects {
		client.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	for _, urlStr := range urls {
		cleanedURL := cleanURL(urlStr)
		htmlContent, err := fetchHTML(client, cleanedURL)
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error fetching URL:", err)
			continue
		}

		if checkCloudflare(htmlContent) {
			fmt.Fprintln(os.Stderr, "Blocked by Cloudflare detected for", cleanedURL, "skipping.")
			continue
		}

		if hasForm(htmlContent) {
			fmt.Println(cleanedURL)
		}
	}
}

func getURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func getURLsFromStdin() []string {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return urls
}

func cleanURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" {
		return "https://" + rawURL
	}
	return rawURL
}

func fetchHTML(client *retryablehttp.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func checkCloudflare(htmlContent string) bool {
	return strings.Contains(htmlContent, "Attention Required") &&
		strings.Contains(htmlContent, "Cloudflare")
}

func hasForm(htmlContent string) bool {
	z := html.NewTokenizer(strings.NewReader(htmlContent))
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return false
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "form" {
				return true
			}
		}
	}
}
