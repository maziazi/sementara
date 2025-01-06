# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
# Build binary dengan CGO disabled untuk memastikan binary statis
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# Run stage
FROM alpine:3.18
WORKDIR /root/
# Salin binary 'main' dari image builder ke image ini
COPY --from=builder /app/main .

# Periksa apakah binary 'main' ada
RUN ls -l /root/

EXPOSE 8080
CMD ["/root/main"]
## Build stage
#FROM golang:1.23 AS builder
#WORKDIR /app
#COPY . .
#RUN go mod tidy
#RUN go build -o /app/main ./cmd  # Menyusun aplikasi dari subdirektori cmd
#
## Run stage
#FROM alpine:3.18
#WORKDIR /root/
## Salin binary 'main' dari image builder ke image ini
#COPY --from=builder /app/main .
#
## Pastikan file main dapat dijalankan
#RUN ls -l /root/   # Periksa apakah file 'main' ada di direktori ini
#
#EXPOSE 8080
#CMD ["/root/main"]


## Build stage
#FROM golang:1.23 AS builder
#WORKDIR /app
#COPY . .
#RUN go mod tidy
#RUN go build -o main ./cmd
#
## Run stage
#FROM alpine:3.18
#WORKDIR /root/
#COPY --from=builder /app/main .
#EXPOSE 8080
#CMD ["./main"]


#FROM golang:1.21
#
## Set working directory
#WORKDIR /app
#
## Copy go.mod dan go.sum lalu download dependencies
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copy seluruh kode ke dalam container
#COPY golang .
#
## Compile aplikasi
#RUN go build -o main .
#
## Jalankan aplikasi
#CMD ["/app/main"]
