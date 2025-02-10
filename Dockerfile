# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /test

# Copy the source code into the container
COPY . .

RUN go get .

# Build the Go application
RUN make build
#RUN go build -o assessment main.go json.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run when the container starts
CMD ["/test/assessment"]
