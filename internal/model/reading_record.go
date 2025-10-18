package model

type ReadingRecord struct {
	ID        *uint64   `json:"id" gorm:"primaryKey"`
	BookId    *uint64   `json:"bookId" gorm:"index"`
	ChapterId *uint64   `json:"chapterId" gorm:"index"`
	UserId    *uint64   `json:"userId" gorm:"index"`
	Index     *int      `json:"index"`
	Offset    *uint64   `json:"offset"`
	CreatedAt *JSONTime `json:"createdAt"`
	UpdatedAt *JSONTime `json:"updatedAt"`
}

type ReadingRecordDataParam struct {
	BookId    *uint64 `json:"bookId"`
	ChapterId *uint64 `json:"chapterId"`
	Index     *int    `json:"index"`
	Offset    *uint64 `json:"offset"`
}

type ReadingRecordQueryParam struct {
	BookId    *uint64 `json:"bookId"`
	ChapterId *uint64 `json:"chapterId"`
}

type ReadingRecordDataVo struct {
	Index  *int    `json:"index"`
	Offset *uint64 `json:"offset"`
}
