package auth

type AccessData struct {
	Ok          bool   `json:"ok"`
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
	AppId       string `json:"app_id"`
	Team        Team   `json:"team"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
