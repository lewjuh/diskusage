FROM golang:1.24

# Install reflex for file watching
RUN go install github.com/cespare/reflex@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add /go/bin to PATH
ENV PATH="/go/bin:${PATH}"

CMD ["reflex", "-r", "\\.go$", "--", "go", "test", "./..."]