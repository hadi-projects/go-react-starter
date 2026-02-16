package service

import (
	"github.com/hadi-projects/go-react-starter/internal/dto"
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/internal/repository"
)

type RoleService interface {
	Create(req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetAll() ([]dto.RoleResponse, error)
	GetByID(id uint) (*dto.RoleResponse, error)
	Update(id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	Delete(id uint) error
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) Create(req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &entity.Role{
		Name: req.Name,
	}

	if err := s.roleRepo.Create(role, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Fetch again to get permissions populated (or we can construct response manually if we trust repo)
	// Better to fetch to be sure.
	createdRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdRole), nil
}

func (s *roleService) GetAll() ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.mapToResponse(&role))
	}
	return responses, nil
}

func (s *roleService) GetByID(id uint) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(role), nil
}

func (s *roleService) Update(id uint, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	if err := s.roleRepo.Update(role, req.PermissionIDs); err != nil {
		return nil, err
	}

	updatedRole, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(updatedRole), nil
}

func (s *roleService) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *roleService) mapToResponse(role *entity.Role) *dto.RoleResponse {
	var permissions []dto.PermissionResponse
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
