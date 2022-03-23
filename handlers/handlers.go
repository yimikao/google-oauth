package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/yimikao/googleoauth/utils"
)

func NewMux() http.Handler {
	m := http.NewServeMux()

	m.Handle("/", http.FileServer(http.Dir("templates/")))

	m.HandleFunc("/auth/google/login", oauthOnLogin)
	m.HandleFunc("/auth/google/callback", oauthOnCallback)

	return m
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)

	// return cookie value
	return state
}

func oauthOnLogin(w http.ResponseWriter, r *http.Request) {

	// add cookie to response and also return it's value
	// 'll be validated that it matches with the state query parameter on redirect callback
	oauthState := generateStateOauthCookie(w)

	// redirect user to googleauth url
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)

}

func oauthOnCallback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthStateCookie, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthStateCookie.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Println("auth code not supplied")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(code)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var ur = new(utils.UserResponse)
	if err := json.Unmarshal(data, ur); err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	t, _ := template.ParseFiles("templates/success.html")
	t.Execute(w, ur)
	// fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// use auth-code to get token and get user info

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	res, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	return bts, nil

}
