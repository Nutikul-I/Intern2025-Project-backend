package service

import (
	"payso-internal-api/handler"
	"payso-internal-api/model"
	"payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

type CategoryService interface {
	GetCategoryService(mid string, page int, row int) (model.CategoryPagination, error)
	CreateCategoryService(body model.CreateCategory) (model.UpdateResponse, error)
	UpdateCategoryService(body model.UpdateCategory) (model.UpdateResponse, error)
	DeleteCategoryService(categoryID int) (model.UpdateResponse, error)
}

type categoryService struct {
	categoryHandler handler.CategoryHandler
}

func NewCategoryService(categoryHandler handler.CategoryHandler) CategoryService {
	return &categoryService{categoryHandler}
}

func (s *categoryService) GetCategoryService(mid string, page int, row int) (model.CategoryPagination, error) {
	log.Infof("==-- GetCategoryService --==")

	var err error
	var CategoryList []model.CategoryPlayload

	// ในที่นี้ยังไม่มี pagination ใน SQL_GET_CATEGORY ต้องเพิ่มใน repository ถ้าจะรองรับ
	CategoryList, err = repository.GetCategoryRepository(mid, page, row)
	if err != nil {
		log.Errorf("Error from GetCategoryRepository: %v", err)
		return model.CategoryPagination{}, err
	}

	TotalPages, err := repository.GetTotalCategoryRepository()
	if err != nil {
		log.Errorf("Error from GetTotalCategoryRepository: %v", err)
		return model.CategoryPagination{}, err
	}

	result := model.CategoryPagination{
		TotalPages:   TotalPages,
		CategoryList: CategoryList,
	}

	return result, err
}

func (s *categoryService) CreateCategoryService(body model.CreateCategory) (model.UpdateResponse, error) {
	log.Infof("==-- CreateCategoryService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.CreateCategoryRepository(body)
	if err != nil {
		log.Errorf("Error from CreateCategoryRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}

func (s *categoryService) UpdateCategoryService(body model.UpdateCategory) (model.UpdateResponse, error) {
	log.Infof("==-- UpdateCategoryService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.UpdateCategoryRepository(body)
	if err != nil {
		log.Errorf("Error from UpdateCategoryRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}

func (s *categoryService) DeleteCategoryService(categoryID int) (model.UpdateResponse, error) {
	log.Infof("==-- DeleteCategoryService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.DeleteCategoryRepository(categoryID)
	if err != nil {
		log.Errorf("Error from DeleteCategoryRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}
