package auth

import "bingobot/pkg/db"

type AuthRepository struct {
	db *db.Database
}

func NewAuthRepository(db *db.Database) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) SaveAccessToken(data AccessData) (int, error) {
	return r.db.Insert(
		"INSERT INTO Tokens (teamId, teamName, accessToken) VALUES ($1, $2, $3)",
		data.Team.ID,
		data.Team.Name,
		data.AccessToken,
	)
}

func (r *AuthRepository) GetAllTokens() ([]AccessData, error) {
	var tokens []AccessData
	rows, err := r.db.Query("SELECT accessToken, teamId, teamName FROM Tokens")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var accessToken string
		var teamId string
		var teamName string
		rows.Scan(&accessToken, &teamId, &teamName)
		tokens = append(tokens, AccessData{
			AccessToken: accessToken,
			Team: Team{
				ID:   teamId,
				Name: teamName,
			},
		})
	}

	return tokens, err
}
