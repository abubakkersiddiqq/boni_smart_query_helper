package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	homeTmpl    *template.Template
	resultsTmpl *template.Template
)

func main() {
	var err error

	homeTmpl, err = template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal("Failed to parse home.html:", err)
	}

	resultsTmpl, err = template.ParseFiles("templates/results.html")
	if err != nil {
		log.Fatal("Failed to parse results.html:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting → open http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := homeTmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ResultsData struct for template
type ResultsData struct {
	Need       string
	Variations []struct {
		Query string
		Link  string
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	need := strings.TrimSpace(r.URL.Query().Get("need"))
	if need == "" {
		fmt.Fprint(w, "<h2>Error: Please enter something!</h2><a href='/'>Back</a>")
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	log.Printf("DEBUG: API Key loaded? %v", apiKey != "")

	if apiKey == "" {
		fmt.Fprint(w, "<h2>Error: OPENROUTER_API_KEY not set!</h2><a href='/'>Back</a>")
		return
	}

	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel("openrouter/free"),
	)
	if err != nil {
		log.Printf("LLM Init Error: %v", err)
		fmt.Fprintf(w, "<h2>LLM Init Error: %v</h2><a href='/'>Back</a>", err)
		return
	}

	prompt := fmt.Sprintf(`You are an expert at creating effective search queries for Bino, a WhatsApp-based local search service in Bangalore.
User input: "%s"
Generate 5 natural, detailed, hyper-local queries that maximize better deals/results (include urgency, budget, location, preferences where relevant).
Output ONLY a numbered list like:
1. best emergency hospital Koramangala today
2. 24/7 hospital near Indiranagar low fees urgent
No extra text.`, need)

	ctx := context.Background()
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Printf("Generate error: %v", err)
		fmt.Fprintf(w, "<h2>Generation Error: %v</h2><a href='/'>Back</a>", err)
		return
	}

	log.Printf("Raw AI response: %s", response)

	// Parse variations
	variations := []string{}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimLeft(line, "12345.-* ")
		line = strings.TrimSpace(line)
		if line != "" {
			variations = append(variations, line)
		}
	}

	if len(variations) == 0 {
		variations = []string{
			need + " urgent today Bangalore",
			need + " best price low cost",
			need + " near me same day",
		}
	}

	// Prepare data for template
	data := ResultsData{
		Need: need,
	}

	for _, q := range variations {
		encoded := url.QueryEscape(q)
		link := "https://wa.me/919800081110?text=" + encoded
		data.Variations = append(data.Variations, struct {
			Query string
			Link  string
		}{Query: q, Link: link})
	}

	// Execute the template with data
	if err := resultsTmpl.Execute(w, data); err != nil {
		log.Printf("Template execute error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}