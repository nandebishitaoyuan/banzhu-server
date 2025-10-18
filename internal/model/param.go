package model

type PageParam struct {
	PageIndex int    `form:"pageIndex" binding:"required,min=1"`
	PageSize  int    `form:"pageSize" binding:"required,min=1"`
	Order     string `form:"order"`
}

type PageResult[T any] struct {
	Total     int64 `json:"total"`     // 总条数
	PageIndex int   `json:"pageIndex"` // 当前页
	PageSize  int   `json:"pageSize"`  // 每页数量
	Data      []T   `json:"data"`      // 数据列表
}

type KeywordParam struct {
	Keyword string `json:"keyword"`
}

type KeywordPageParam[T any] struct {
	PageParam
	Keyword T `json:"keyword"`
}
