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
var severitymap map[string]map[string][]pc.Vulnerability
var totalScans []pc.ScanArray
var mu sync.RWMutex // Mutex for concurrency

/*
Description: Function to get Jsonfiles from github
Functionality:
	1.Create the github link in raw format
	2.Attempt accessing the github for a maximum of GitMaxRetries
	3.If the connection is successful, read the content and Unmarshal the payload
	4.store the unmarshalled payload to a global store
Returns: Nothing
*/
func fetchJSONFromGitHub(githubRepo, filePath string, wg *sync.WaitGroup) {
        defer wg.Done() 
        rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/refs/heads/main/%s",githubRepo, filePath)
	fmt.Printf("Accessing File :%s in URL :%s \n", filePath, rawURL)

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

/*
Description: Function to handle scan endpoint
Functionality:
	1.Basic error condition checks
	2.Read the content insde http request and unmarshal the data
	3.Process Files concurrently with go threads and wait groups if there are many files
	4.Each go thread does a different file from github and store it in a common global map
	5.Go through the global map and store the needed information that can be used during query API
	6.Send the Response to the post message	
Returns: Nothing
*/
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
			if _,ok := severitymap[vul.Severity]; !ok {
				severitymap[vul.Severity] = make(map[string][]pc.Vulnerability)
			}
			/* I am assuming ID changes if there's content change in vulnerability
			If ID remain same, it has not changed to what is already cached */
			if _, exists := severitymap[vul.Severity][vul.ID]; !exists {
                       		 severitymap[vul.Severity][vul.ID] = append(severitymap[vul.Severity][vul.ID], vul)
			}
                }
        }
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Repository file contents fetched and store locally")
}

/*
Description: Function to handle query endpoint
Functionality:
	1.Basic error condition checks
	2.Read the content insde http request and unmarshal the data
	3.check the severity filter from the http request and match all the 
	  vulnerabilities of that severity
	6.Send the Response to the post message	
Returns: Nothing
*/
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
			for _,value := range result {
				sevResults = append(sevResults, value...)
			}
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
