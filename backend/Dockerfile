# ── Build stage ────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Copia módulos e baixa dependências (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# ── Run stage ──────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Certificados e timezone
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=America/Sao_Paulo

# Copia binário do builder
COPY --from=builder /app/server .

# Usuário não-root para segurança
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
