FROM public.ecr.aws/docker/library/golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o weather-api ./cmd

FROM public.ecr.aws/docker/library/golang:latest 
WORKDIR /app

COPY --from=builder /app/weather-api .

EXPOSE 8080

CMD ["./weather-api"]