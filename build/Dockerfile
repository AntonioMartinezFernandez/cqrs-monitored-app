#####################################
#          Builder image            #
#####################################

FROM golang:1.24.1 AS builder

WORKDIR /cqrs-monitored-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Download modules and delete the ssh key
RUN go mod download

WORKDIR /cqrs-monitored-app/cmd/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o cqrs-monitored-app

#####################################
#        Deployable image           #
#####################################

FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S cqrs-monitored-app -G appgroup -h /home/cqrs-monitored-app
USER cqrs-monitored-app

WORKDIR /home/cqrs-monitored-app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder cqrs-monitored-app/cmd/app/cqrs-monitored-app .
COPY --from=builder cqrs-monitored-app/schemas ./schemas

ENTRYPOINT ["/home/cqrs-monitored-app/cqrs-monitored-app"]
