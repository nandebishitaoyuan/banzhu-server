package model

type Book struct {
	ID        *uint64   `json:"id" gorm:"primaryKey;type:bigint"`
	Name      *string   `json:"name" gorm:"index;size:64" like:"true"`
	Path      *string   `json:"path" gorm:"size:255"`
	Author    *string   `json:"author" gorm:"size:100" like:"true"`
	CreatedAt *JSONTime `json:"createdAt"`
	UpdatedAt *JSONTime `json:"updatedAt"`
}

type BookDataParam struct {
	ID     *uint64 `json:"id"`
	Name   *string `json:"name"`
	Path   *string `json:"path"`
	Author *string `json:"author"`
}

type BookDataVo struct {
	ID               *uint64   `json:"id"`
	Name             *string   `json:"name"`
	Path             *string   `json:"path"`
	Author           *string   `json:"author"`
	NumberOfChapters *int      `json:"numberOfChapters"`
	WordCount        *int      `json:"wordCount"`
	UpdatedAt        *JSONTime `json:"updatedAt"`
}
