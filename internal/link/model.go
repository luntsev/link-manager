package link

import (
	"link-manager/internal/stat"
	"link-manager/pkg/token"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Name  string      `json:"name"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url, name string) *Link {
	link := Link{
		Url:  url,
		Name: name,
	}
	link.GenHash(3)
	return &link
}

func (l *Link) GenHash(n int) {
	l.Hash = token.GenToken(n)
}
