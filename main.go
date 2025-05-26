package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func fetchHTML(targetURL string) *html.Node {
	// Add rate limiting
	time.Sleep(1 * time.Second)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", targetURL, err)
		return nil
	}

	// Add user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching %s: %v", targetURL, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Non-200 status code for %s: %d", targetURL, resp.StatusCode)
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML from %s: %v", targetURL, err)
		return nil
	}
	return doc
}

func extractText(n *html.Node, words *[]string) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		// Improved text extraction - filter out common noise
		if len(text) > 2 && len(text) < 50 {
			// Skip common noise words
			noise := []string{"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by", "from", "up", "about", "into", "through", "during", "before", "after", "above", "below", "between", "among", "within", "without", "under", "over"}
			isNoise := false
			lowerText := strings.ToLower(text)
			for _, n := range noise {
				if lowerText == n {
					isNoise = true
					break
				}
			}
			if !isNoise {
				*words = append(*words, text)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, words)
	}
}

func filterWords(raw []string) []string {
	seen := map[string]bool{}
	valid := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	result := []string{}
	for _, word := range raw {
		w := strings.ToLower(word)
		if valid.MatchString(w) && !seen[w] && len(w) > 2 {
			seen[w] = true
			result = append(result, w)
		}
	}
	return result
}

// Leetspeak transformation
func toLeetspeak(s string) string {
	replacements := map[string]string{
		"a": "@", "A": "@",
		"e": "3", "E": "3",
		"i": "1", "I": "1",
		"o": "0", "O": "0",
		"s": "$", "S": "$",
		"t": "7", "T": "7",
		"l": "1", "L": "1",
		"g": "9", "G": "9",
	}

	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

func generateCombinations(base string, context []string) []string {
	base = strings.ReplaceAll(strings.ToLower(base), " ", "")

	// Expanded suffixes for CTF
	suffixes := []string{"123", "!", "2024", "2025", "@", "#", "1", "01", "321", "1337", "2023", "2022", "2021", "2020", "00", "99", "$$", "admin", "user", "test", "flag", "ctf"}

	// Common years and patterns
	years := []string{"2020", "2021", "2022", "2023", "2024", "2025", "1999", "2000", "1337"}

	// Keyboard patterns
	keyboards := []string{"qwerty", "asdf", "123456", "password", "admin", "root"}

	// CTF-specific terms
	ctfTerms := []string{"flag", "admin", "secret", "hidden", "key", "pass", "login", "auth", "token", "hash", "crypto", "encode", "decode", "ctf", "challenge", "pwn", "web", "misc", "forensics", "reverse", "binary"}

	results := []string{}

	// Base word variations
	results = append(results, base)
	results = append(results, strings.ToUpper(base), strings.Title(base))
	results = append(results, reverse(base))
	results = append(results, toLeetspeak(base))

	// Common word transformations
	results = append(results, base+"s", base+"es", base+"ing", base+"ed")
	results = append(results, base+"ly", base+"ful", base+"less", base+"ish")
	results = append(results, base+"er", base+"est", base+"tion", base+"ment")
	results = append(results, base+"ity", base+"ness", base+"ism", base+"ist")
	results = append(results, base+"ize", base+"ise", base+"ation", base+"ification")
	results = append(results, base+"ology", base+"graphy", base+"ics")
	results = append(results, base+"verse", base+"world", base+"net", base+"sys")

	// CTF-specific patterns
	for _, term := range ctfTerms {
		results = append(results, base+term, term+base, base+"_"+term, term+"_"+base)
	}

	// Year combinations
	for _, year := range years {
		results = append(results, base+year, year+base)
	}

	// Keyboard patterns
	for _, kb := range keyboards {
		results = append(results, base+kb, kb+base)
	}

	// Number patterns
	for i := 0; i <= 99; i++ {
		if i < 10 {
			results = append(results, base+fmt.Sprintf("0%d", i))
		}
		results = append(results, base+fmt.Sprintf("%d", i))
	}

	// Common prefixes
	prefixes := []string{"admin", "user", "test", "super", "root", "guest", "demo", "temp", "new", "old", "backup", "hidden", "secret"}
	for _, prefix := range prefixes {
		results = append(results, prefix+base, prefix+"_"+base, prefix+"-"+base)
	}

	// Context word combinations
	for _, word := range context {
		if len(word) < 3 || len(word) > 20 {
			continue
		}

		combos := []string{
			word,
			strings.ToLower(word),
			strings.ToUpper(word),
			strings.Title(word),
			reverse(word),
			toLeetspeak(word),
			base + word,
			word + base,
			base + "_" + word,
			word + "_" + base,
			base + "-" + word,
			word + "-" + base,
			base + "." + word,
			word + "." + base,
		}

		for _, combo := range combos {
			results = append(results, combo)
			// Add suffixes to combinations
			for _, s := range suffixes {
				results = append(results, combo+s)
			}
			// Add years to combinations
			for _, year := range years {
				results = append(results, combo+year)
			}
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, word := range results {
		if !seen[word] && word != "" {
			seen[word] = true
			unique = append(unique, word)
		}
	}

	return unique
}

func writeToFile(filename string, words []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, w := range words {
		fmt.Fprintln(f, w)
	}
}

func buildGoogleSearchURL(base string) string {
	q := url.QueryEscape(base + " CTF challenge")
	return "https://www.google.com/search?q=" + q
}

func buildWikipediaSearchURL(base string) string {
	q := strings.ReplaceAll(strings.ToLower(base), " ", "+")
	return "https://en.wikipedia.org/w/index.php?search=" + q
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"<base word>\"")
		fmt.Println("Example: go run main.go \"hackthebox\"")
		os.Exit(1)
	}

	base := os.Args[1]
	fmt.Printf("Generating wordlist for: %s\n", base)

	wikiURL := buildWikipediaSearchURL(base)
	googleURL := buildGoogleSearchURL(base)

	fmt.Println("Fetching context from Wikipedia...")
	doc1 := fetchHTML(wikiURL)

	fmt.Println("Fetching context from Google...")
	doc2 := fetchHTML(googleURL)

	var raw []string
	if doc1 != nil {
		extractText(doc1, &raw)
	}
	if doc2 != nil {
		extractText(doc2, &raw)
	}

	fmt.Printf("Extracted %d raw words\n", len(raw))

	context := filterWords(raw)
	fmt.Printf("Filtered to %d context words\n", len(context))

	wordlist := generateCombinations(base, context)
	fmt.Printf("Generated %d total combinations\n", len(wordlist))

	filename := strings.ReplaceAll(strings.ToLower(base), " ", "") + "_wordlist.txt"
	writeToFile(filename, wordlist)
	fmt.Printf("Wordlist saved to %s\n", filename)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
