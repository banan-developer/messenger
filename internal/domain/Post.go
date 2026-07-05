package domain

// для получения данных
type WallPost struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Img       string `json:"img"`
	CreatedAt string `json:"created_at"`
}

// для отправки данных
type CreateWallRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Img   string `json:"img"`
}
