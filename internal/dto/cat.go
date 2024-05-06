package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

// 	"race": "", /** not null, enum of:
// 			- "Persian"
// 			- "Maine Coon"
// 			- "Siamese"
// 			- "Ragdoll"
// 			- "Bengal"
// 			- "Sphynx"
// 			- "British Shorthair"
// 			- "Abyssinian"
// 			- "Scottish Fold"
// 			- "Birman" */

type (
	ReqAddOrUpdateCat struct {
		Name        string   `json:"name" validate:"required,min=1,max=30"`
		Race        string   `json:"race" validate:"required,oneof=Persian Maine_Coon Siamese Ragdoll Bengal Sphynx British_Shorthair Abyssinian Scottish_Fold Birman"`
		Sex         string   `json:"sex" validate:"required,oneof=male female"`
		AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
		Description string   `json:"description" validate:"required,min=1,max=200"`
		ImageUrls   []string `json:"imageUrls" validate:"required"`
	}
	ParamGetCat struct {
		ID         string `json:"id"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
		Race       string `json:"race"`
		Sex        string `json:"sex"`
		HasMatched bool   `json:"hasMatched"`
		AgeInMonth string `json:"ageInMonth"`
		Owned      bool   `json:"owned"`
		Search     string `json:"search"`
	}

	ResAddCat struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}
	ResGetCat struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Race        string   `json:"race"`
		Sex         string   `json:"sex"`
		AgeInMonth  int      `json:"ageInMonth"`
		Description string   `json:"description"`
		ImageUrls   []string `json:"imageUrls"`
		HasMatched  bool     `json:"hasMatched"`
		CreatedAt   string   `json:"createdAt"`
	}
)

func (d *ReqAddOrUpdateCat) ToCatEntity(userId string) entity.Cat {
	return entity.Cat{
		Name:        d.Name,
		Race:        d.Race,
		Sex:         d.Sex,
		AgeInMonth:  d.AgeInMonth,
		Description: d.Description,
		ImageUrls:   d.ImageUrls,
		HasMatched:  false,
		UserID:      userId,
	}
}
