FROM golang:1.24

WORKDIR /app
COPY . .

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# things like templ generate should already have been ran
# see setup-dev.sh if they are missing
RUN CGO_ENABLED=0 GOOS=linux go build -o /roommates-app

CMD ["/roommates-app"]