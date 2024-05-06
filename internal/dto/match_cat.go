package dto

import (
	"time"

	"github.com/sprint-id/eniqilo-server/internal/entity"
)

// ReqMatchCat is a struct to represent request payload for match cat
// {
// 	"matchCatId": "",
// 	"userCatId": "",
// 	"message": "" // not null, minLength: 5, maxLength: 120
// }

// ResGetMatchCat is a struct to represent response payload for get match cat
// {
// 	"message": "success",
// 	"data": [ // ordered by newest first
// 		{
// 			"id": "", // use whatever id
// 			"issuedBy": {
// 				"name": "",
// 				"email": "",
// 				"createdAt": "" // should in ISO 8601 format
// 			},
// 			"matchCatDetail": {
// 				"id": "",
// 				"name": "",
// 				"race": "",
// 				"sex": "",
// 				"description":"",
// 				"ageInMonth": 1,
// 				"imageUrls":[
// 					"","",""
// 				],
// 				"hasMatched": false,
// 				"createdAt": "" // should in ISO 8601 format
// 			},
// 			"userCatDetail": {
// 				"id": "",
// 				"name": "",
// 				"race": "",
// 				"sex": "",
// 				"description":"",
// 				"ageInMonth": 1,
// 				"imageUrls": [
// 					"","",""
// 				],
// 				"hasMatched": false,
// 				"createdAt": "" // should in ISO 8601 format
// 			},
// 			"message": "",
// 			"createdAt": "" // should in ISO 8601 format
// 		}
// 	]
// }

// ReqApproveOrRejectMatchCat is a struct to represent request payload for approve or reject match cat
// {
// 	"matchId":""
// }

type (
	ReqMatchCat struct {
		MatchCatId string `json:"matchCatId" validate:"required"`
		UserCatId  string `json:"userCatId" validate:"required"`
		Message    string `json:"message" validate:"required,min=5,max=120"`
	}

	ResGetMatchCat struct {
		ID             string          `json:"id"`
		IssuedBy       entity.IssuedBy `json:"issuedBy"`
		MatchCatDetail entity.Cat      `json:"matchCatDetail"`
		UserCatDetail  entity.Cat      `json:"userCatDetail"`
		Message        string          `json:"message"`
		CreatedAt      string          `json:"createdAt"`
	}

	ReqApproveOrRejectMatchCat struct {
		MatchId string `json:"matchId" validate:"required"`
	}

	ResMatchCat struct {
		MatchId string `json:"matchId"`
	}
)

// ToMatchCatEntity is a function to convert ReqMatchCat to MatchCat entity
func (d *ReqMatchCat) ToMatchCatEntity(name, email string) entity.MatchCat {
	return entity.MatchCat{
		IssuedBy: entity.IssuedBy{
			Name:      name,
			Email:     email,
			CreatedAt: time.Now().Format(time.RFC3339), // should in ISO 8601 format for time : now
		},
		MatchCatId: d.MatchCatId,
		UserCatId:  d.UserCatId,
		Message:    d.Message,
	}
}
