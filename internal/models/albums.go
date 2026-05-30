package models

type Album struct {
	ID     int     `json:"id" db:"id"`
	Title  string  `json:"title" db:"title"`
	Artist string  `json:"artist" db:"artist"`
	Price  float64 `json:"price" db:"price"`
}

type NewAlbum struct {
	Title  string  `json:"title" binding:"required,min=1,max=255"`
	Artist string  `json:"artist" binding:"required,min=1,max=255"`
	Price  float64 `json:"price" binding:"required,gte=0"`
}

// UpdateAlbum — для обновления (все поля опциональны)
type UpdateAlbum struct {
	Title  *string  `json:"title" binding:"omitempty,min=1,max=255"`
	Artist *string  `json:"artist" binding:"omitempty,min=1,max=255"`
	Price  *float64 `json:"price" binding:"omitempty,gte=0"`
}
