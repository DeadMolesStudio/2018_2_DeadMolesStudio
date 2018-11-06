package database

import (
	"sort"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

func GetUserPositionsDescendingPaginated(p *models.FetchScoreboardPage) (
	[]models.Position, int, error) {
	var records []models.Position
	if p.Page*p.Limit >= uint(len(users)) {
		return []models.Position{}, len(users), nil
	}
	// example: page 0 = 0..9 positions (limit = 10)
	if p.Limit != 0 {
		for i := p.Page * p.Limit; i < uint(len(users)) && i < (p.Page+1)*p.Limit; i++ {
			records = append(records, models.Position{
				ID:       users[i].UserID,
				Nickname: users[i].Nickname,
				Points:   users[i].Record,
			})
		}
	} else {
		for _, v := range users {
			records = append(records, models.Position{
				ID:       v.UserID,
				Nickname: v.Nickname,
				Points:   v.Record,
			})
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Points > records[j].Points
	})

	return records, len(users), nil
}
