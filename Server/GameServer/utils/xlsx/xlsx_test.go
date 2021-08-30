package xlsx

import (
	"battle_rabbit/service/log"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"testing"
)

const (
	RoleDataXlsx = "D:\\project\\go\\BattleRabbit\\Server\\GameServer\\bin\\xlsx\\RoleData.xlsx"
)

var (
	PublicAttribute    = new(Attribute)
	AttributeFieldInfo = getFieldInfos(reflect.TypeOf(PublicAttribute).Elem())
	RoleDataConfig = make(map[string]*RoleData)
)



type RoleData struct {
	Id         string     `xlsx:"id"`
	Name       string     `xlsx:"name"`
	Describe   string     `xlsx:"describe"`
	Relation   string     `xlsx:"relation"`
	Star_level int        `xlsx:"star_level"`
	Level      int        `xlsx:"level"`
	Exp        int        `xlsx:"exp"`
	Attribute  Attribute `xlsx:"attribute"`
}

type Attribute struct {
	Life_a          int     `xlsx:"life_a"`
	Life_c          float32 `xlsx:"life_c"`
	Defense_a       int     `xlsx:"defense_a"`
	Defense_c       float32 `xlsx:"defense_c"`
	Attack_a        int     `xlsx:"attack_a"`
	Attack_c        float32 `xlsx:"attack_c"`
	Attack_speed_a  int     `xlsx:"attack_speed_a"`
	Attack_speed_c  float32 `xlsx:"attack_speed_c"`
	Move_speed_a    int     `xlsx:"move_speed_a"`
	Move_speed_c    float32 `xlsx:"move_speed_c"`
	Dodge_b         int     `xlsx:"dodge_b"`
	Dodge_c         int     `xlsx:"dodge_c"`
	Critical_b      int     `xlsx:"critical_b"`
	Critical_c      int     `xlsx:"critical_c"`
	Buff_time_a     float32 `xlsx:"buff_time_a"`
	Buff_time_c     float32 `xlsx:"buff_time_c"`
	Boss_hurt_a     float32 `xlsx:"boss_hurt_a"`
	Boss_hurt_c     float32 `xlsx:"boss_hurt_c"`
	Critical_hurt_a float32 `xlsx:"critical_hurt_a"`
	Critical_hurt_c float32 `xlsx:"critical_hurt_c"`
	Gold_add_c      float32 `xlsx:"gold_add_c"`
}

func TestLoadXlsxConfig(t *testing.T) {

	//var iface  = new(RoleData)
	var iface  = &RoleData{}
	err := LoadXlsxFile(RoleDataXlsx,iface,"Sheet2", func(v interface{}) {
		val := v.(*RoleData)
		RoleDataConfig[val.Id] = val
	})

	if err != nil {
		log.Fatal(err)
	}

	m,err  := jsoniter.Marshal(RoleDataConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(m))

	//xlsx, err := excelize.OpenFile(RoleDataXlsx)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//rows, err := xlsx.GetRows("Sheet2")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//indexMap := make(map[int]string)
	//for i, row := range rows {
	//	if i < 2 {
	//		continue
	//	}
	//
	//	if i == 2 {
	//		for i, key := range row {
	//			indexMap[i] = key
	//		}
	//		continue
	//	}
	//	var iface interface{} = new(RoleData)
	//	outVal, outTy := getConcreteReflectValueAndType(iface)
	//	outTagMap := getFieldInfos(outTy)
	//	var loadSub = false
	//	for i, val := range row {
	//		if fieldIndex, ok := outTagMap[indexMap[i]]; ok { // 第一层
	//			err := xlsx2.setField(outVal.Field(fieldIndex), val, true)
	//			if err != nil {
	//				log.Error(err)
	//				return
	//			}
	//		} else if !loadSub {
	//			for j := 0; j < outVal.NumField(); j++ {
	//				if outVal.Field(j).Kind() == reflect.Struct || outVal.Field(j).Kind() == reflect.Ptr {
	//					out := outVal.Field(j)
	//					if out.Kind() == reflect.Ptr {
	//						out = reflect.New(out.Type())
	//					}
	//					v, err := getSubValue(out.Interface(), indexMap, row)
	//					if err != nil {
	//						log.Error(err)
	//						return
	//					}
	//					loadSub = true
	//					outVal.Field(j).Set(v)
	//					break
	//				}
	//			}
	//		}
	//	}
	//
	//	s, _ := jsoniter.Marshal(iface)
	//
	//	log.Debug(string(s))
	//}

}

//func getSubValue(out interface{}, tagMap map[int]string, orig []string) (reflect.Value, error) {
//	value := reflect.ValueOf(out)
//	if value.Kind() == reflect.Ptr {
//		value = value.Elem()
//	}
//	isP := false
//	innerTy := value.Type()
//	if innerTy.Kind() == reflect.Ptr {
//		isP = true
//		innerTy = innerTy.Elem()
//	}
//	subVal := reflect.New(innerTy)
//	subInner := subVal.Elem()
//
//	subTagMap := getFieldInfos(innerTy)
//	for k, v := range orig {
//		if subFieldIndex, ok := subTagMap[tagMap[k]]; ok {
//			err := xlsx2.setField(subInner.Field(subFieldIndex), v, true)
//			if err != nil {
//				return reflect.Value{}, err
//			}
//		}
//	}
//	if isP {
//		return subVal, nil
//	}
//	return subInner, nil
//}
//
//func getFieldInfos(rType reflect.Type) map[string]int {
//	fieldInfo := make(map[string]int)
//	fieldsCount := rType.NumField()
//	for i := 0; i < fieldsCount; i++ {
//		field := rType.Field(i)
//		if field.PkgPath != "" {
//			continue
//		}
//
//		filteredTags := field.Tag.Get(TagName)
//		if filteredTags == "" || filteredTags == "-" {
//			continue
//		}
//		fieldInfo[filteredTags] = i
//	}
//	return fieldInfo
//}
//
//func getConcreteReflectValueAndType(in interface{}) (reflect.Value, reflect.Type) {
//	value := reflect.ValueOf(in)
//	if value.Kind() == reflect.Ptr {
//		value = value.Elem()
//	}
//	return value, value.Type()
//}
