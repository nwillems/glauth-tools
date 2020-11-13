FROM golang:latest AS build

# Build and jazz
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build

FROM scratch

COPY --from=build /app/glauth-tools /glauth-tools
ENTRYPOINT ["/glauth-tools"]
