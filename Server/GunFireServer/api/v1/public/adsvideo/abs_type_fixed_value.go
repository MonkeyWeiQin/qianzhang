package adsvideo

import "com.xv.admin.server/model/db"

type FixedValue struct{}

func (*FixedValue) handle(AdsConfig *db.AdvertisingVideoConfigModel) {
	panic("implement me")
}

