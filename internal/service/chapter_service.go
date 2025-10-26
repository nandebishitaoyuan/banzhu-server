package service

import (
	"errors"
	"fmt"
	"httpServerTest/internal/config"
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"
	"httpServerTest/pkg/file"
	"httpServerTest/pkg/logger"
	"httpServerTest/pkg/snowflake"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ChapterService struct {
	cfg *config.Config
}

func NewChapterService(cfg *config.Config) *ChapterService {
	return &ChapterService{
		cfg: cfg,
	}
}

func (s ChapterService) GetPage(param model.KeywordPageParam[uint64]) (*model.PageResult[model.Chapter], error) {
	condition := database.DB.Where("book_id = ?", param.Keyword)
	pageResult, err := database.Paginate[model.Chapter](condition, &param.PageParam)
	if err != nil {
		return nil, err
	}
	return pageResult, err
}

func (s ChapterService) GetList(bookId uint64) ([]*model.Chapter, error) {
	var chapter []*model.Chapter
	err := database.DB.Where("book_id = ?", bookId).Order("sort").Find(&chapter).Error
	return chapter, err
}

func (s ChapterService) DeleteChapter(id uint64) error {

	var chapter *model.Chapter
	err := database.DB.Where("id = ?", id).First(&chapter).Error
	if err != nil {
		return err
	}
	err = database.DB.Delete(&model.Chapter{}, id).Error
	if err != nil {
		return err
	}
	err = file.DeleteFile(*chapter.Path)
	if err != nil {
		return err
	}

	return nil
}

func (s ChapterService) DeleteChapterByBookId(bookId uint64) error {

	var chapter []*model.Chapter
	err := database.DB.Where("book_id = ?", bookId).Find(&chapter).Error
	if err != nil {
		return err
	}
	err = database.DB.Where("book_id = ?", bookId).Delete(&model.Chapter{}).Error
	return err
}

func (s ChapterService) syncChapter(books []model.Book) error {
	if books == nil || len(books) == 0 {
		return nil
	}

	ExtractPrefix := func(filename string) (int, error) {
		// 使用正则表达式提取文件名前缀中的数字部分
		re := regexp.MustCompile(`^\d+`)
		matches := re.FindString(filename)
		if matches == "" {
			return 0, fmt.Errorf("文件名中没有数字前缀")
		}

		// 将数字部分转换为整数
		var prefixNum int
		_, err := fmt.Sscanf(matches, "%d", &prefixNum)
		if err != nil {
			return 0, err
		}
		return prefixNum, nil
	}

	var allFileChapter []model.Chapter

	for _, book := range books {
		fileList, err := file.ListDir(*book.Path)
		if err != nil {
			logger.New().Fatalf("读取目录失败: %s", err.Error())
			continue
		}

		// 按照文件名前缀数字进行排序
		sort.Slice(fileList, func(i, j int) bool {
			numI, _ := ExtractPrefix(fileList[i].Name())
			numJ, _ := ExtractPrefix(fileList[j].Name())
			return numI < numJ
		})

		var fileChapter []model.Chapter
		for i, info := range fileList {
			if !info.IsDir() {

				bookPath := *book.Path
				name := strings.ReplaceAll(info.Name(), ".txt", "")
				relPath := bookPath + "/" + info.Name()
				absPath, err := filepath.Abs(relPath)
				if err != nil {
					return fmt.Errorf("获取绝对路径失败: %w", err)
				}
				id := snowflake.GenerateID()
				chapterSort := i + 1

				charCount, err := file.CountText(absPath)

				fileChapter = append(fileChapter, model.Chapter{
					BookID:    book.ID,
					Name:      &name,
					Path:      &absPath,
					ID:        &id,
					WordCount: &charCount,
					Sort:      &chapterSort,
				})
			}
		}
		allFileChapter = append(allFileChapter, fileChapter...)
	}

	// 插入新增的章节
	if len(allFileChapter) > 0 {
		if err := database.DB.Create(&allFileChapter).Error; err != nil {
			return fmt.Errorf("新增失败: %w", err)
		}
		fmt.Printf("同步完成：成功新增：%d章节\n", len(allFileChapter))
	}

	return nil
}

func (s ChapterService) GetChapterById(id64 uint64) (*model.Chapter, error) {
	var chapter *model.Chapter
	err := database.DB.Model(&model.Chapter{}).Where("id = ?", id64).First(&chapter).Error

	if err != nil {
		return nil, errors.New("章节不存在")
	}
	return chapter, nil
}

func (s ChapterService) GetChapterByIdList(id64List []uint64) ([]*model.Chapter, error) {
	var chapter []*model.Chapter
	err := database.DB.Model(&model.Chapter{}).Where("book_id in (?)", id64List).Find(&chapter).Error

	if err != nil {
		return nil, errors.New("章节不存在")
	}
	return chapter, nil
}
