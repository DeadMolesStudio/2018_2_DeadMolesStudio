package database

import (
	"api/models"
)

func GetUserPositionsDescendingPaginated(p *models.FetchScoreboardPage) (
	[]models.Position, int, error) {
	records := []models.Position{}
	total, err := GetCountOfUsers()
	if err != nil {
		return records, total, err
	}

	// TODO: optimize it
	rows, err := db.Queryx(`
		SELECT user_id, nickname, record FROM user_profile
		ORDER BY record DESC
		LIMIT $1
		OFFSET $2`,
		p.Limit, p.Limit*p.Page)
	if err != nil {

	}
	if err := rows.Err(); err != nil {
		return records, total, err
	}
	r := models.Position{}
	for rows.Next() {
		err := rows.StructScan(&r)
		if err != nil {
			return records, total, err
		}
		records = append(records, r)
	}

	return records, total, nil
}
