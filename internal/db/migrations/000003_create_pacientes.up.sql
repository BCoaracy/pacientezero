CREATE TABLE pacientes (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    nome             VARCHAR(255) NOT NULL,
    data_nascimento  DATE         NOT NULL,
    altura_cm        INTEGER,
    peso_kg          NUMERIC(5,2),
    sexo             sexo,
    anamnese         JSONB,
    criado_em        TIMESTAMP    NOT NULL DEFAULT NOW(),
    atualizado_em    TIMESTAMP    NOT NULL DEFAULT NOW()
);