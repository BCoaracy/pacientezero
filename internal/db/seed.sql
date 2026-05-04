-- Seed de desenvolvimento — não executar em produção
-- Rodar manualmente: psql $DATABASE_URL -f internal/db/seed.sql

INSERT INTO usuarios (nome, email, password_hash, role) VALUES
    ('Admin Dev', 'admin@dev.local', '$2a$12$placeholder_hash_admin', 'admin'),
    ('Usuário Dev', 'usuario@dev.local', '$2a$12$placeholder_hash_user', 'usuario')
ON CONFLICT (email) DO NOTHING;

INSERT INTO pacientes (nome, data_nascimento, sexo, anamnese) VALUES
    ('Paciente Teste', '1990-05-15', 'M', '{"alergias": ["dipirona"], "medicamentos": []}'),
    ('Paciente Dois',  '1985-11-22', 'F', '{"alergias": [], "medicamentos": ["losartana 50mg"]}')
ON CONFLICT DO NOTHING;