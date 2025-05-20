package service

import (
	"payso-internal-api/model"
	"payso-internal-api/repository"

	log "github.com/sirupsen/logrus"
)

type PermissionService interface {
	GetPermissionService(id int, page int, row int) (model.PermissionResponse, error)
	CreatePermissionService(body model.CreatePermissionPayload) (model.UpdateResponse, error)
	UpdatePermissionService(body model.UpdatePermissionPayload) (model.UpdateResponse, error)
	DeletePermissionService(id int) (model.UpdateResponse, error)
	GetPermissionByIDService(id int) (model.Permission, error)
}

type permissionService struct {
	// ไม่ต้องใช้ permissionHandler ตรงนี้
}

func NewPermissionService() PermissionService {
	return &permissionService{}
}

func (s *permissionService) GetPermissionService(id int, page int, row int) (model.PermissionResponse, error) {
	log.Infof("==-- GetPermissionService --==")

	var err error
	var PermissionList []model.Permission

	PermissionList, err = repository.GetPermissionsRepository(id, page, row)
	if err != nil {
		log.Errorf("Error from GetPermissionsRepository: %v", err)
		return model.PermissionResponse{}, err
	}

	TotalPages, err := repository.GetTotalPermissionsRepository(id)
	if err != nil {
		log.Errorf("Error from GetTotalPermissionsRepository: %v", err)
		return model.PermissionResponse{}, err
	}

	PermissionResponse := model.PermissionResponse{
		TotalPages:     TotalPages,
		PermissionList: PermissionList,
	}

	return PermissionResponse, err
}

func (s *permissionService) CreatePermissionService(body model.CreatePermissionPayload) (model.UpdateResponse, error) {
	log.Infof("==-- CreatePermissionService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = repository.CreatePermissionRepository(body)
	if err != nil {
		log.Errorf("Error from CreatePermissionRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *permissionService) UpdatePermissionService(body model.UpdatePermissionPayload) (model.UpdateResponse, error) {
	log.Infof("==-- UpdatePermissionService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = repository.UpdatePermissionRepository(body)
	if err != nil {
		log.Errorf("Error from UpdatePermissionRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *permissionService) DeletePermissionService(id int) (model.UpdateResponse, error) {
	log.Infof("==-- DeletePermissionService --==")

	var err error
	var Result model.UpdateResponse

	Result, err = repository.DeletePermissionRepository(id)
	if err != nil {
		log.Errorf("Error from DeletePermissionRepository: %v", err)
		return model.UpdateResponse{}, err
	}

	return Result, err
}

func (s *permissionService) GetPermissionByIDService(id int) (model.Permission, error) {
	log.Infof("==-- GetPermissionByIDService --==")

	var err error
	var Permission model.Permission

	Permission, err = repository.GetPermissionByIDRepository(id)
	if err != nil {
		log.Errorf("Error from GetPermissionByIDRepository: %v", err)
		return model.Permission{}, err
	}

	return Permission, err
}
