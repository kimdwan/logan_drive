package pilemodel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteFriend struct {
	gorm.Model
	Friend_id uuid.UUID `gorm:"type:uuid;not null;"`
	Friend_1  uuid.UUID `gorm:"type:uuid;not null;"`
	Friend_2  uuid.UUID `gorm:"type:uuid;not null;"`
}

func (DeleteFriend) TableName() string {
	return "DELETE_LOGAN_FRIEND_TB"
}
