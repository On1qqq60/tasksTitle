FROM golang:1.27 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/notes .

FROM postgres:17

ENV POSTGRES_USER=notes
ENV POSTGRES_PASSWORD=notes
ENV POSTGRES_DB=notes
ENV DATABASE_URL=postgres://notes:notes@127.0.0.1:5432/notes?sslmode=disable

COPY --from=build /out/notes /usr/local/bin/notes
COPY start.sh /start.sh
RUN sed -i 's/\r$//' /start.sh && chmod +x /start.sh

EXPOSE 8080

ENTRYPOINT ["/start.sh"]
