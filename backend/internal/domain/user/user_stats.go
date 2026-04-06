package user

// UserStats armazena as estatísticas de jogo de um utilizador.
type UserStats struct {
	UserId string `json:"user_id"`
	Games  int    `json:"games"`
	Wins   int    `json:"wins"`
	Elo    int    `json:"elo"`
}

// NewUserStats cria estatísticas iniciais com ELO base de 1000.
func NewUserStats(userId string) *UserStats {
	return &UserStats{
		UserId: userId,
		Games:  0,
		Wins:   0,
		Elo:    1000,
	}
}

// RecordGame regista o resultado de um jogo, atualizando contagens e ELO (±10).
func (s *UserStats) RecordGame(won bool) {
	s.Games++
	if won {
		s.Wins++
		s.Elo += 10
	} else {
		s.Elo -= 10
	}
}
