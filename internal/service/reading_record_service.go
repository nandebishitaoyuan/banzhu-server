package service

import (
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"
	"httpServerTest/internal/model/convert"
	"httpServerTest/pkg/snowflake"
)

type ReadingRecordService struct{}

func NewReadingRecordService() *ReadingRecordService {
	return &ReadingRecordService{}
}

func (s ReadingRecordService) Update(param model.ReadingRecordDataParam, userId uint64) error {

	data := model.ReadingRecord{
		BookId:    param.BookId,
		ChapterId: param.ChapterId,
		Index:     param.Index,
		Offset:    param.Offset,
		UserId:    &userId,
	}
	err := database.DB.Where("book_id = ? AND chapter_id = ? AND user_id = ?", param.BookId, param.ChapterId, userId).Delete(&model.ReadingRecord{}).Error
	if err != nil {
		return err
	}
	id := snowflake.GenerateID()
	data.ID = &id
	return database.DB.Create(&data).Error
}

func (s ReadingRecordService) Get(param model.ReadingRecordQueryParam, userId uint64) (*model.ReadingRecordDataVo, error) {
	var dbData model.ReadingRecord
	err := database.DB.Where("book_id = ? AND chapter_id = ? AND user_id = ?", param.BookId, param.ChapterId, userId).Find(&dbData).Error
	if err != nil {
		return nil, err
	}
	vo := convert.ReadingRecordConvert{}.EntityToVo(dbData)
	return &vo, nil
}

func (s ReadingRecordService) GetListByBookId(bookId, userId uint64) (*[]model.ReadingRecordDataVo, error) {
	var dbData []model.ReadingRecord
	err := database.DB.Where("book_id = ? AND user_id = ?", bookId, userId).Find(&dbData).Error
	if err != nil {
		return nil, err
	}
	vo := convert.ReadingRecordConvert{}.EntityListToVoList(dbData)
	return &vo, nil
}
