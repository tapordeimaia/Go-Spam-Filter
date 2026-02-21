# Go Spam Filter API

A lightweight, RESTful API that classifies text as "Spam" or "Ham" using a Naive Bayes probabilistic model. Built entirely from scratch in Go without external Machine Learning libraries.

## 🚀 Tech Stack
* **Language:** Go (Golang)
* **Concepts:** Machine Learning (Naive Bayes), NLP (Bag of Words, Tokenization), RESTful APIs, Log-probability scaling.

## 🧠 How it Works
1. **Training:** The model ingests the Enron Email Dataset, tokenizes the text, and calculates word frequencies.
2. **Math:** It uses Laplace smoothing and sums logarithmic probabilities to prevent floating-point underflow when handling large datasets.
3. **Serving:** Exposes a local HTTP server to accept JSON payloads and return classification scores in real-time.

## 💻 How to Run Locally

1. Clone the repository.
2. Download the [Enron-Spam Dataset (enron1)](https://www2.aueb.gr/users/ion/data/enron-spam/) and extract it into an `enron1/` directory at the root level.
3. Start the server:
   ```bash
   go run main.go
4. Send a test POST request:
   ```bash
   Invoke-RestMethod -Uri "http://localhost:8080/api/classify" -Method Post -ContentType "application/json" -Body '{"text": "Congratulations! You won a free $1000 gift card."}'
