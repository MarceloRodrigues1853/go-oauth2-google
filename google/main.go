// Login com Google usando Goth.
//
// Documentação:
//   - Goth: https://github.com/markbates/goth
//   - OAuth 2.0 do Google: https://developers.google.com/identity/protocols/oauth2/web-server
//   - OpenID Connect: https://developers.google.com/identity/openid-connect/openid-connect
package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env para o sistema
	// Isso evita que credenciais sensíveis fiquem hardcoded (chumbadas) no código fonte.
	godotenv.Load()

	// Configura o provedor de identidade (Google) com as credenciais obtidas no Google Cloud Console.
	// Os escopos "openid", "email" e "profile" definem quais dados do usuário estamos pedindo permissão para acessar.
	provider := google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK_URL"),
		"openid",
		"email",
		"profile",
	)

	// "online" significa que não estamos pedindo um refresh token para acesso offline prolongado.
	provider.SetAccessType("online")

	// Registra o provedor no pacote Goth para que ele gerencie o fluxo.
	goth.UseProviders(provider)

	// Geração de uma chave aleatória de 32 bytes para assinar os cookies de sessão.
	// O uso de crypto/rand garante que a chave seja criptograficamente segura e imprevisível.
	sessionKey := make([]byte, 32)
	if _, err := rand.Read(sessionKey); err != nil {
		log.Fatal(err)
	}

	// Inicializa o gerenciador de sessões do Gothic usando a nossa chave segura.
	gothic.Store = newSessionStore(sessionKey)

	log.Println("acesse http://localhost:3000")
	// Inicia o servidor web na porta 3000, utilizando as rotas definidas abaixo.
	log.Fatal(http.ListenAndServe(":3000", routes()))
}

// newSessionStore configura como os cookies de sessão serão armazenados no navegador do usuário.
func newSessionStore(sessionKey []byte) *sessions.CookieStore {
	store := sessions.NewCookieStore(sessionKey)

	// Configurações de segurança e comportamento do cookie
	store.Options = &sessions.Options{
		Path:     "/",                  // O cookie é válido para todo o site
		MaxAge:   10 * 60,              // Tempo de vida da sessão: 10 minutos (em segundos)
		HttpOnly: true,                 // Proteção contra ataques XSS: impede que o JavaScript do lado do cliente leia o cookie
		Secure:   false,                // Em ambiente de produção (com HTTPS), isso DEVE ser true. False permite uso em localhost (HTTP).
		SameSite: http.SameSiteLaxMode, // Proteção contra ataques CSRF: controla envio de cookies em requisições cross-site
	}
	return store
}

// routes centraliza o mapeamento das URLs para as suas respectivas funções de tratamento (handlers).
func routes() http.Handler {
	mux := http.NewServeMux()

	// Rota principal que serve a página HTML estática
	mux.HandleFunc("GET /", home)

	// Rota que inicia o fluxo de autenticação. O {provider} será substituído por "google" na URL.
	mux.HandleFunc("GET /auth/{provider}", gothic.BeginAuthHandler)

	// Rota de callback onde o Google redireciona o usuário de volta com o código de autorização.
	mux.HandleFunc("GET /auth/{provider}/callback", authCallback)

	return mux
}

// home apenas entrega o arquivo index.html para o navegador.
func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")
}

// authCallback processa a resposta final do Google, extraindo os dados do usuário.
func authCallback(w http.ResponseWriter, r *http.Request) {
	// CompleteUserAuth faz o "trabalho sujo" de trocar o código de autorização pelo token de acesso
	// e depois buscar os dados do perfil do usuário no Google.
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "não foi possível concluir o login com Google", http.StatusUnauthorized)
		return
	}

	// Loga os dados no terminal do servidor (útil para debug)
	fmt.Println("Usuário retornado pelo Google")
	fmt.Println("Provider:", user.Provider)
	fmt.Println("ID:", user.UserID)
	fmt.Println("Nome:", user.Name)
	fmt.Println("E-mail:", user.Email)

	// Responde ao navegador com os dados em texto puro
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Login concluído!\n\nProvider: %s\nID: %s\nNome: %s\nE-mail: %s\n",
		user.Provider, user.UserID, user.Name, user.Email)
}
