package stat

import (
	"link-manager/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	DataBase *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{DataBase: db}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.DataBase.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 {
		repo.DataBase.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks++
		repo.DataBase.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	repo.DataBase.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}
