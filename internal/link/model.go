package link

import (
	"gorm.io/gorm"
	"link-manager/pkg/token"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	link := Link{
		Url: url,
	}
	link.GenHash()
	return &link
}

func (l *Link) GenHash() {
	l.Hash = token.GenToken(10)
}
