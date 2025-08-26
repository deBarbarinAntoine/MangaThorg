FROM golang:1.24 AS builder
LABEL authors="thorgan"

WORKDIR /app

#COPY go.mod go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -trimpath -o ./mangathorg ./cmd/

FROM alpine:3.22

WORKDIR /app

# Create a non-root user for security
RUN addgroup --system --gid 1001 mangathorg
RUN adduser --system --uid 1001 mangathorg

# Create the data directory and files
RUN mkdir /app/data/
RUN echo "[]" > /app/data/users.json

# Set permissions and change user
RUN chown -R mangathorg:mangathorg /app
USER mangathorg

COPY --from=builder --chown=mangathorg:mangathorg /app/mangathorg ./mangathorg
COPY --from=builder --chown=mangathorg:mangathorg /app/config ./config
COPY --from=builder --chown=mangathorg:mangathorg /app/cache ./cache
COPY --from=builder --chown=mangathorg:mangathorg /app/assets ./assets
COPY --from=builder --chown=mangathorg:mangathorg /app/templates ./templates
#COPY --from=builder --chown=mangathorg:mangathorg /app/vendor ./vendor

EXPOSE 8080

ENTRYPOINT ["/app/mangathorg"]