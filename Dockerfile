# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang


RUN go get github.com/maddyonline/gotutorial

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/gotutorial -base="/go/src/github.com/maddyonline/gotutorial"

# Document that the service listens on port 8080.
EXPOSE 8080
