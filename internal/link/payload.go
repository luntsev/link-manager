package link

type CreateRequest struct {
	Url  string `json:"url" validate:"required,url"`
	Name string `json:"name"`
}

type UpdateRequest struct {
	Url  string `json:"url" validate:"required,url"`
	Name string `json:"name"`
}
