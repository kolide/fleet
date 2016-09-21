package kolide

// AppConfigStore contains method for saving and retrieving
// application configuration
type AppConfigStore interface {
	NewOrgInfo(info *OrgInfo) (*OrgInfo, error)
	OrgInfo() (*OrgInfo, error)
	SaveOrgInfo(info *OrgInfo) error
}

// OrgInfo holds information about the current
// organization using Kolide
type OrgInfo struct {
	ID         uint `gorm:"primary_key"`
	OrgName    string
	OrgLogoURL string
}
