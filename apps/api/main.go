package main 

import (
		"database/sql"
		"fmt"
		"log"
		"net/http"

		_ "github.com/lib/pq" // Importa o driver do PostGres silenciosamente 
)

func main() {
	// 1. Configuração da string de conexão com o banco
	// Usamos a porta 5433 que configuramos no Docker para evitar o conflito do Codespaces!
	connStr := "user=root password=root dbname=vanguard_os port=5433 sslmode=disable host=127.0.0.1"

	// 2. Abrindo a conexão
	db, err := sql.Open("postgres",  connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	// 3. Testando se o banco realmente responde (Ping)
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database! Is Docker Running? Error: ", err)
	}

	fmt.Println("🚀 [VANGUARD.OS] SUCCESSFULLY CONNECTED TO PostgreSQL!")

	// PATCH DA ISSUE #9: Garante que a coluna de streak exista na tabela
	_, err = db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS current_streak INT DEFAULT 0;`)
	if err != nil {
		log.Printf("Aviso: Falha ao verificar/criar coluna current_streak: %v\n", err)
	}

	// PATCH DA ISSUE #9 (Parte 2): Ensina o banco de dados a aceitar 'PENALTY'
	_, err = db.Exec(`
		ALTER TABLE ledger_entries DROP CONSTRAINT IF EXISTS ledger_entries_entry_type_check;
		ALTER TABLE ledger_entries ADD CONSTRAINT ledger_entries_entry_type_check CHECK (entry_type IN ('GAINED', 'SPENT', 'PENALTY'));
	`)
	if err != nil {
		log.Printf("Aviso: Falha ao atualizar a trava de segurança do ledger: %v\n", err)
	}

	// Inicializa o serviço de rotinas em background
	cronService := StartCronJobs(db)
	defer cronService.Stop()

	RunMigrations(db)

	// 4. Criando uma rota de teste (Health Check)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("VANGUARD.OS Engine is online and ready."))
	})

	http.HandleFunc("/auth/register", RegisterHandler(db))
	http.HandleFunc("/auth/login", LoginHandler(db))

	// Rotas Protegidas (Envelopadas pelo AuthMiddleware)
	http.HandleFunc("/protocols", AuthMiddleware(func(w http.ResponseWriter, r *http.Request){
		// O Go não tem roteador embutido avançado, então separamos o POST do GET manualmente aqui
		if r.Method == http.MethodPost {
			CreateProtocolHandler(db) (w, r)
		} else if r.Method == http.MethodGet {
			GetProtocolsHandler(db) (w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Rota para concluir o hábito (Injeta XP)
	http.HandleFunc("/protocols/complete", AuthMiddleware(CompleteProtocolHandler(db)))

	http.HandleFunc("/engine/reward", AuthMiddleware(AdminRewardHandler(db)))

	fmt.Println("⚡ API Server running on http://localhost:8080")

	// 5. Subindo o servidor na porta 8080
	fmt.Println("⚡ API Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server crashed: ", err)

	}
}

