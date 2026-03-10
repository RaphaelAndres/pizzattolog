# 📋 Planejamento de Requisitos — PizzattoLog Licenças

## 1. Visão Geral

**Produto:** Sistema de Gerenciamento de Licenças Operacionais  
**Cliente:** PizzattoLog  
**Stack:** Go (backend) + React (frontend) + MySQL + MinIO  

---

## 2. Personas e Usuários

| Persona      | Descrição                                                   | Permissões                          |
|--------------|-------------------------------------------------------------|-------------------------------------|
| **Admin**    | TI ou gestor geral, configura usuários e visualiza tudo     | Full access                         |
| **Gestor**   | Responsável por subir e gerenciar licenças                  | CRUD de licenças, visualizar alertas|
| **Visualizador** | Colaborador que apenas consulta o status das licenças   | Somente leitura                     |

---

## 3. Requisitos Funcionais

### RF-01 — Autenticação e Autorização
- RF-01.1: Login com e-mail e senha (JWT)
- RF-01.2: Logout (invalida token no client)
- RF-01.3: Refresh de token antes da expiração
- RF-01.4: Registro de novos usuários (restrito ao Admin)
- RF-01.5: Gerenciamento de roles (Admin, Gestor, Visualizador)
- RF-01.6: Senha com hash bcrypt

### RF-02 — Gerenciamento de Licenças
- RF-02.1: Cadastro de licença com: nome, tipo, órgão emissor, número, data de emissão, data de validade
- RF-02.2: Upload do arquivo da licença (PDF ou imagem) para o MinIO
- RF-02.3: Listagem de licenças com filtros por: tipo, status, órgão, intervalo de data
- RF-02.4: Visualização de detalhes da licença
- RF-02.5: Edição de dados da licença
- RF-02.6: Exclusão lógica (soft delete) de licença
- RF-02.7: Download do arquivo original da licença
- RF-02.8: Pré-visualização do PDF/imagem inline
- RF-02.9: Histórico de alterações por licença

### RF-03 — Tipos de Licença (categorias iniciais)
- Ambiental (IBAMA, SEMA, etc.)
- Polícia Civil / Segurança
- Sanitária (ANVISA, Vigilância Sanitária)
- Bombeiros
- Prefeitura / Alvará
- Outro (campo livre)

### RF-04 — Sistema de Alertas e Notificações
- RF-04.1: Cron job diário que verifica licenças próximas ao vencimento
- RF-04.2: Envio de e-mail para gestores 30, 15 e 7 dias antes do vencimento
- RF-04.3: Badge/contador de alertas no header do dashboard
- RF-04.4: Página de alertas com todas as licenças em estado crítico
- RF-04.5: Status automático: `ativa` → `próxima_vencimento` → `vencida`

### RF-05 — Dashboard
- RF-05.1: Contador de licenças por status (ativa, próxima do vencimento, vencida)
- RF-05.2: Lista das próximas licenças a vencer (top 5)
- RF-05.3: Gráfico por tipo de licença
- RF-05.4: Linha do tempo de vencimentos (calendário simplificado)

### RF-06 — Busca e Filtros
- RF-06.1: Busca global por nome, número ou órgão emissor
- RF-06.2: Filtro por tipo, status e intervalo de data
- RF-06.3: Ordenação por nome, validade ou data de criação

---

## 4. Requisitos Não Funcionais

| ID     | Requisito                                                                |
|--------|--------------------------------------------------------------------------|
| RNF-01 | API responde em < 300ms para operações comuns (sem upload)               |
| RNF-02 | Arquivos de até 50MB por upload                                          |
| RNF-03 | Autenticação via JWT com expiração de 24h                                |
| RNF-04 | Toda comunicação via HTTPS em produção                                   |
| RNF-05 | Cobertura de testes ≥ 70% no backend                                    |
| RNF-06 | Deploy via Docker Compose (dev) e Docker Swarm/K8s (prod)               |
| RNF-07 | Logs estruturados (JSON) com nível configurável                         |
| RNF-08 | Variáveis sensíveis exclusivamente via `.env` / secrets                 |

---

## 5. Modelo de Dados (MySQL)

### Tabela: `usuarios`
```sql
id, nome, email, senha_hash, role, ativo, criado_em, atualizado_em
```

### Tabela: `licencas`
```sql
id, nome, tipo, orgao_emissor, numero, descricao,
data_emissao, data_validade, status, arquivo_key (MinIO),
arquivo_nome, arquivo_tamanho, criado_por (FK usuarios),
deletado_em (soft delete), criado_em, atualizado_em
```

### Tabela: `historico_licencas`
```sql
id, licenca_id, usuario_id, acao (criou/editou/deletou), dados_anteriores (JSON), criado_em
```

### Tabela: `alertas_enviados`
```sql
id, licenca_id, tipo_alerta (30d/15d/7d), enviado_em, destinatarios (JSON)
```

---

## 6. Roadmap

### 🟢 MVP (Este documento)
- [x] Estrutura do projeto + Docker Compose
- [x] Auth completo (login, JWT, roles)
- [x] CRUD de licenças + upload MinIO
- [x] Dashboard com resumo
- [x] Alertas por e-mail (cron)
- [x] Interface React responsiva
- [x] Testes unitários básicos

### 🔵 Versão 1.1
- [ ] Notificações in-app em tempo real (WebSocket)
- [ ] Histórico de alterações por licença
- [ ] Exportação de relatórios em PDF
- [ ] Pré-visualização inline de PDFs

### 🟣 Versão 2.0
- [ ] Multi-empresa (multi-tenant)
- [ ] Integração com sistemas de e-mail corporativo
- [ ] App mobile (React Native)
- [ ] API pública documentada (Swagger)
- [ ] Auditoria completa de ações
