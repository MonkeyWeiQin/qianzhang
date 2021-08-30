package response

import (
	"battle_rabbit/define"
	"battle_rabbit/model"
)

type GetEquipmentListResponse struct {
	List          []*model.EquipmentModel `json:"list"`
	EquipmentType define.EquipageType  `json:"equipmentType"`
}
