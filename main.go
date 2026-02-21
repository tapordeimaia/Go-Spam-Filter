package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type Model struct {
	HamCounts map[string]int
	SpamCounts map[string]int
	HamTotal int
	SpamTotal int
}

// This function takes a long string of text and breaks it down into a list of words
// It also removes punctuation and makes everything lowercase
func tokenize(text string) []string {
	// 1. Make everything lowercase
	text = strings.ToLower(text)

	// 2. Define a function that tells Go what the separator will be aka anything that is 	// not a letter or a number
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	// 3. Split the string based on the function f
	words := strings.FieldsFunc(text, f)
	return words
}

// This function will read every file in a specific folder and count how many times words // appear, and how many words there a re in total
func countWords(folderPath string) (map[string]int, int, error) {
	counts := make(map[string]int)
	total := 0

	// Read the directory
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, 0, err	
	}

	// Loop through every file in the folder
	for _, file := range files {
		// Read the file content
		content, err := os.ReadFile(folderPath + "/" + file.Name())
		if err != nil {
			// If one file fails, we just skip it
			continue
		}

		// Turn the content into words
		words := tokenize(string(content))

		// Add to out count map
		for _, word := range words {
			// Skip the tiny link words
			if len(word) > 2 {
				counts[word]++
				total++
			}
		}
	}
	return counts, total, nil
}

func (m *Model) Classify(text string) (string, float64, float64) {
	words := tokenize(text)
	
	vocabSize := 40000.0

	// Calculate the initial scores based on how much spam/ham data we have in total
	spamScore := math.Log(float64(m.SpamTotal) / float64(m.SpamTotal + m.HamTotal))
	hamScore := math.Log(float64(m.HamTotal) / float64(m.SpamTotal + m.HamTotal))

	for _, word := range words {
		// Calculate the spam probability for this word
		// We do the +1 so we prevent crashing if we have new words
		spamCount := m.SpamCounts[word]
		pSpam := float64(spamCount + 1) / float64(m.SpamTotal + int(vocabSize))
		spamScore += math.Log(pSpam)

		hamCount := m.HamCounts[word]
		pHam := float64(hamCount + 1) / float64(m.HamTotal + int(vocabSize))
		hamScore += math.Log(pHam)
	}

	fmt.Printf("DEBUG: Spam Score: %.2f | Ham Score: %.2f\n", spamScore, hamScore)
	if spamScore > hamScore {
		return "SPAM", spamScore, hamScore
	}
	return "HAM", spamScore, hamScore
}

func (m *Model) SaveToFile(filename string) error {
	// Marshal coverts the struct into a JSON byte array
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// Write the byte array to a file with standard permissions
	return os.WriteFile(filename, data, 0644)
}

func LoadFromFile(filename string) (*Model, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var m Model
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

type RequestPayload struct {
	Text string `json:"text"`
}

type ResponsePayload struct {
	Result string `json:"result"`
	SpamScore float64 `json:"spam_score"`
	HamScore float64 `json:"ham_score"`
}

var globalModel *Model

func classifyHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Text == "" {
		http.Error(w, "Invalid JSON or empty text", http.StatusBadRequest)
		return
	}

	result, spamScore, hamScore := globalModel.Classify(req.Text)

	res := ResponsePayload{Result: result, SpamScore: spamScore, HamScore: hamScore,}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	modelFile := "model.json"

	// Check if the model file already exists
	if _, err := os.Stat(modelFile); err == nil {
		// File exists
		fmt.Println("Found existing model.json. Loading model into memory...")
		globalModel, err = LoadFromFile(modelFile)
		if err != nil {
			panic("Failed to load model: " + err.Error())
		}
		fmt.Println("Model loaded instantly!")

	} else {
		fmt.Println("No existing model found. Training from Enron dataset...")
		hamCounts, hamTotal, _ := countWords("enron1/ham")
		spamCounts, spamTotal, _ := countWords("enron1/spam")

		globalModel = &Model{
			HamCounts: hamCounts, SpamCounts: spamCounts,
			HamTotal: hamTotal,   SpamTotal: spamTotal,
		}

		fmt.Println("Training complete! Saving to model.json...")
		err = globalModel.SaveToFile(modelFile)
		if err != nil {
			fmt.Println("Warning: Could not save model to file:", err)
		} else {
			fmt.Println("Model successfully saved!")
		}
	}

	// Start the server
	http.HandleFunc("/api/classify", classifyHandler)
	fmt.Println("\nServer is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
