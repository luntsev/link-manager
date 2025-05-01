package link

import (
	"link-manager/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	DataBase *db.Db
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Create(link)
	return link, result.Error
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "hash = ?", hash)
	return &link, result.Error
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(link)
	return link, result.Error
}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.DataBase.DB.Delete(&Link{}, id)
	return result.Error
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "ID = ?", id)
	return &link, result.Error
}

func (repo *LinkRepository) GetCount() int64 {
	var count int64
	repo.DataBase.
		Table("links").
		Where("deleted_at is null").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(page, pageSize int) ([]Link, int64) {
	var links []Link
	var count int64

	query := repo.DataBase.
		Table("links").
		Where("deleted_at is null").
		Session(&gorm.Session{})

	query.
		Order("id asc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&links)

	query.
		Count(&count)

	return links, count
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		DataBase: database,
	}
}
