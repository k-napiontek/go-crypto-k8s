FROM golang:tip-alpine3.23 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o server

FROM scratch

COPY --from=build /app/server /server

ENTRYPOINT ["/server"]