package adsvideo

import (
	"com.xv.admin.server/model/db"
	"math/rand"
)

type Random struct{}

func (*Random) handle(AdsConfig *db.AdvertisingVideoConfigModel) float64 {
	return randFloats(AdsConfig.MinGold,AdsConfig.MinGold)
}

func randFloats(min float64, max float64) float64 {
	return min + rand.Float64() * (max - min)
}