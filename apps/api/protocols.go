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
