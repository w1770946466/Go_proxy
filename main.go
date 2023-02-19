package main

import (
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("https://t.me/s/masco899")
	if err != nil {
		fmt.Printf("http.Get() failed with '%s'\n", err)
		return
	}
	defer resp.Body.Close()

	links, err := ExtractLinks(resp.Body)
	if err != nil {
		fmt.Printf("ExtractLinks() failed with '%s'\n", err)
		return
	}

	for _, link := range links {
		fmt.Println(link)
	}
}

// ExtractLinks returns all the href links found in the provided HTML body
func ExtractLinks(body io.Reader) ([]string, error) {
	links := []string{}

	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return links, nil
			}
			return nil, err

		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data != "a" {
				continue
			}

			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}
