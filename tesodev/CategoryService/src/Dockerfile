FROM golang:1.17-alpine AS build_base

RUN apk update && apk --no-cache add ca-certificates git && apk add tzdata
WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app ./

FROM scratch
COPY --from=build_base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV TZ Europe/Istanbul

COPY --from=build_base /bin/app /bin/app
EXPOSE 8082
ENTRYPOINT ["/bin/app"]
