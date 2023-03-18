FROM golang:alpine

LABEL maintainer="Mr Beyond"


# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app

WORKDIR /app

COPY . .

COPY .env .

RUN go get -d -v ./...

RUN go install -v ./...

# Build the Go app
RUN go build -o /build

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]
