package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting → open http://localhost:%s\n", port)


	http.ListenAndServe("0.0.0.0:"+port, nil)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Boni Smart Query Helper</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
				input { padding: 10px; width: 300px; }
				button { padding: 10px 20px; background: #4CAF50; color: white; border: none; cursor: pointer; }
				button:hover { background: #45a049; }
				ul { list-style-type: none; padding: 0; }
				li { margin: 10px 0; }
				a { color: #0066cc; text-decoration: none; }
				a:hover { text-decoration: underline; }
			</style>
		</head>
		<body>
			<h1>Boni Smart Query Helper 🚀</h1>
			<p>Message <strong>+91 98000 81110</strong> on WhatsApp with detailed queries to get better, faster deals from competing vendors in Bangalore!</p>
			<p>Examples of smart queries (click to try):</p>
			<ul>
				<li><a href="https://wa.me/919800081110?text=Haircut%20in%20Koramangala%20under%20500%20today" target="_blank">Haircut in Koramangala under 500 today</a></li>
				<li><a href="https://wa.me/919800081110?text=Urgent%20AC%20repair%20Indiranagar%20same%20day%20best%20price" target="_blank">Urgent AC repair Indiranagar same day best price</a></li>
				<li><a href="https://wa.me/919800081110?text=Cheap%20flight%20to%20Goa%20from%20Bangalore%20this%20weekend" target="_blank">Cheap flight to Goa from Bangalore this weekend</a></li>
			</ul>

			<h2>Make your own smart query:</h2>
			<form method="GET" action="/generate">
				<input type="text" name="need" placeholder="e.g. plumber Whitefield or biryani MG Road" required>
				<button type="submit">Generate Smart Links</button>
			</form>

			<p style="margin-top: 40px; font-size: 0.9em;">Built to promote Boni – try detailed queries for the best results!</p>
		</body>
		</html>
	`)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	need := strings.TrimSpace(r.URL.Query().Get("need"))
	if need == "" {
		fmt.Fprint(w, "<h2>Error: Please enter something!</h2><a href='/'>Back</a>")
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	fmt.Println("DEBUG: API Key loaded?", apiKey != "") // Check terminal
	fmt.Println("DEBUG: User input:", need)

	if apiKey == "" {
		fmt.Fprint(w, "<h2>Error: OPENROUTER_API_KEY not set!</h2><p>Run in terminal: export OPENROUTER_API_KEY=sk-or-v1-...</p><a href='/'>Back</a>")
		return
	}

	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithModel("openrouter/free"), // Change to another :free model if rate-limited
	)
	if err != nil {
		fmt.Fprintf(w, "<h2>LLM Init Error: %v</h2><a href='/'>Back</a>", err)
		fmt.Println("DEBUG: Init error:", err)
		return
	}

	prompt := fmt.Sprintf(`You are an expert at creating effective search queries for Boni, a WhatsApp-based local search service in Bangalore.
User input: "%s"
Generate 5 natural, detailed, hyper-local queries that maximize better deals/results (include urgency, budget, location, preferences where relevant).
Output ONLY a numbered list like:
1. best emergency hospital Koramangala today
2. 24/7 hospital near Indiranagar low fees urgent
No extra text.`, need)

	ctx := context.Background()
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt) // <-- FIXED: Use llms. prefix + pass llm as arg
	if err != nil {
		fmt.Fprintf(w, "<h2>Generation Error: %v</h2><a href='/'>Back</a>", err)
		fmt.Println("DEBUG: Generate error:", err)
		return
	}

	fmt.Println("DEBUG: Raw AI response:", response) // Check what came back
	// Flexible parsing
	variations := []string{}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove number/dot/space or dashes
		line = strings.TrimLeft(line, "12345.-* ")
		line = strings.TrimSpace(line)
		if line != "" {
			variations = append(variations, line)
		}
	}

	// Fallback
	if len(variations) == 0 {
		variations = []string{
			need + " urgent today Bangalore",
			need + " best price low cost",
			need + " near me same day",
			need + " top rated",
			need + " cheap and good",
		}
	}

	// Results page
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Your Smart Boni Queries</title>
			<style>body { font-family: Arial; max-width: 800px; margin: 0 auto; padding: 20px; } ul { list-style-type: none; padding: 0; } li { margin: 15px 0; font-size: 1.1em; } a { color: #0066cc; }</style>
		</head>
		<body>
			<h1>Your Smart Queries for Boni</h1>
			<p>Original: <strong>`+need+`</strong></p>
			<p>Click to send to WhatsApp:</p>
			<ul>
	`)

	for _, q := range variations {
		encoded := url.QueryEscape(q)
		link := "https://wa.me/919800081110?text=" + encoded
		fmt.Fprintf(w, `<li><a href="%s" target="_blank">%s</a></li>`, link, q)
	}

	fmt.Fprint(w, `
			</ul>
			<p><a href="/">Back</a></p>
			<p style="font-size:0.9em;">Powered by AI (OpenRouter) for better Boni searches!</p>
		</body>
		</html>
	`)
}