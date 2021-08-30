package response

type SkinList struct {
	SkinId string `json:"skinId"` // 当前皮肤的ID
	Use    bool   `json:"use"`    // 是否已装备
}
