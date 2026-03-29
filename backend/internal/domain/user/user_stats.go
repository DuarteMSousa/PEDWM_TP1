package user

type UserStats struct {
	UserId string `json:"user_id"`
	Games  int    `json:"games"`
	Wins   int    `json:"wins"`
	Elo    int    `json:"elo"`
}

func NewUserStats(userId string) *UserStats {
	return &UserStats{
		UserId: userId,
		Games:  0,
		Wins:   0,
		Elo:    1000,
	}
}

func (s *UserStats) RecordGame(won bool) {
	s.Games++
	if won {
		s.Wins++
		s.Elo += 10
	} else {
		s.Elo -= 10
	}
}
