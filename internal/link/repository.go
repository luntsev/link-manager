package link

import "link-manager/pkg/db"

type LinkRepository struct {
	DataBase *db.Db
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.DataBase.DB.Create(link)
	return link, result.Error
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.Find(&link, "hash = ?", hash)
	return &link, result.Error
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		DataBase: database,
	}
}
