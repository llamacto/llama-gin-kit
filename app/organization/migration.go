package organization

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// GetMigrations returns the organization module migrations
func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202506181_create_organizations",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Organization{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("organizations")
			},
		},
		{
			ID: "202506182_create_teams",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Team{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("teams")
			},
		},
		{
			ID: "202506183_create_roles",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Role{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("organization_roles")
			},
		},
		{
			ID: "202506184_create_permissions",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Permission{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("organization_permissions")
			},
		},
		{
			ID: "202506185_create_members",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Member{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("organization_members")
			},
		},
		{
			ID: "202506186_create_invitations",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(&Invitation{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("organization_invitations")
			},
		},
		{
			ID: "202506187_create_default_roles",
			Migrate: func(db *gorm.DB) error {
				// Create default system roles (not associated with any organization)
				adminRole := &Role{
					Name:        "admin",
					DisplayName: "Administrator",
					Description: "Full system administrator with all permissions",
					Permissions: `{"*":"*"}`, // Wildcard for all permissions
					IsDefault:   false,
				}
				
				memberRole := &Role{
					Name:        "member",
					DisplayName: "Member",
					Description: "Regular member with limited permissions",
					Permissions: `{
						"organization.view": true,
						"team.view": true,
						"member.view": true
					}`,
					IsDefault: true,
				}
				
				// Create default organization roles
				ownerRole := &Role{
					Name:        "owner",
					DisplayName: "Owner",
					Description: "Organization owner with all organization permissions",
					Permissions: `{
						"organization.*": true,
						"team.*": true,
						"member.*": true,
						"role.*": true,
						"invitation.*": true
					}`,
					IsDefault: false,
				}
				
				managerRole := &Role{
					Name:        "manager",
					DisplayName: "Manager",
					Description: "Organization manager with management permissions",
					Permissions: `{
						"organization.view": true,
						"team.*": true,
						"member.*": true,
						"invitation.*": true
					}`,
					IsDefault: false,
				}
				
				// Add default permissions
				result := db.Create([]*Role{adminRole, memberRole, ownerRole, managerRole})
				return result.Error
			},
			Rollback: func(db *gorm.DB) error {
				return db.Where("name IN ?", []string{"admin", "member", "owner", "manager"}).Delete(&Role{}).Error
			},
		},
	}
}
