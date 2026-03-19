package main

import (
	"database/sql"
	"log"

	"github.com/robfig/cron/v3"
)

// StartCronJobs inicializa e agenda as rotinas em segundo plano
func StartCronJobs(db *sql.DB) *cron.Cron {
	c := cron.New()

	// Expressão Cron: "0 0 * * *" significa todos os dias às 00:00
	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("[CRON] Iniciando processo de decaimento de XP...")
		executeDecay(db)
	})

	if err != nil {
		log.Fatalf("[CRON] Erro ao agendar rotina: %v", err)
	}

	c.Start()
	log.Println("[CRON] Serviço de agendamento iniciado.")
	return c
}

// executeDecay processa as penalidades de inatividade
func executeDecay(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[CRON] Falha ao iniciar transação: %v\n", err)
		return
	}
	defer tx.Rollback()

	// 1. Aplica penalidade: -500 XP (mínimo 0), atualiza Level e zera a ofensiva (streak)
	// Nota de Arquitetura: A cláusula WHERE foi comentada. Ela deve ser ajustada 
	// com a modelagem exata da sua tabela de tarefas quando implementada.
	updateQuery := `
		UPDATE users 
		SET 
			current_xp = GREATEST(current_xp - 500, 0),
			current_level = (GREATEST(current_xp - 500, 0) / 1000) + 1,
			current_streak = 0
		-- WHERE id IN (SELECT user_id FROM daily_tasks WHERE status != 'COMPLETED' AND date = CURRENT_DATE)
	`
	_, err = tx.Exec(updateQuery)
	if err != nil {
		log.Printf("[CRON] Falha ao atualizar usuários: %v\n", err)
		return
	}

	// 2. Registra a auditoria no extrato (ledger)
	insertLedgerQuery := `
		INSERT INTO ledger_entries (user_id, amount, entry_type, description)
		SELECT id, 500, 'PENALTY', 'Penalidade por inatividade'
		FROM users
		-- WHERE ... (utilizar a mesma condição de filtragem da query acima)
	`
	_, err = tx.Exec(insertLedgerQuery)
	if err != nil {
		log.Printf("[CRON] Falha ao inserir no ledger: %v\n", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[CRON] Falha ao efetivar transação: %v\n", err)
		return
	}

	log.Println("[CRON] Decaimento processado com sucesso.")
}