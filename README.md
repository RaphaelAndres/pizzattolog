# 🍕 PizzattoLog — Gerenciador de Licenças

> Sistema de gerenciamento de licenças ambientais, comerciais e regulatórias para a **PizzattoLog**.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat-square&logo=mysql)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker)

---

## 📋 Sobre o Projeto

O **PizzattoLog** é uma plataforma web para armazenamento, visualização e controle de validade de licenças operacionais (ambientais, de polícia civil, sanitárias, etc). O sistema envia alertas automáticos quando as licenças estão próximas do vencimento, evitando multas e irregularidades.

### Funcionalidades principais

- 📂 **Upload de licenças** em PDF/imagem com armazenamento via MinIO
- 📅 **Controle de validade** com dashboard visual e linha do tempo
- 🔔 **Alertas automáticos** por e-mail e notificação in-app (30, 15 e 7 dias antes do vencimento)
- 🔐 **Autenticação segura** com JWT (login, registro, roles)
- 👥 **Multi-usuário** com diferentes níveis de acesso (admin, gestor, visualizador)
- 📊 **Dashboard** com resumo de licenças ativas, vencidas e próximas ao vencimento

---

## 🏗️ Arquitetura

```
pizzattolog/
├── backend/           # API REST em Go (Gin + GORM)
│   ├── cmd/server/    # Ponto de entrada
│   ├── internal/
│   │   ├── auth/      # JWT, bcrypt
│   │   ├── handlers/  # HTTP handlers
│   │   ├── middleware/ # Auth, CORS, logging
│   │   ├── models/    # Entidades (GORM)
│   │   ├── repository/ # Camada de dados
│   │   └── services/  # Regras de negócio
│   ├── migrations/    # Migrations SQL
│   └── tests/         # Testes unitários e integração
│
├── frontend/          # SPA em React + TypeScript + Vite
│   └── src/
│       ├── components/ # Componentes reutilizáveis
│       ├── pages/     # Telas (Login, Dashboard, Licenças)
│       ├── hooks/     # Custom hooks
│       ├── services/  # Chamadas à API
│       └── types/     # Tipos TypeScript
│
├── docker-compose.yml # Orquestração local
└── docs/              # Documentação e planejamento
```

---

## 🚀 Rodando localmente com Docker

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) instalado
- [Docker Compose](https://docs.docker.com/compose/install/) v2+

### 1. Clone o repositório

```bash
git clone https://github.com/pizzattolog/licencas.git
cd licencas
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
# Edite o .env conforme necessário
```

### 3. Suba todos os serviços

```bash
docker compose up --build
```

### 4. Acesse a aplicação

| Serviço     | URL                          |
|-------------|------------------------------|
| Frontend    | http://localhost:3002        |
| API (Go)    | http://localhost:8081        |
| MinIO UI    | http://localhost:9101        |
| MySQL       | localhost:3307               |

> **Credenciais padrão (desenvolvimento):**
> - Usuário admin: `admin@pizzattolog.com.br` / `Admin@123`
> - MinIO: `minioadmin` / `minioadmin`

---

## ⚙️ Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Aplicação
APP_ENV=development
APP_PORT=8080

# Banco de dados
DB_HOST=mysql
DB_PORT=3306
DB_USER=pizzatto
DB_PASSWORD=pizzatto123
DB_NAME=pizzattolog

# JWT
JWT_SECRET=sua-chave-super-secreta-aqui
JWT_EXPIRY_HOURS=24

# MinIO (armazenamento de arquivos)
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=licencas
MINIO_USE_SSL=false

# E-mail (alertas)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@pizzattolog.com.br
SMTP_PASSWORD=sua-senha-smtp

# Frontend
VITE_API_URL=http://localhost:8080/api/v1
```

---

## 🧪 Rodando os Testes

### Backend (Go)

```bash
cd backend
go test ./... -v
```

### Com cobertura

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Frontend

```bash
cd frontend
npm test
```

---

## 📡 Endpoints da API

### Autenticação

| Método | Rota                | Descrição              |
|--------|---------------------|------------------------|
| POST   | `/api/v1/auth/login`    | Login do usuário       |
| POST   | `/api/v1/auth/register` | Registro (admin only)  |
| POST   | `/api/v1/auth/refresh`  | Renovar token JWT      |
| POST   | `/api/v1/auth/logout`   | Logout                 |

### Licenças

| Método | Rota                        | Descrição                     |
|--------|-----------------------------|-------------------------------|
| GET    | `/api/v1/licencas`          | Listar todas as licenças      |
| POST   | `/api/v1/licencas`          | Criar nova licença + upload   |
| GET    | `/api/v1/licencas/:id`      | Detalhes de uma licença       |
| PUT    | `/api/v1/licencas/:id`      | Atualizar licença             |
| DELETE | `/api/v1/licencas/:id`      | Remover licença               |
| GET    | `/api/v1/licencas/:id/arquivo` | Download do arquivo        |

### Dashboard

| Método | Rota                    | Descrição                        |
|--------|-------------------------|----------------------------------|
| GET    | `/api/v1/dashboard`     | Resumo geral                     |
| GET    | `/api/v1/alertas`       | Licenças próximas ao vencimento  |

### Usuários

| Método | Rota                    | Descrição                |
|--------|-------------------------|--------------------------|
| GET    | `/api/v1/usuarios`      | Listar usuários (admin)  |
| PUT    | `/api/v1/usuarios/:id`  | Atualizar usuário        |
| DELETE | `/api/v1/usuarios/:id`  | Remover usuário (admin)  |

---

## 🗂️ Modelo de Dados

### Licença

```json
{
  "id": 1,
  "nome": "Licença Ambiental IBAMA",
  "tipo": "ambiental",
  "orgao_emissor": "IBAMA",
  "numero": "LA-2024-001234",
  "data_emissao": "2024-01-15",
  "data_validade": "2026-01-15",
  "status": "ativa",
  "arquivo_url": "https://minio.../licencas/la-2024-001234.pdf",
  "tags": ["ambiental", "federal"],
  "criado_em": "2024-01-20T10:30:00Z",
  "atualizado_em": "2024-01-20T10:30:00Z"
}
```

---

## 🔔 Sistema de Alertas

O sistema verifica automaticamente as licenças vencendo e dispara notificações:

- **30 dias antes**: Alerta de atenção (e-mail + badge no dashboard)
- **15 dias antes**: Alerta de urgência (e-mail + notificação in-app)
- **7 dias antes**: Alerta crítico (e-mail para todos os gestores)
- **Vencida**: Status automático atualizado para `vencida`

---

## 🛠️ Stack Tecnológica

| Camada       | Tecnologia               | Versão  |
|--------------|--------------------------|---------|
| Backend      | Go + Gin                 | 1.22    |
| ORM          | GORM                     | v2      |
| Banco        | MySQL                    | 8.0     |
| Armazenamento| MinIO                    | latest  |
| Frontend     | React + TypeScript + Vite| 18 / 5  |
| Estilo       | Tailwind CSS             | v3      |
| Auth         | JWT (golang-jwt)         | v5      |
| Testes (Go)  | testify + httptest       | —       |
| Testes (JS)  | Vitest + Testing Library | —       |
| CI/CD        | GitHub Actions           | —       |
| Container    | Docker + Compose         | v2      |

---

## 📄 Licença

Projeto proprietário — © 2025 PizzattoLog. Todos os direitos reservados.
