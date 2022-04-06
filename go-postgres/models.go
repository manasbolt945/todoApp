package models

import(
	"time"
	"github.com/google/uuid"
)

type Todo struct{
	ID uuid.UUID 'gorm:"type:uuid;default:uuid_generate_v4();primaryKey;"'
	Desc string
	Status string
	Created_At  time.Time 'gorm:"autoCreateTime"'
	Updated_At time.Time  'gorm:"autoUpdateTime'
}

