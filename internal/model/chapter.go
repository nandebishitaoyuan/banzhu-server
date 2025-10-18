package model

type Chapter struct {
	ID        *uint64   `json:"id" gorm:"primaryKey;type:bigint"`
	BookID    *uint64   `json:"bookId" gorm:"index"`
	Name      *string   `json:"name" gorm:"index;size:64" like:"true"`
	Path      *string   `json:"path" gorm:"size:255"`
	WordCount *int      `json:"wordCount"`
	Sort      *int      `json:"sort"`
	CreatedAt *JSONTime `json:"createdAt"`
	UpdatedAt *JSONTime `json:"updatedAt"`
}
