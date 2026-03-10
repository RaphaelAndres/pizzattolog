-- Migrations iniciais do PizzattoLog
-- Executado automaticamente pelo Docker no primeiro boot

CREATE DATABASE IF NOT EXISTS pizzattolog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE pizzattolog;

-- As tabelas são criadas via AutoMigrate do GORM.
-- Este script garante apenas que o banco existe com o charset correto.

-- Você pode adicionar aqui seeds de desenvolvimento:
-- INSERT INTO usuarios (nome, email, senha_hash, role, ativo) VALUES (...)
