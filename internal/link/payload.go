package link

type CreateRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type CreateResponse struct {
	Url  string `json:"url" gorm:"required"`
	Hash string `json:"hash" gorm:"required"`
}

type ReadRequest struct {
	Hash string `json:"hash" gorm:"required"`
}

type ReadResponse struct {
	Url  string `json:"url" gorm:"required"`
	Hash string `json:"hash" gorm:"required"`
}

type UpdateRequest struct {
	Hash string `json:"hash" gorm:"required"`
}

type UpdateResponse struct {
	Url  string `json:"url" gorm:"required"`
	Hash string `json:"hash" gorm:"required"`
}
