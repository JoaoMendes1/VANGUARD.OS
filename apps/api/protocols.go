package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// 1. CRIAR UM HÁBITO (POST)
func CreateProtocolHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Recupera o ID do usuário que o nosso Middleware guardou no bolso (Context)
		userID := r.Context().Value(userContextKey).(string)

		var p Protocol
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		// Insere no banco. Note que injetamos o userID direto, sem o usuário precisar mandar!
		query := `INSERT INTO protocols (user_id, title, attribute_type) VALUES ($1, $2, $3) RETURNING id`
		var newProtocolID string

		err = db.QueryRow(query, userID, p.Title, p.AttributeType).Scan(&newProtocolID)
		if err != nil {
			http.Error(w, "Failed to create protocol", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Protocol initialized!",
			"id":      newProtocolID,
		})

	}

}

// 2. LISTAR OS HÁBITOS DO USUÁRIO (GET)
func GetProtocolsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return 
		}

		// Descobre de quem é o token 
		userID := r.Context().Value(userContextKey).(string)
		
		//Busca APENAS os hábitos deste usuário específico 
		query := `SELECT id, title, streak_count, attribute_type, is_active FROM protocols WHERE  user_id = $1`
		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Cria uma lista vazia de protocolos 
		var protocols []Protocol

		// Faz um loop (for) lendo cada linha que o banco devolveu
		for rows.Next() {
			var p Protocol
			err := rows.Scan(&p.ID, &p.Title, &p.StreakCount, &p.AttributeType, &p.IsActive)
			if err != nil {
				continue // Se der erro numa linha, pula pra próxima
			}
			protocols = append(protocols, p) // Adiciona o hábito na lista 
		}

		// Se alista estiver vazia (nil), envia um array vazio [] ao invés de null para o React
		if protocols == nil {
			protocols = []Protocol{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(protocols)
	}
}

// Estruturra para receber o ID do protocolo que será concluído
type CompleteProtocolPayload struct {
	ProtocolID string `json:"protocol_id"`
}

// 3. CONCLUIR UM HÁBITO (POST /protocols/complete)
func CompleteProtocolHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.Context().Value(userContextKey).(string)

		var payload CompleteProtocolPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil || payload.ProtocolID == "" {
			http.Error(w, "Invalid payload. 'protocol_id' is required.", http.StatusBadRequest)
			return
		}

		// 1. Inicia a Transação
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database transaction failed", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// 2. Atualiza o Protocolo (Anti-Cheat Integrado no SQL)
		// Só atualiza se last_completed_at for nulo ou menor que a data atual (hoje)
		var updatedProtocolID string
		updateProtocolQuery := `
			UPDATE protocols 
			SET streak_count = streak_count + 1, last_completed_at = CURRENT_DATE 
			WHERE id = $1 AND user_id = $2 AND (last_completed_at IS NULL OR last_completed_at < CURRENT_DATE)
			RETURNING id
		`
		err = tx.QueryRow(updateProtocolQuery, payload.ProtocolID, userID).Scan(&updatedProtocolID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Protocol not found or already completed today", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to update protocol", http.StatusInternalServerError)
			return
		}

		// 3. Injeta 50 XP e calcula o Level Up diretamente no SQL
		var currentXP, currentLevel int
		updateUserQuery := `
			UPDATE users 
			SET current_xp = current_xp + 50,
			    current_level = ((current_xp + 50) / 1000) + 1
			WHERE id = $1 
			RETURNING current_xp, current_level
		`
		err = tx.QueryRow(updateUserQuery, userID).Scan(&currentXP, &currentLevel)
		if err != nil {
			http.Error(w, "Failed to update user XP", http.StatusInternalServerError)
			return
		}

		// 4. Grava no Extrato (Ledger)
		insertLedgerQuery := `
			INSERT INTO ledger_entries (user_id, amount, entry_type, description) 
			VALUES ($1, $2, 'GAINED', $3)
		`
		_, err = tx.Exec(insertLedgerQuery, userID, 50, "Protocol Completed")
		if err != nil {
			http.Error(w, "Failed to write to ledger", http.StatusInternalServerError)
			return
		}

		// 5. Confirma a Transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Protocol completed successfully!",
			"xp_gained":     50,
			"new_total_xp":  currentXP,
			"current_level": currentLevel,
		})
	}
}
