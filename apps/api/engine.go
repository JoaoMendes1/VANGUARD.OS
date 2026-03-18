package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type RewardPayload struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func AdminRewardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.Context().Value(userContextKey).(string)

		var payload RewardPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		// Inicia a Transação Segura
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // Se der erro no meio, desfaz tudo

		// Dá o XP
		var currentXP int
		err = tx.QueryRow(`UPDATE users SET current_xp = current_xp + $1 WHERE id = $2 RETURNING current_xp`, payload.Amount, userID).Scan(&currentXP)
		if err != nil {
			http.Error(w, "Failed to update user XP", http.StatusInternalServerError)
			return
		}

		// Calcula o Level (A cada 1000 XP = 1 Nível)
		newLevel := (currentXP / 1000) + 1
		_, err = tx.Exec(`UPDATE users SET current_level = $1 WHERE id = $2`, newLevel, userID)
		if err != nil {
			http.Error(w, "Failed to update level", http.StatusInternalServerError)
			return
		}

		// Salva no Extrato
		_, err = tx.Exec(`INSERT INTO ledger_entries (user_id, amount, entry_type, description) VALUES ($1, $2, 'GAINED', $3)`, userID, payload.Amount, payload.Description)
		if err != nil {
			http.Error(w, "Failed to write to ledger", http.StatusInternalServerError)
			return
		}

		// Confirma a transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "XP Injected Successfully!",
			"xp_gained":     payload.Amount,
			"new_total_xp":  currentXP,
			"current_level": newLevel,
		})
	}
}