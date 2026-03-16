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

	// 5. Subindo o servidor na porta 8080
	fmt.Println("⚡ API Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server crashed: ", err)

	}
}

