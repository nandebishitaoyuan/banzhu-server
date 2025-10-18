package convert

import "httpServerTest/internal/model"

type ReadingRecordConvert struct{}

func (ReadingRecordConvert) EntityToVo(entity model.ReadingRecord) model.ReadingRecordDataVo {
	return model.ReadingRecordDataVo{
		Index:  entity.Index,
		Offset: entity.Offset,
	}
}
