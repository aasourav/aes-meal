FROM golang:1.20-bullseye

RUN apt update

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Change .env with docker
# COPY .env.docker ./.env

ENV SERVER_ADDR=localhost
ENV SERVER_PORT=8080
ENV MONGO_URI=mongodb://localhost:27017
ENV MONGO_DATABASE=exampledb
ENV JWT_SECRET=My.Ultra.Secure.Password
ENV JWT_ACCESS_EXPIRATION_MINUTES=1440
ENV JWT_REFRESH_EXPIRATION_DAYS=7
ENV MODE=release

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]