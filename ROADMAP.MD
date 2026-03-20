# 🗺️ VANGUARD.OS - Roadmap de Desenvolvimento

Este documento descreve a visão técnica de alto nível e as etapas granulares de construção do VANGUARD.OS. Ele serve como guia principal para qualquer desenvolvedor ingressando no projeto.

---

## ✅ Fase 1: Fundação & Infraestrutura (Backend) - CONCLUÍDO

O objetivo desta fase é colocar a estrutura base do projeto de pé, garantindo que o banco de dados e as tecnologias conversem entre si.

### 1.1: Setup do Repositório
- [x] Levantamento de Requisitos e Definição de Escopo (PRD concluído).
- [x] Validação Visual de UI/UX (Mockup interativo construído em React).
- [x] Inicializar o Monorepo com Turborepo.
- [x] Configurar padrão de commits (Husky, Commitlint).
- [x] Criar arquivo `docker-compose.yml` para rodar o PostgreSQL localmente.

### 1.2: Modelagem de Dados e Go Server
- [x] Estruturar o esquema SQL (Tabelas: `users`, `hobbies`, `protocols`, `operations`, `ledger_entries`, `sprints`).
- [x] Inicializar o servidor Backend em Go (`go mod init`).
- [x] Configurar conexão do Go com o banco PostgreSQL.
- [x] Criar a estrutura de pastas do Backend (Controllers, Models, Routes).

### 1.3: Sistema de Autenticação (MVP)
- [x] Criar API de Cadastro de Usuário (Criptografia de senha com bcrypt).
- [x] Criar API de Login com geração de JWT (JSON Web Token).
- [x] Implementar Middleware de autenticação no Go para proteger rotas.

---

## 🚧 Fase 2: Motor de Regras (Core API) - EM ANDAMENTO

Nesta fase, programaremos a matemática do jogo. Todo o ganho, perda de XP e lógica de hábitos viverá exclusivamente no servidor (Backend) para evitar fraudes (*cheats*).

### 2.1: APIs de Operações e Protocolos (O PRÓXIMO ALVO - ISSUE #10)
- [ ] CRUD de Protocolos (Criar, Listar, Concluir hábito, Resetar Streak).
- [ ] CRUD de Operações (Criar tarefa, Listar por data, Concluir, Deletar).
- [ ] Rota para buscar o Estado atual do Jogador (XP, Level, Créditos, Atributos).

### 2.2: O Motor de XP e Economia
- [x] **Issue #8:** Implementar função global de distribuição de recompensa (Calcula XP, Atualiza V-Credits, Atualiza Ledger).
- [ ] Implementar lógica de Tag de Hobbies (Se operação tem *tag*, envia XP para a tabela do Hobby).
- [x] **Issue #9:** Desenvolver o *Cronjob* da Meia-Noite (Worker que roda 00:00 para verificar protocolos falhos e aplicar punição de *Decay*).

---

## 💻 Fase 3: Dashboard Web (Frontend Next.js) - 📅 PLANEJADO

Nesta fase, pegaremos o código React da nossa prévia validada e o conectaremos com os dados reais vindos do servidor Go.

### 3.1: Fundação Front-end
- [ ] Inicializar App Next.js com Tailwind CSS.
- [ ] Criar contexto global de estado (ou Zustand) para guardar os dados do Player logado.
- [ ] Implementar telas de Login e Registro.

### 3.2: Integração do Painel Principal
- [ ] Conectar os cards de Protocolos e Operações com os endpoints da API.
- [ ] Dinamizar o Radar de Atributos SVG com os dados reais do banco.
- [ ] Dinamizar as barras de progresso do Subsistema de Hobbies.
- [ ] Popular o Balancete (Ledger) com o histórico de transações real do usuário.

---

## 📱 Fase 4: Expansão Mobile (React Native) - 📅 PLANEJADO

Levar a experiência do VANGUARD.OS para o bolso do usuário, focando na facilidade de registrar os hábitos no meio da rotina.

- [ ] Inicializar projeto Expo no Monorepo.
- [ ] Compartilhar lógicas de consumo de API entre Web e Mobile.
- [ ] Criar interface mobile simplificada (Foco exclusivo em marcar as tarefas como concluídas rapidamente).
- [ ] Implementar Push Notifications locais (Aviso de Deterioração iminente às 22h).

---

## 🌌 Fase 5: Melhorias & Extras (Backlog) - 📦 FUTURO

Recursos que estão fora da Fase 1, mas já foram planejados na visão de produto de longo prazo.

- [ ] **Integração Social:** Sistema de Guildas e Batalhas entre jogadores.
- [ ] **The Cortex (AI):** Integração com Inteligência Artificial para análise de produtividade e geração automática de missões para combater pontos fracos.
- [ ] **WhatsApp Bot:** Permite concluir hábitos e adicionar tarefas enviando mensagens de áudio ou texto via WhatsApp.
- [ ] **Pagamentos (Stripe):** Preparação de *Tiers* pagos (SaaS) oferecendo skins exclusivas de avatares e relatórios avançados de dados.