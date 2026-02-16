package seeder

import (
	"github.com/hadi-projects/go-react-starter/internal/entity"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"gorm.io/gorm"
)

type RolePermission struct {
	Role       string
	Permission string
}

var defaultRoles = []RolePermission{
	{Role: "admin", Permission: "create-user"},
	{Role: "admin", Permission: "delete-user"},
	{Role: "admin", Permission: "edit-user"},
	{Role: "admin", Permission: "get-user"},
	{Role: "admin", Permission: "create-role"},
	{Role: "admin", Permission: "delete-role"},
	{Role: "admin", Permission: "edit-role"},
	{Role: "admin", Permission: "get-role"},
	{Role: "admin", Permission: "create-permission"},
	{Role: "admin", Permission: "delete-permission"},
	{Role: "admin", Permission: "edit-permission"},
	{Role: "admin", Permission: "get-permission"},
	{Role: "admin", Permission: "manage-cache"},
	{Role: "auditor", Permission: "get-audit-log"},
	{Role: "auditor", Permission: "get-auth-log"},
	{Role: "user", Permission: "get-profile"}, // Basic user role
}

func SeedRole(db *gorm.DB) {
	for _, rp := range defaultRoles {
		// 1. Create Permission if not exists
		var perm entity.Permission
		if err := db.FirstOrCreate(&perm, entity.Permission{Name: rp.Permission}).Error; err != nil {
			logger.SystemLogger.Error().Err(err).Msgf("Failed to seed permission: %s", rp.Permission)
			continue
		}

		// 2. Create Role if not exists
		var role entity.Role
		if err := db.FirstOrCreate(&role, entity.Role{Name: rp.Role}).Error; err != nil {
			logger.SystemLogger.Error().Err(err).Msgf("Failed to seed role: %s", rp.Role)
			continue
		}

		// 3. Assign Permission to Role (idempotent)
		var count int64
		// Check association table 'role_has_permissions'
		if err := db.Table("role_has_permissions").
			Where("role_id = ? AND permission_id = ?", role.ID, perm.ID).
			Count(&count).Error; err != nil {
			logger.SystemLogger.Error().Err(err).Msg("Failed to check role permission association")
			continue
		}

		if count == 0 {
			if err := db.Model(&role).Association("Permissions").Append(&perm); err != nil {
				logger.SystemLogger.Error().Err(err).Msgf("Failed to assign permission %s to role %s", rp.Permission, rp.Role)
			}
		}
	}
	logger.SystemLogger.Info().Msg("Role & Permission Seeding Completed!")
}
