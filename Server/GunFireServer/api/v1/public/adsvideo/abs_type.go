package adsvideo

import (
	"com.xv.admin.server/model/db"
)

const (
	ABS_TYPE_RANDOM      int = 1 // 随机
	ABS_TYPE_MULTIPLE    int = 2 // 倍数
	ABS_TYPE_AWARDS      int = 3 // 奖励次数
	ABS_TYPE_FIXED_VALUE int = 4 // 固定值
	ABS_TYPE_MUST_SEE    int = 5 // 必须看
	ABS_TYPE_DAILY_TIMES int = 6 // 每天观看次数
)

type AbsType struct {
	strategy AbsTypeStrategy
}

type AbsTypeStrategy interface {
	handle(AdsConfig *db.AdvertisingVideoConfigModel) float64
}

func language(strategy AbsTypeStrategy) *AbsType {
	return &AbsType{
		strategy: strategy,
	}
}

// Handle
//func Handle(AdsConfig *db.AdvertisingVideoConfigModel, userId int64) (gold float64, times int, err error) {
//	if AdsConfig.Times > 0 {
//		UserAdvertisingVideoDataModel := new(db.UserAdvertisingVideoDataModel)
//		num, GetTotalNumErr := UserAdvertisingVideoDataModel.GetTotalNum(userId, AdsConfig.AvType, time.Now().Format("2006-01-02"))
//		if GetTotalNumErr != nil {
//			return 0, 0, GetTotalNumErr
//		}
//		if int(num) >= AdsConfig.Times {
//			return 0, 0, errors.New("次数上限")
//		}
//	}
//	userModel := new(db.UserModel)
//	_, _ = userModel.GetUserByUserId(userId)
//	switch AdsConfig.AvType {
//		case ABS_TYPE_RANDOM:
//			gold = language(&Random{}).strategy.handle(AdsConfig)
//			break
//		case ABS_TYPE_MULTIPLE:
//			gold = language(&Multiple{}).strategy.handle(AdsConfig)
//			break
//		//case ABS_TYPE_AWARDS:
//		//	language(&Awards{}).strategy.handle(AdsConfig)
//		//	break
//		//case ABS_TYPE_FIXED_VALUE:
//		//	language(&FixedValue{}).strategy.handle(AdsConfig)
//		//	break
//		//case ABS_TYPE_MUST_SEE:
//		//	language(&MustSee{}).strategy.handle(AdsConfig)
//		//	break
//		//case ABS_TYPE_DAILY_TIMES:
//		//	language(&DailyTimes{}).strategy.handle(AdsConfig)
//		//	break
//		default:
//			return 0, 0, errors.New("广告类型错误")
//	}
//
//
//	return gold, 99, nil
//}

func VerificationType(adsType int) bool {
	GetAbsType := GetAbsType()
	for _, value := range GetAbsType {
		if adsType == value {
			return true
		}
	}
	return false
}

func GetAbsType() []int {
	return []int{ABS_TYPE_RANDOM, ABS_TYPE_MULTIPLE, ABS_TYPE_AWARDS, ABS_TYPE_FIXED_VALUE, ABS_TYPE_MUST_SEE, ABS_TYPE_DAILY_TIMES}
}
