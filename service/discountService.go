package service

import (
	"payso-internal-api/handler"
	"payso-internal-api/model"
	"payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

type DiscountService interface {
	GetDiscountService(mid string, page int, row int) (model.DiscountPagination, error)
	CreateDiscountService(body model.CreateDiscount) (model.UpdateResponse, error)
	UpdateDiscountService(body model.UpdateDiscount) (model.UpdateResponse, error)
	DeleteDiscountService(discountID int) (model.UpdateResponse, error)
}

type discountService struct {
	discountHandler handler.DiscountHandler
}

func NewDiscountService(discountHandler handler.DiscountHandler) DiscountService {
	return &discountService{discountHandler}
}

func (s *discountService) GetDiscountService(mid string, page int, row int) (model.DiscountPagination, error) {
	log.Infof("==-- GetDiscountService --==")

	var err error
	var DiscountList []model.DiscountPlayload

	DiscountList, err = repository.GetDiscountRepository(mid, page, row)
	if err != nil {
		log.Errorf("Error from GetDiscountRepository: %v", err)
		return model.DiscountPagination{}, err
	}

	TotalPages, err := repository.GetTotalDiscountRepository()
	if err != nil {
		log.Errorf("Error from GetTotalDiscountRepository: %v", err)
		return model.DiscountPagination{}, err
	}

	result := model.DiscountPagination{
		TotalPages:   TotalPages,
		DiscountList: DiscountList,
	}

	return result, err
}

func (s *discountService) CreateDiscountService(body model.CreateDiscount) (model.UpdateResponse, error) {
	log.Infof("==-- CreateDiscountService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.CreateDiscountRepository(body)
	if err != nil {
		log.Errorf("Error from CreateDiscountRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}

func (s *discountService) UpdateDiscountService(body model.UpdateDiscount) (model.UpdateResponse, error) {
	log.Infof("==-- UpdateDiscountService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.UpdateDiscountRepository(body)
	if err != nil {
		log.Errorf("Error from UpdateDiscountRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}

func (s *discountService) DeleteDiscountService(discountID int) (model.UpdateResponse, error) {
	log.Infof("==-- DeleteDiscountService --==")

	var err error
	var result model.UpdateResponse

	result, err = repository.DeleteDiscountRepository(discountID)
	if err != nil {
		log.Errorf("Error from DeleteDiscountRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return result, err
}
