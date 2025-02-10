package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"gopkg.in/yaml.v3"
	pc "github.com/sompalli-hub/assessment/payloadcontent"
)

// **Function to Read YAML Config**
func ReadConfig(filename string) (*pc.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config pc.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}


func main() {
	//Read config file
	var err error
	Appconfig, err = ReadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	//fmt.Println("value of config are", Appconfig.ServerAddr, Appconfig.ServerPort, Appconfig.GitMaxRetries)
	//Start a https server and expose two end points
	severitymap = make(map[string][]pc.Vulnerability)	
	http.HandleFunc("/scan", handleScan)
	http.HandleFunc("/query", handleQuery)

	serverparam := Appconfig.ServerAddr + ":" + Appconfig.ServerPort

	// Start the server on port 8080
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(serverparam, nil))

}
