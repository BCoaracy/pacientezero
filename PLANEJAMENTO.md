# Planejamento do projeto — CRUD Clínico

## Visão geral

Aplicação com dois CRUDs principais — **Usuário** e **Paciente** — com foco em segurança desde a fundação. O backend é uma API REST em Go, conectada a um PostgreSQL. Toda a infraestrutura (API + banco) roda via Docker Compose. O frontend ainda será definido.

---

## Stack tecnológica

| Camada | Tecnologia | Observação |
|---|---|---|
| Backend | Go (Golang) | API REST |
| Banco de dados | PostgreSQL | Container Docker |
| Infraestrutura | Docker + Docker Compose | API + banco |
| Autenticação | JWT stateless | Cookie `httpOnly` ou header `Authorization` |
| Frontend | A definir | Frameworks compatíveis com Go a avaliar |

---

## Arquitetura

### Padrão em camadas (Go)

```
Handler (HTTP)       →  lógica de roteamento, bind/validação do request, resposta HTTP
Service              →  lógica de negócio, orquestração
Repository           →  queries SQL parametrizadas ao PostgreSQL
Model                →  structs de domínio, tipos de request/response
```

### Fluxo de dados

```
Cliente HTTP
  → handler/*.go  (valida input, chama service)
  → service/*.go  (lógica de negócio)
  → repository/*.go (queries SQL — database/sql ou pgx)
  → PostgreSQL (container Docker)
```

### Autenticação — JWT stateless

- Endpoint `/login` valida credenciais e emite JWT assinado
- JWT carregado via cookie `httpOnly` ou header `Authorization: Bearer`
- Middleware de autenticação protege todas as rotas privadas
- Tabela `usuarios` com `password_hash` (bcrypt, cost factor 12)
- Sem tabela de sessão no banco

---

## Entidades e schema

### Enums PostgreSQL

```sql
CREATE TYPE role AS ENUM ('admin', 'usuario');
CREATE TYPE sexo AS ENUM ('M', 'F', 'O');
```

### Tabela: `usuarios`

| Coluna | Tipo PG | Restrições | Observação |
|---|---|---|---|
| `id` | `uuid` | PK, `gen_random_uuid()` | |
| `nome` | `varchar(255)` | NOT NULL | |
| `email` | `varchar(255)` | NOT NULL, UNIQUE | Lower-case antes de persistir |
| `password_hash` | `varchar(255)` | NOT NULL | bcrypt — **nunca retornar na API** |
| `role` | `role` | NOT NULL, default `'usuario'` | |
| `email_verified` | `boolean` | NOT NULL, default `false` | Bloquear login se false (futuro) |
| `criado_em` | `timestamp` | NOT NULL, default `now()` | |
| `atualizado_em` | `timestamp` | NOT NULL, default `now()` | Atualizar via trigger ou manualmente |

### Tabela: `pacientes`

| Coluna | Tipo PG | Restrições | Observação |
|---|---|---|---|
| `id` | `uuid` | PK, `gen_random_uuid()` | |
| `nome` | `varchar(255)` | NOT NULL | |
| `data_nascimento` | `date` | NOT NULL | Sem timezone |
| `altura_cm` | `integer` | nullable | Em centímetros (ex: 175) |
| `peso_kg` | `numeric(5,2)` | nullable | Em kg (ex: 70.50) |
| `sexo` | `sexo` | nullable | Enum: M, F, O |
| `anamnese` | `jsonb` | nullable | Estrutura flexível para crescer sem migrations |
| `criado_em` | `timestamp` | NOT NULL, default `now()` | |
| `atualizado_em` | `timestamp` | NOT NULL, default `now()` | |

> A anamnese usa `jsonb` para permitir crescimento do schema (alergias, medicamentos, histórico) sem novas migrations — basta ampliar a struct Go correspondente.

---

## Estrutura de pastas (Go)

```
pacientezero/
├── cmd/
│   └── api/
│       └── main.go               # Entrypoint — inicializa DB, router, servidor
├── internal/
│   ├── handler/                  # Bind do request, validação, resposta HTTP
│   │   ├── usuario.go
│   │   └── paciente.go
│   ├── service/                  # Lógica de negócio
│   │   ├── usuario.go
│   │   └── paciente.go
│   ├── repository/               # Queries SQL ao PostgreSQL
│   │   ├── usuario.go
│   │   └── paciente.go
│   ├── model/                    # Structs de domínio e DTOs
│   │   ├── usuario.go
│   │   └── paciente.go
│   ├── middleware/               # Auth JWT, rate limit, logging
│   │   ├── auth.go
│   │   └── rate_limit.go
│   └── db/
│       ├── db.go                 # Inicialização do pool de conexões
│       └── migrations/           # Arquivos .sql de migration
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── .env                          # Nunca commitado
└── go.mod
```

