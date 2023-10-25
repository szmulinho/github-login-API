package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	oauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func LoggedHandler(w http.ResponseWriter, r *http.Request, githubData string) {

	cookie, err := r.Cookie("jwt")
	if err != nil {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("tajny_klucz_do_podpisu_tokena"), nil
	})

	if err != nil || !token.Valid {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "Błąd parsowania danych użytkownika!")
		return
	}

	username, usernameExists := claims["username"].(string)
	email, emailExists := claims["email"].(string)

	if !usernameExists || !emailExists {
		fmt.Fprintf(w, "Brak danych użytkownika w tokenie!")
		return
	}
	fmt.Fprintf(w, "Zalogowany użytkownik: %s, Email: %s", username, email)

	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func (h *handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/github/login/">LOGIN</a>`)
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := oauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)

	http.Redirect(w, r, redirectURL, 301)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal("OAuth exchange failed:", err)
	}

	githubData := getGithubData(token.AccessToken)

	LoggedHandler(w, r, githubData)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": token.AccessToken,
	})

	tokenString, err := jwtToken.SignedString([]byte("token"))
	if err != nil {
		log.Fatal("Błąd podpisywania tokena JWT:", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: tokenString,
	})
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}
