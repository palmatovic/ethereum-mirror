package sync

import "gorm.io/gorm"

type Sync struct {
	database *gorm.DB
}

func NewSync(database *gorm.DB) *Sync {
	return &Sync{database: database}
}

// Sync is cron job that disable user due to expired password
func (s *Sync) Sync() {

}
