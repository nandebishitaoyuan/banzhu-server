package convert

import (
	"httpServerTest/internal/model"
)

type BookConvert struct{}

func (BookConvert) EntityToVo(entity model.Book) model.BookDataVo {
	return model.BookDataVo{
		ID:        entity.ID,
		Name:      entity.Name,
		Path:      entity.Path,
		Author:    entity.Author,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (BookConvert) EntityToVoList(entityList []model.Book) []model.BookDataVo {
	result := make([]model.BookDataVo, len(entityList))
	for i, item := range entityList {
		result[i] = BookConvert{}.EntityToVo(item)
	}
	return result
}

func (BookConvert) EntityPageToVoPage(pageResult *model.PageResult[model.Book]) *model.PageResult[model.BookDataVo] {
	voList := BookConvert{}.EntityToVoList(pageResult.Data)
	vo := &model.PageResult[model.BookDataVo]{
		Total:     pageResult.Total,
		Data:      voList,
		PageSize:  pageResult.PageSize,
		PageIndex: pageResult.PageIndex,
	}
	return vo
}

func (BookConvert) VoToEntity(vo model.BookDataVo) model.Book {
	return model.Book{
		ID:     vo.ID,
		Name:   vo.Name,
		Path:   vo.Path,
		Author: vo.Author,
	}
}

func (BookConvert) ParamToEntity(param model.BookDataParam) model.Book {
	return model.Book{
		ID:     param.ID,
		Name:   param.Name,
		Path:   param.Path,
		Author: param.Author,
	}
}

func (BookConvert) EntityToParam(entity model.Book) model.BookDataParam {
	return model.BookDataParam{
		ID:     entity.ID,
		Name:   entity.Name,
		Path:   entity.Path,
		Author: entity.Author,
	}
}
