# Development stage for Golang using Alpine
FROM golang:1.20-alpine AS development

# Set timezone if needed
ENV TZ=Asia/Bangkok
RUN apk add --update tzdata && \
  cp /usr/share/zoneinfo/$TZ /etc/localtime && \
  echo $TZ > /etc/timezone

# Install git, build-base (for C libraries), protoc, and protobuf-dev
# RUN apk add --no-cache git build-base protobuf protobuf-dev
RUN apk add --no-cache git build-base protobuf

WORKDIR /usr/src/app

# Disable CGO
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Install air for live reload during development
# Install SQL migration tool, Go protobuf plugin, and Go gRPC plugin
RUN go install github.com/cosmtrek/air@v1.49.0 && \
  go install github.com/rubenv/sql-migrate/...@v1.6.1 && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Copy all the code and the air.toml configuration
COPY . .

# Command to run our application using air for hot-reloading
CMD ["air", "-c", ".air.toml"]

# Production stage using alpine for a smaller image
FROM golang:1.20-alpine AS production

# Set timezone if needed
ENV TZ=Asia/Bangkok
RUN apk add --update tzdata && \
  cp /usr/share/zoneinfo/$TZ /etc/localtime && \
  echo $TZ > /etc/timezone

WORKDIR /usr/src/app
COPY --from=development /usr/src/app/go.mod /usr/src/app/go.sum ./
RUN go mod download && go mod tidy
COPY --from=development /usr/src/app ./
RUN go build -o /bin/app

ENTRYPOINT [ "/bin/app" ]
