package main 

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Criamos uma "chave" única para guardar o ID do usuário no Contexto da requisição
type contextKey string 
const userContextKey = contextKey("userID")

// AuthMiddleware é o nosso "Segurança". Ele abraça outras rotas.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Pega o "Crachá" no cabeçalho da requisição
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Access Denied: Missing Token", http.StatusUnauthorized)
			return
		}
		// 2. O padrão da web é enviar "Bearer <token>". Nós cortamos a palavra "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Lemos o token usando a nossa senha secreta
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Se o token for falso ou estiver vencido, bloqueia!
		if err != nil || !token.Valid {
			http.Error(w, "Access Denied: Invalid Token", http.StatusUnauthorized)
			return
		}

		// 4. A Mágica do Go: Guardamos o ID do usuário no "bolso" da requisição (Context)
		ctx := context.WithValue(r.Context(), userContextKey, claims.UserID)

		// 5. Deixa o usuário passar para a rota que ele queria, agora com o ID no bolso
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
