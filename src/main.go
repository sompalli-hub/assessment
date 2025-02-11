package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"gopkg.in/yaml.v3"
	pc "github.com/sompalli-hub/assessment/payloadcontent"
)

/*
description: Read the config file from internal path 
Functionality:
	1.Read the file and unmarshal the config parameters
Return: Nothing
*/
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

/*
description: Start point of the application
Functionality:
	1.Read the config file and save the values locally
	2.Start the server on configured port and listen on scan and query endpoints
Return: Nothing
*/
func main() {
	var err error
	//Read config file
	Appconfig, err = ReadConfig("../env/config.yaml")
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
