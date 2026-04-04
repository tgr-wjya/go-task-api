FROM golang:1.26-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

FROM scratch

COPY --from=build /app/server /server

EXPOSE 8080

CMD ["/server"]