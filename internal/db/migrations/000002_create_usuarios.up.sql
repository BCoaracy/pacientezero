CREATE TABLE usuarios (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    nome           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hash  VARCHAR(255) NOT NULL,
    role           role         NOT NULL DEFAULT 'usuario',
    email_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    criado_em      TIMESTAMP    NOT NULL DEFAULT NOW(),
    atualizado_em  TIMESTAMP    NOT NULL DEFAULT NOW()
);