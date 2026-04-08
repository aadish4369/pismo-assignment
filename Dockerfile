FROM golang:bookworm

WORKDIR /app
ENV CGO_ENABLED=1 GOTOOLCHAIN=auto
ENV DATABASE_PATH=/data/app.db

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
CMD ["go", "run", "."]
