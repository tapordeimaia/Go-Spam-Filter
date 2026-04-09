# Go Spam Filter API

A lightweight, RESTful API that classifies text as "Spam" or "Ham" using a Naive Bayes probabilistic model. Built entirely from scratch in Go without external Machine Learning libraries, containerized with Docker, and automatically deployed to Google Cloud.

## 🚀 Tech Stack
* **Language:** Go (Golang)
* **Concepts:** Machine Learning (Naive Bayes), NLP (Bag of Words, Tokenization), RESTful APIs, Log-probability scaling.
* **Infrastructure & DevOps:** Docker, Google Cloud Platform (Cloud Run, Artifact Registry), CI/CD (GitHub Actions).

## 🌍 Live API Endpoint
The API is containerized and deployed publicly on Google Cloud Run. You can test it from any terminal right now.

**Send a test request:**
```powershell
Invoke-RestMethod -Uri "[https://go-spam-api-7aj26y5hta-lm.a.run.app/api/classify](https://go-spam-api-7aj26y5hta-lm.a.run.app/api/classify)" -Method Post -ContentType "application/json" -Body '{"text": "Congratulations! You won a free $1000 gift card."}'
