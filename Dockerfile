docker build -t plexauthz .
docker run -p 8080:8080 plexauthz

# Use the official Golang image from the DockerHub
FROM golang:1.19 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Run make build
RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux go build -o plexauthz-server .
RUN make build

######## Start a new stage from scratch #######
FROM alpine:latest

# Set an environment variable for the HTTP and GRPC ports
ENV GRPC_PORT 7777
ENV HTTP_PORT 7778

WORKDIR /app

# Copy the plexauthz-server binary from the previous stage
COPY --from=builder /app/plexauthz-server /app/

# Expose ports to the outside world
EXPOSE $GRPC_PORT
EXPOSE $HTTP_PORT

# Command to run the executable
CMD ["./plexauthz-server"]

#docker build -t plexauthz .
#docker run -e GRPC_PORT=7777 -e HTTP_PORT=7778 plexauthz
