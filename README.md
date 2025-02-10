# assessment
Project assessment

Steps to clone the Repository
=============================

1.create a code folder in your home directory

2.cd code

3.git clone git@github.com:sompalli-hub/assessment.git

4.cd assessment

5.go get . <=== To get all dependencies

6.Make build <=== To get the app binary


Test the changes Locally
========================

1.Execute the Binary ( ./assessment)

  |_ It will host the sever on localhost:8080 port.

  |_ The port and server details can be changed if required in config.yaml file

2. To scan the files in Github repository, please use the following curl command :

  |_ curl -X POST -H "Content-Type: application/json" -d '{ "repo": "velancio/vulnerability_scans", "files": ["vulnscan1011.json", "vulnscan1213.json", "vulnscan15.json", "vulnscan19.json"] }' http://localhost:8080/scan 

3. To query the stored json content with some filter, use the following curl command:

  |_ curl -X POST -H "Content-Type: application/json" -d '{ "filters": {"severity": "LOW"} }' http://localhost:8080/query | json_pp


To Build the docker image
==========================

The dockerfile is already present in the directory. docker image name can be tagged as per your choice. Here's the Example:

1. sudo docker build -t assessment_image . <=== To get the docker Image

2. sudo docker run -p 8080:8080 assessment_image <== To run the docker image 

Once the docker image is run, the same curl commands that's simulated above can be used to hit the container port with http traffic.
