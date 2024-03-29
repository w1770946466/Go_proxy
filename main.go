
package main

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "net/http"
    "regexp"
    "os"
)

func main() {
    url := "https://t.me/s/ZDYZ2"
    links, err := getLinks(url)
    if err != nil {
        fmt.Printf("Failed to retrieve %s: %v\n", url, err)
        return
    }

    var base64Links []string
    for _, link := range links {
        body, err := getWebpageContent(link)
        if err != nil {
            fmt.Printf("Failed to retrieve %s: %v\n", link, err)
            continue
        }
        if isBase64(string(body)) {
            base64Links = append(base64Links, link)
        }
    }

    fmt.Printf("Links containing base64-encoded data: %v\n", base64Links)

    // Write links to sub file
    file, err := os.Create("sub")
    if err != nil {
        fmt.Printf("Failed to create file: %v\n", err)
        return
    }
    defer file.Close()

    for _, link := range base64Links {
        fmt.Fprintf(file, "%s\n", link)
    }

    fmt.Println("Links written to sub")
}

func getLinks(url string) ([]string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve %s: %v", url, err)
    }
    defer resp.Body.Close()

    // Define regular expression to match URLs
    // This pattern may not match all possible URLs, and might require some tweaks
    // depending on the specific HTML content being processed.
    urlPattern := `(?i)<a\s+(?:[^>]*?\s+)?href="([^"]*)"`
    r := regexp.MustCompile(urlPattern)

    var buf bytes.Buffer
    if _, err := buf.ReadFrom(resp.Body); err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    matches := r.FindAllStringSubmatch(buf.String(), -1)

    // Extract URLs from matched substrings
    var links []string
    for _, match := range matches {
        if len(match) >= 2 {
            links = append(links, match[1])
        }
    }

    return links, nil
}

func getWebpageContent(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve %s: %v", url, err)
    }
    defer resp.Body.Close()

    var buf bytes.Buffer
    if _, err := buf.ReadFrom(resp.Body); err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    return buf.Bytes(), nil
}

func isBase64(str string) bool {
    decoded, err := base64.StdEncoding.DecodeString(str)
    if err != nil {
        return false
    }
    for _, b := range decoded {
        if b > 127 {
            return false
        }
    }
    return true
}
