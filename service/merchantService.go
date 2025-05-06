package service

import (
	"payso-internal-api/handler"
	"payso-internal-api/model"
	merchantRepository "payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

type MerchantService interface {
	GetMerchantService(mid string, page int, row int) (model.MerchantPagination, error)
	CreateMerchantService(body model.CreateMerchantPayload, ipAddress string) (model.UpdateResponse, error)
	DeleteMerchantService(ReqMasterMerchantID string, ReqMerchantID string) (model.UpdateResponse, error)
}

type merchantService struct {
	merchantHandler handler.MerchantHandler
}

func NewMerchantService(merchantHandler handler.MerchantHandler) MerchantService {
	return &merchantService{merchantHandler}
}

func (s *merchantService) GetMerchantService(mid string, page int, row int) (model.MerchantPagination, error) {
	log.Infof("==-- GetMerchantService --==")

	var err error
	var MerchantList []model.MerchantPlayload

	MerchantList, err = merchantRepository.GetMerchantRepository(mid, page, row)
	if err != nil {
		log.Errorf("Error from GetMerchantRepository: %v", err)
		return model.MerchantPagination{}, err
	}

	TotalPages, err := merchantRepository.GetTotalMerchantRepository(mid)
	if err != nil {
		log.Errorf("Error from GetMerchantRepository: %v", err)
		return model.MerchantPagination{}, err
	}

	MerchantPagination := model.MerchantPagination{
		TotalPages:   TotalPages,
		MerchantList: MerchantList}

	return MerchantPagination, err
}

func (s *merchantService) CreateMerchantService(body model.CreateMerchantPayload, ipAddress string) (model.UpdateResponse, error) {
	log.Infof("==-- CreateMerchantService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = merchantRepository.CreateMerchantRepository(body)
	if err != nil {
		log.Errorf("Error from CreateMerchantRepository: %v", err)
		return model.UpdateResponse{}, err
	}
	return Result, err
}

func (s *merchantService) DeleteMerchantService(ReqMasterMerchantID string, ReqMerchantID string) (model.UpdateResponse, error) {
	log.Infof("==-- DeleteMerchantService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = merchantRepository.DeleteMerchantRepository(ReqMasterMerchantID, ReqMerchantID)
	if err != nil {
		log.Errorf("Error from DeleteMerchantRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}
