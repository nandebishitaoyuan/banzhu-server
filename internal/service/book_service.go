package service

import (
	"fmt"
	"httpServerTest/internal/config"
	"httpServerTest/internal/database"
	"httpServerTest/internal/model"
	"httpServerTest/internal/model/convert"
	"httpServerTest/pkg/file"
	"httpServerTest/pkg/list"
	"httpServerTest/pkg/snowflake"
	"path/filepath"
	"strings"
)

type BookService struct {
	cfg *config.Config
}

func NewBookService(cfg *config.Config) *BookService {
	return &BookService{
		cfg: cfg,
	}
}

func (s *BookService) GetById(bookId uint64) *model.Book {
	var book model.Book
	if database.DB.Where("id = ?", bookId).First(&book).Error != nil {
		return nil
	}
	return &book
}

func (s *BookService) GetPage(param model.KeywordPageParam[string]) (*model.PageResult[model.BookDataVo], error) {

	condition := database.DB.Debug()

	if param.Keyword != "" {
		condition = database.ApplyKeywordSearch[model.Book](condition, param.Keyword)
	}

	pageResult, err := database.Paginate[model.Book](condition, &param.PageParam)
	if err != nil {
		return nil, err
	}
	vo := convert.BookConvert{}.EntityPageToVoPage(pageResult)

	bookIdList := list.Map(vo.Data, func(t model.BookDataVo) uint64 {
		return *t.ID
	})
	chapterList, err := NewChapterService(s.cfg).GetChapterByIdList(bookIdList)
	if err != nil {
		return nil, err
	}

	chapterListMap := list.GroupBy(chapterList, func(t *model.Chapter) uint64 {
		return *t.BookID
	})

	bookVoList := list.Map(vo.Data, func(t model.BookDataVo) model.BookDataVo {
		chapters := chapterListMap[*t.ID]

		numberOfChapters := len(chapters)

		t.NumberOfChapters = &numberOfChapters

		wordCountList := list.Map(chapters, func(t1 *model.Chapter) *int {
			return t1.WordCount
		})
		totalWordCount := 0
		for _, count := range wordCountList {
			if count != nil {
				totalWordCount += *count
			}
		}
		t.WordCount = &totalWordCount
		return t
	})

	vo.Data = bookVoList

	return vo, err
}

func (s *BookService) DeleteBook(id uint64) error {
	book := s.GetById(id)
	var err error
	if book != nil {
		err = file.DeleteDir(*book.Path)
	}
	// 删除章节
	err = ChapterService{}.DeleteChapterByBookId(id)
	if err != nil {
		return err
	}
	// 删除阅读记录
	err = ReadingRecordService{}.DeleteByBookId(id)
	if err != nil {
		return err
	}
	err = database.DB.Delete(&model.Book{}, id).Error
	return err
}

func (s *BookService) SyncBook() error {
	// 读取文件目录
	list, err := file.ListDir(s.cfg.Path.Book)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 查询数据库已有书籍
	var dbBooks []model.Book
	if err := database.DB.Model(&model.Book{}).Find(&dbBooks).Error; err != nil {
		return fmt.Errorf("查询数据库失败: %w", err)
	}

	// 构建文件系统书籍列表
	var fsBooks []model.Book
	for _, info := range list {
		if info.IsDir() {
			name := info.Name()
			relPath := strings.Replace(s.cfg.Path.Chapter, "{}", name, 1)
			absPath, err := filepath.Abs(relPath)
			if err != nil {
				return fmt.Errorf("获取绝对路径失败: %w", err)
			}
			id := snowflake.GenerateID()
			fsBooks = append(fsBooks, model.Book{
				ID:   &id,
				Name: &name,
				Path: &absPath,
			})
		}
	}

	// 注意对比方向：
	// toAdd: 文件系统有，但数据库没有 → 需要插入数据库
	// toDelete: 数据库有，但文件系统没有 → 需要删除数据库
	toAdd, toDelete := s.CompareListsByName(dbBooks, fsBooks)

	if len(toAdd) == 0 && len(toDelete) == 0 {
		fmt.Println("同步完成：没有需要新增或删除的书籍！")
	}

	// 使用事务执行新增和删除操作
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 插入新增的书籍
	if len(toAdd) > 0 {
		if err := tx.Create(&toAdd).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("新增失败: %w", err)
		}
		fmt.Printf("同步完成：成功新增：%d本\n", len(toAdd))
	}

	// 删除不存在的书籍
	if len(toDelete) > 0 {
		ids := make([]uint64, 0, len(toDelete))
		for _, book := range toDelete {
			if book.ID != nil {
				ids = append(ids, *book.ID)
			}
		}
		if err := tx.Where("id IN ?", ids).Delete(&model.Book{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除失败: %w", err)
		}
		fmt.Printf("同步完成：成功删除：%d本\n", len(toDelete))
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %w", err)
	}

	err = NewChapterService(s.cfg).syncChapter(toAdd)
	if err != nil {
		return err
	}

	return nil
}

// CompareListsByName 比较两个列表中对象的 Name（类型为 *stringUtil）
// 返回：需要新增的项、需要删除的项
func (s *BookService) CompareListsByName(dbList, apiList []model.Book) (toAdd, toDelete []model.Book) {
	dbMap := make(map[string]model.Book)
	apiMap := make(map[string]model.Book)

	// Helper 函数：安全获取 Name 值
	getName := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	// 构建 map：数据库数据
	for _, item := range dbList {
		dbMap[getName(item.Name)] = item
	}

	// 构建 map：接口数据
	for _, item := range apiList {
		apiMap[getName(item.Name)] = item
	}

	// 找出要新增的（在接口中有，但数据库中没有）
	for _, item := range apiList {
		name := getName(item.Name)
		if _, exists := dbMap[name]; !exists {
			toAdd = append(toAdd, item)
		}
	}

	// 找出要删除的（在数据库中有，但接口中没有）
	for _, item := range dbList {
		name := getName(item.Name)
		if _, exists := apiMap[name]; !exists {
			toDelete = append(toDelete, item)
		}
	}

	return
}

func (s *BookService) SyncBookById(id uint64) error {
	err := database.DB.Delete(&model.Book{}, id).Error
	if err != nil {
		return err
	}
	err = ChapterService{}.DeleteChapterByBookId(id)
	if err != nil {
		return err
	}
	// 删除阅读记录
	err = ReadingRecordService{}.DeleteByBookId(id)
	if err != nil {
		return err
	}
	err = s.SyncBook()
	return err
}
