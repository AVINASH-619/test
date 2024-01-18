package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Replace these variables with your GitHub token, username, repository, and branch.
const (
	githubToken    = "YOUR_GITHUB_TOKEN"
	githubUsername = "USERNAME"
	githubRepo     = "REPOSITORY"
	githubBranch   = "BRANCH"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	port := 3000
	fmt.Printf("Server is running at http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Encode file content to base64
	fileContentBase64 := base64.StdEncoding.EncodeToString(fileContent)

	// Create a GitHub API URL for creating or updating a file
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/temp_file.txt", githubUsername, githubRepo)

	// Create a PUT request to GitHub API
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBufferString(fmt.Sprintf(`{
		"message": "Add file",
		"content": "%s",
		"branch": "%s"
	}`, fileContentBase64, githubBranch)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set headers, including the GitHub personal access token
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to GitHub API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		fmt.Println("File uploaded successfully!")
		w.WriteHeader(http.StatusOK)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to upload file. Status: %d, Body: %s\n", resp.StatusCode, body)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
	}
}
