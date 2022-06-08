package auth

import (
	"bingobot/pkg/bot"
	"bingobot/pkg/db"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

type authServer struct {
	repository *AuthRepository
	db         *db.Database
}

func StartAuthServer(db *db.Database) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3050"
	}

	server := authServer{
		repository: NewAuthRepository(db),
		db:         db,
	}

	http.HandleFunc("/auth", server.handleAuth)
	http.HandleFunc("/", server.handleHomePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("www"))))

	fmt.Println("Starting server on " + host + ":" + port)
	http.ListenAndServe(host+":"+port, nil)
}

func (s *authServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	accessData, err := getAccessToken(code)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		somethingWentWrong(w)
		return
	}

	_, err = s.repository.SaveAccessToken(accessData)
	if err != nil {
		fmt.Println("Error saving access token:", err)
		somethingWentWrong(w)
		return
	}

	go bot.NewBot(s.db, accessData.AccessToken, accessData.Team.ID).Run()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("App is successfully installed!"))
}

func getAccessToken(code string) (AccessData, error) {
	accessData := AccessData{}

	reqBody := fmt.Sprintf(
		"client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		os.Getenv("SLACK_REDIRECT_URL"),
	)

	res, err := http.Post("https://slack.com/api/oauth.v2.access", "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {
		fmt.Println("failed to fetch access token", err)
		return accessData, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&accessData)
	if !accessData.Ok {
		return accessData, fmt.Errorf(accessData.Error)
	}

	return accessData, nil
}

func (s *authServer) handleHomePage(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("README.md")
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		somethingWentWrong(w)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Error parsing index.html:", err)
		somethingWentWrong(w)
		return
	}

	content = markdown.NormalizeNewlines(content)
	output := markdown.ToHTML(content, nil, nil)

	tmpl.Execute(w, template.HTML(output))
}

func somethingWentWrong(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Something went wrong"))
}
