# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /test

# Copy the source code into the container
COPY . .

#Move to src directory to get/download the packages
WORKDIR /test/src

RUN go get .

#Move to build directory to run the Makefile
WORKDIR /test/build

# Build the Go application
RUN make build
#RUN go build -o assessment main.go json.go

#Move to bin directory to run the application
WORKDIR /test/bin/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run when the container starts
CMD ["/test/bin/assessment"]
