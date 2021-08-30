package adsvideo

import "com.xv.admin.server/model/db"

type Multiple struct{}

func (m *Multiple) handle(AdsConfig *db.AdvertisingVideoConfigModel) float64 {
	return 999 * AdsConfig.Multiple
}


