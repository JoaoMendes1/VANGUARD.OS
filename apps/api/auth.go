package main 

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Chave secreta do nosso servidor (Nunca deixe isso exposto num projeto real em produção!)
var jwtKey = []byte("vanguard_super_secret_key_2026")

// Estrutura que o Frontend vai enviar
type Credentials struct {
		Alias    string `json:"alias"`
		Email    string `json:"email"` //Usado apenas no Cadastro 
		Password string `json:"password"`
}

// Estrutura do "Crachá Digital" (JWT)
type Claims struct {
	UserID string `json:"user_id"`
	Alias  string `json:"alias"`
	jwt.RegisteredClaims
}

// 1. ROTA DE CADASTRO 
func RegisterHandler(db *sql.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds Credentials 
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return 
		}

		// Criptografando a senha (custo 10 é o padrão de mercado)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 10)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return 
		}

		// Salvando no Banco de Dados 
		query := `INSERT INTO users (alias, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
		var newUserID string 
		err = db.QueryRow(query, creds.Alias, creds.Email, string(hashedPassword)).Scan(&newUserID)
		if err != nil {
			http.Error(w, "User already exists or database error", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Operative registred successfully!", "id": newUserID})

	}

}

// 2. Rota de login 
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		// Buscando o usuário no banco pelo Alias
		var storedUser User 
		query := `SELECT id, alias, password_hash FROM users WHERE alias = $1`
		err := db.QueryRow(query, creds.Alias).Scan(&storedUser.ID, &storedUser.Alias, &storedUser.PasswordHash)
		if err != nil {
			http.Error(w, "Operative not found", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(creds.Password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Gerando o Crachá Digital (Token expira em 24h)
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID: storedUser.ID,
			Alias:  storedUser.Alias,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime), 
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		//Retornando o Token para o usuário 
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}