package convert

import "httpServerTest/internal/model"

type ReadingRecordConvert struct{}

func (ReadingRecordConvert) EntityToVo(entity model.ReadingRecord) model.ReadingRecordDataVo {
	return model.ReadingRecordDataVo{
		Index:  entity.Index,
		Offset: entity.Offset,
	}
}

func (ReadingRecordConvert) EntityListToVoList(entityList []model.ReadingRecord) []model.ReadingRecordDataVo {
	var voList []model.ReadingRecordDataVo
	for _, entity := range entityList {
		vo := ReadingRecordConvert{}.EntityToVo(entity)
		voList = append(voList, vo)
	}
	return voList
}
