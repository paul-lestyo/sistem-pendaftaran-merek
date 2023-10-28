FROM golang:1.19.3-alpine as Build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags='-s -w' -o main .

FROM scratch

COPY --from=Build ["/build/","/build/main", "/build/.env", "/"]

#we tell docker what to run when this image is run and run it as executable.
ENTRYPOINT [ "/main" ]