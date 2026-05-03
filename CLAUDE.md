# CLAUDE.md — pacientezero

Instruções para o Claude Code trabalhar neste projeto.

## Visão geral

API REST clínica em Go com dois CRUDs principais: **Usuário** e **Paciente**. Infraestrutura via Docker Compose (API + PostgreSQL). Frontend ainda não definido.

## Stack

- **Backend:** Go (Golang)
- **Banco:** PostgreSQL 16 (container Docker)
- **Infraestrutura:** Docker + Docker Compose
- **Autenticação:** JWT stateless, bcrypt para senhas
- **Frontend:** a definir

## Estrutura de pastas

```
cmd/api/main.go          # entrypoint
internal/handler/        # HTTP: bind, validação, resposta
internal/service/        # lógica de negócio
internal/repository/     # queries SQL parametrizadas
internal/model/          # structs de domínio e DTOs
internal/middleware/      # auth JWT, rate limit
internal/db/             # pool de conexões e migrations
```

## Como rodar

```bash
# Subir banco + API
docker compose up

# Apenas o banco (útil durante desenvolvimento)
docker compose up db

# Build local do binário
go build ./cmd/api/...

# Testes
go test ./...
```

## Convenções obrigatórias

- Queries SQL sempre parametrizadas — nunca concatenar strings para montar SQL
- `password_hash` nunca incluído em structs de response — usar DTOs separados
- Variáveis de ambiente lidas de `.env` — nunca valores hardcoded
- Erros tipados em `internal/model/errors.go` — não retornar `errors.New` com strings soltas nos handlers
- Respostas de erro em produção sem stack trace — apenas código e mensagem genérica

## Segurança

- Todo handler privado protegido pelo middleware JWT (`internal/middleware/auth.go`)
- UUIDs validados antes de qualquer query ao banco
- Headers de segurança configurados no middleware (HSTS, X-Frame-Options, nosniff, CORS restritivo)
- Rate limiting por IP em todas as rotas

## Banco de dados

- Migrations em SQL puro em `internal/db/migrations/`
- Enums PostgreSQL: `role` (`admin`, `usuario`) e `sexo` (`M`, `F`, `O`)
- Campo `anamnese` é `jsonb` — permite evolução do schema sem novas migrations
- SSL obrigatório na connection string em produção (`sslmode=require`)

## Variáveis de ambiente

Copiar `.env.example` para `.env` antes de rodar. Nunca commitar `.env`.

| Variável | Descrição |
|---|---|
| `DB_HOST` | Host do PostgreSQL (padrão: `db` no Docker) |
| `DB_PORT` | Porta do PostgreSQL (padrão: `5432`) |
| `DB_USER` | Usuário do banco |
| `DB_PASSWORD` | Senha do banco |
| `DB_NAME` | Nome do banco |
| `PORT` | Porta da API (padrão: `8080`) |
| `JWT_SECRET` | Segredo para assinar tokens JWT |

## Planejamento completo

Ver [PLANEJAMENTO.md](PLANEJAMENTO.md) para fases do projeto, schema detalhado das tabelas e decisões de arquitetura.