---

## Docker Compose

```yaml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  api:
    build: .
    depends_on:
      - db
    ports:
      - "8080:8080"
    env_file:
      - .env

volumes:
  pgdata:
```

---

## Segurança por camada

### Network / Headers
- `Strict-Transport-Security` em produção
- `X-Frame-Options: DENY`
- `X-Content-Type-Options: nosniff`
- `Referrer-Policy: strict-origin-when-cross-origin`
- CORS restritivo configurado no middleware

### Handler / API
- Todo handler protegido pelo middleware de autenticação JWT
- Validação de UUID antes de qualquer query
- Respostas de erro sem stack trace em produção
- Rate limiting por IP
- Bind e validação de input antes de chamar o service

### Repository / DB
- Queries 100% parametrizadas via `database/sql` ou `pgx` — sem SQL injection
- SSL obrigatório na connection string em produção (`sslmode=require`)
- `password_hash` nunca incluído em structs de response
- Variáveis de ambiente validadas na inicialização

---

## Variáveis de ambiente

### `.env.example` (commitado)
```bash
# Banco de dados
DB_HOST=db
DB_PORT=5432
DB_USER=
DB_PASSWORD=
DB_NAME=

# API
PORT=8080
JWT_SECRET=
```

### `.env` (nunca commitado — adicionar ao .gitignore)
```bash
DB_HOST=db
DB_PORT=5432
DB_USER=pacientezero
DB_PASSWORD=senha_forte_aqui
DB_NAME=pacientezero_db

PORT=8080
JWT_SECRET=<gerado com: openssl rand -base64 32>
```

---

## Fases do projeto

### Fase 1 — Setup & infraestrutura
- [x] Criar repositório no GitHub
- [x] Inicializar módulo Go (`go mod init`)
- [x] Criar `docker-compose.yml` com serviços `db` e `api`
- [x] Criar `Dockerfile` multi-stage para o binário Go
- [ ] Validar que `docker compose up` sobe o banco e a API
- [x] Criar estrutura de pastas (`cmd/`, `internal/`)
- [x] Configurar linter (`golangci-lint`)
- [x] Criar `.env.example` e adicionar `.env` ao `.gitignore`

### Fase 2 — Model: schema, migrations e structs
- [x] Escrever migrations SQL iniciais (enums + tabelas `usuarios` e `pacientes`)
- [x] Definir structs de domínio em `internal/model/`
- [x] Criar script de seed para dev/test
- [x] Decidir biblioteca de migration (golang-migrate)

### Fase 3 — Repository: acesso ao banco
- [x] Criar `internal/db/db.go` — pool de conexões com `pgx` ou `database/sql`
- [x] Implementar `internal/repository/usuario.go` — CRUD com queries parametrizadas
- [x] Implementar `internal/repository/paciente.go`
- [x] Escrever testes de integração dos repositórios (banco real em Docker)

### Fase 4 — Service: lógica de negócio
- [x] Implementar `internal/service/usuario.go` — hash de senha, validações de negócio
- [x] Implementar `internal/service/paciente.go`
- [x] Criar `internal/model/errors.go` — erros tipados (`NotFound`, `Unauthorized`, etc.)
- [x] Escrever testes unitários dos services (mock dos repositories)

### Fase 5 — Handler: API REST e autenticação
- [x] Escolher router HTTP (net/http padrão — Go 1.22+)
- [x] Implementar handlers CRUD para `usuarios` e `pacientes`
- [x] Implementar endpoint `/login` — valida credenciais, emite JWT
- [x] Criar `internal/middleware/auth.go` — valida JWT em rotas protegidas
- [x] Criar `internal/middleware/rate_limit.go`
- [x] Configurar headers de segurança (X-Frame-Options, nosniff, Referrer-Policy)
- [x] Configurar CORS restritivo
- [x] Escrever testes de integração dos handlers

### Fase 6 — Frontend (a definir)
- [ ] Avaliar frameworks compatíveis com Go (HTMX + Templ, Next.js standalone, SvelteKit, etc.)
- [ ] Decidir abordagem (SSR servido pelo Go, SPA separada, ou híbrido)
- [ ] Integrar ao Docker Compose

### Fase 7 — CI/CD + hardening
- [ ] Configurar GitHub Actions: lint + testes em todo PR
- [ ] Adicionar análise de vulnerabilidades (`govulncheck`)
- [ ] Validar headers de segurança em produção
- [ ] Smoke test E2E — CRUD completo

---

*Documento em evolução — atualizar conforme o projeto avança.*
