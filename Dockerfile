FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/web ./cmd/web \
 && CGO_ENABLED=0 go build -o /out/indexer ./cmd/indexer
# Bake the timetable into the image so a cold start only has to read a 1 MB file.
RUN /out/indexer -out /out/index.gob

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/web /out/index.gob /app/
ENV DATA_PATH=/app/index.gob GIN_MODE=release
EXPOSE 8000
ENTRYPOINT ["/app/web"]
