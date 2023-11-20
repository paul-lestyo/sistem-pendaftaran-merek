FROM golang:1.19.3-alpine as Build
RUN apk --no-cache add ca-certificates
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .


ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags='-s -w' -o main .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/views/ /views/
COPY --from=build /build/assets/ /assets/
COPY --from=Build ["/build/main", "/build/.env", "/"]

#we tell docker what to run when this image is run and run it as executable.
ENTRYPOINT [ "/main" ]