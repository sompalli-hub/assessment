package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"sync"
	"time"
	pc "github.com/sompalli-hub/assessment/payloadcontent"
)

var Appconfig *pc.Config
var severitymap map[string][]pc.Vulnerability
var totalScans []pc.ScanArray
var mu sync.RWMutex // Mutex for concurrency

// Function to fetch and parse a JSON file from GitHub
func fetchJSONFromGitHub(githubRepo, filePath string, wg *sync.WaitGroup) {
        defer wg.Done() 
        rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/refs/heads/main/%s",githubRepo, filePath)
	fmt.Printf("Accessing File :%s \n", filePath)

	for attempts := 1; attempts <= Appconfig.GitMaxRetries; attempts++ {
	        resp, err := http.Get(rawURL)
	        if err != nil {
			log.Printf("Error fetching file %s: %v attempt:%d", filePath, err, attempts)
			attempts += 1
			time.Sleep(1 * time.Second)
	        } else {
		        defer resp.Body.Close()
		
		        if resp.StatusCode != http.StatusOK {
		                log.Printf("Failed to fetch file %s: Status %d", filePath, resp.StatusCode)
		                return
		        }
		
		        body, err := ioutil.ReadAll(resp.Body)
		        if err != nil {
		                log.Printf("Failed to fetch file %s: Status %d", filePath, resp.StatusCode)
		                return
		        }
		
		        var scans []pc.ScanArray
		        err = json.Unmarshal(body, &scans)
		        if err != nil {
		                log.Printf("Error reading file %s: %v", filePath, err)
		                return
		        }
		
		        // Store scan results safely using a mutex
		        mu.Lock()
		        for _, scan := range scans {
		                totalScans = append(totalScans, scan)
		        }
		        mu.Unlock()
			break
		}
	}
}

// Handle the POST /scan endpoint to fetch JSON from GitHub and process it
func handleScan(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var scanRequest pc.PostScan
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &scanRequest); err != nil {
		http.Error(w, "Error unmarshaling request", http.StatusBadRequest)
		return
	}

	// Process multiple files concurrently
	var wg sync.WaitGroup
	for _, fileName := range scanRequest.Files {
		wg.Add(1)
		go fetchJSONFromGitHub(scanRequest.Repo, fileName, &wg)
	}
	wg.Wait() // Wait for all Goroutines to complete

	mu.Lock()
        for _, scan := range totalScans {
		for _, vul := range scan.KeyScanResult.Vulnerabilities {
                        //fmt.Printf("Severity: %s \n",vul.Severity)
                        severitymap[vul.Severity] = append(severitymap[vul.Severity], vul)
                }
        }
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Repository file contents fetched and store locally")
}

// Handle the POST /query endpoint to filter scan results
func handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var query pc.PostQuery
	if err := json.Unmarshal(body, &query); err != nil {
		http.Error(w, "Error unmarshaling filter request", http.StatusInternalServerError)
		return
	}

	severity := query.Filters.Severity

	var sevResults []pc.Vulnerability
	mu.RLock()
	for sev,result := range severitymap {
		if sev == severity {
			sevResults = result
		}
	}
	mu.RUnlock()

	queryrsp, err := json.Marshal(sevResults)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(queryrsp)
}
