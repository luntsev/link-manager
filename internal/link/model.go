package link

import (
	"gorm.io/gorm"
	"link-manager/pkg/token"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
	Name string `json:"name"`
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
