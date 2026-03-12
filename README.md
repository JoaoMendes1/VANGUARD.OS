# 🌐 VANGUARD.OS

Um Ecossistema de Produtividade Gamificado de alta performance, construído com Next.js, React Native, Go (Golang) e PostgreSQL.

---

## 🎯 Visão do Produto

O VANGUARD.OS não é um simples "to-do list". É um Sistema Operacional de Vida (*Life OS*) projetado para resolver o problema da inconsistência humana. Ele substitui listas de tarefas monótonas por uma simulação de evolução cibernética e biológica com riscos reais. 

Utilizando gatilhos neurológicos como **Aversão à Perda** e **Recompensas Variáveis**, o sistema obriga o usuário a manter o foco em seus objetivos diários, penalizando severamente a procrastinação.

## ✨ Features (MVP - Fase 1)

* **Motor de Atributos:** Toda ação concluída gera XP para 5 pilares biológicos/neurais: *Kinetic (Corpo), Neural (Foco), Core (Saúde Base), Sync (Social)* e *Logic (Estratégia)*.
* **Protocolos (Hábitos):** Ações binárias diárias focadas em *Streaks* (ofensivas). Falhar um protocolo reseta a ofensiva e engatilha o Sistema de Deterioração.
* **Operações (Tarefas):** Missões únicas com prazos e prioridades (Alta, Média, Baixa).
* **Subsistema de Hobbies (Micro-Skills):** Sistema de *Tags* que permite aos usuários evoluir níveis em habilidades específicas (Ex: Nível 10 em Corrida, Nível 5 em Leitura) ganhando títulos e emblemas.
* **Sistema de Deterioração (Decay):** Falhar nos protocolos diários até a meia-noite gera uma penalidade massiva de XP, ativando a aversão à perda.
* **Balancete de XP (The Ledger):** Um extrato financeiro transparente mostrando o fluxo exato de XP (Ganho, Gasto e Perdido por Deterioração).
* **Loja do Sistema (Override Store):** O usuário ganha *V-Credits* para comprar permissões do mundo real (ex: Sessão de TV sem culpa) ou itens do jogo (ex: Escudo de 24h contra perda de XP).
* **Sprints (Operações Especiais):** Desafios de 7, 21 ou 30 dias (Ex: *Monk Mode Sprint*) com acompanhamento de taxa de sucesso diária.

## 🛠️ Tech Stack (Arquitetura Monorepo)

O projeto adota uma arquitetura robusta separando a inteligência central (Backend) das interfaces de consumo (Clients).

* **Gerenciamento de Workspace:** Turborepo
* **Frontend Web:** Next.js (React), Tailwind CSS, Lucide Icons, Chart.js (Radar/Ledger).
* **Frontend Mobile:** React Native (Expo).
* **Backend API:** Go (Golang) utilizando a Standard Library.
* **Banco de Dados:** PostgreSQL (Relacional) modelado com UUIDs rígidos.
* **Infraestrutura Local:** Docker & Docker Compose.

## 🗺️ Roadmap

Nosso progresso detalhado de desenvolvimento é gerenciado no documento oficial do projeto.  
👉 [Consulte o ROADMAP.md completo aqui](./ROADMAP.md)

## 🏁 Como Rodar (Setup Local)

*(Esta seção será documentada assim que o setup do Docker e do Monorepo forem inicializados na Fase 1.1 do Roadmap).*
