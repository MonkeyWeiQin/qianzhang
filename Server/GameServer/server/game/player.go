package game

import (
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
	"battle_rabbit/model"
	"battle_rabbit/protocol/response"
	"battle_rabbit/service/log"
	"battle_rabbit/utils"
	"battle_rabbit/utils/xid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"runtime/debug"
	"sync"
)

type Player struct {
	Uid        int
	game       *Game
	Account    *model.UserModel
	Role       *model.RoleModel
	Stage      *model.StageModel
	Equip      map[string]*model.EquipmentModel
	Package    map[string]model.IBackPack
	Hero       map[string]*model.HeroModel
	useHero    map[string]*model.HeroModel
	useEquip   map[string]*model.EquipmentModel
	packIndex  map[string]model.IBackPack
	heroIndex  map[string]*model.HeroModel      // 索引,方便查找 key 是配置表中的关联ID, val 是数据表唯一ID
	equipIndex map[string]*model.EquipmentModel // 索引,方便查找 key 是配置表中的关联ID, val 是数据表唯一ID
	Task       *model.TaskModel
	roleAttr   *global.Attribute
	mainAttr   *global.Attribute
	subAttr    *global.Attribute
	mtx        sync.RWMutex
}

func NewPlayer(uid int, g *Game) *Player {
	return &Player{
		Uid:        uid,
		game:       g,
		useHero:    make(map[string]*model.HeroModel),
		useEquip:   make(map[string]*model.EquipmentModel),
		packIndex:  make(map[string]model.IBackPack),
		heroIndex:  make(map[string]*model.HeroModel),
		equipIndex: make(map[string]*model.EquipmentModel),
	}
}

func (p *Player) LoadDataFromDB() (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
			log.Error(string(debug.Stack()))
		}
	}()

	if p.Account, err = model.GetUserCollection().LoadUser(nil, p.Uid); err != nil {
		return
	}

	if p.Role, err = model.RoleCollection().LoadRoleFromDB(p.Uid); err != nil {
		return
	}

	if p.Stage, err = model.GetStageColl().LoadStageFromDB(p.Uid); err != nil {
		return
	}

	if p.Equip, err = model.GetEquipmentCollection().LoadEquipmentFromDB(p.Uid); err != nil {
		return
	}
	for _, equipmentModel := range p.Equip {
		p.equipIndex[equipmentModel.RelationId] = equipmentModel
		if equipmentModel.Use {
			p.useEquip[equipmentModel.Id] = equipmentModel
		}
	}
	if p.Package, err = model.GetBackPackColl().LoadBackpackFromDB(p.Uid); err != nil {
		return
	}

	for _, pack := range p.Package {
		met := pack.(*model.UpgradeCard)
		p.packIndex[met.TabId] = met
	}

	if p.Hero, err = model.GetHeroCollection().LoadHeroFromDB(p.Uid); err != nil {
		return
	}
	for _, heroModel := range p.Hero {
		p.heroIndex[heroModel.RelationId] = heroModel
		if heroModel.Use {
			p.useHero[heroModel.Id] = heroModel
		}
	}

	//if p.Task,err = model.GetTaskCollection().LoadTaskFromDB(p.Uid) ; err != nil {
	//	return err
	//}else{
	//	flushDayTask(p)
	//}

	return nil
}

type AttrObj struct {
	Obj  interface{}
	Attr *global.Attribute
}

type TotalAttributeHelper struct {
	Main     *AttrObj
	Sub      *AttrObj
	Armor    *AttrObj
	Ornament *AttrObj
	Role     *AttrObj
}

func (p *Player) TotalAttribute() error{
	// 天赋技能 + 角色的基础属性
	pubAttr := p.Role.TalentLv.AddAttribute(p.Role.RLv, *excel.RoleDataConf[excel.GetConfigId(p.Role.RelationId, p.Role.RLv)].Attribute)
	var (
		res = new(TotalAttributeHelper)
	)

	for _, equipment := range p.useEquip {
		if equipment == nil {
			continue
		}
		switch equipment.Type {
		case define.ArmorType: //护甲的属性
			res.Armor = &AttrObj{
				Attr: excel.RoleArmorDataConf[excel.GetConfigId(equipment.RelationId, equipment.Level)].Attribute,
				Obj:  equipment,
			}
		case define.OrnamentType: //饰品的属性
			res.Ornament = &AttrObj{
				Attr: excel.RoleOrnamentsDataConf[excel.GetConfigId(equipment.RelationId, equipment.Level)].Attribute,
				Obj:  equipment,
			}
		case define.MainWeaponType:
			res.Main = &AttrObj{
				Attr: excel.RoleMainWeaponConf[excel.GetConfigId(equipment.RelationId, equipment.Level)].Attribute,
				Obj:  equipment,
			}
		case define.SubWeaponType:
			res.Sub = &AttrObj{
				Attr: excel.RoleSubWeaponConf[excel.GetConfigId(equipment.RelationId, equipment.Level)].Attribute,
				Obj:  equipment,
			}
		}
	}

	if res.Main == nil {
		return  fmt.Errorf("主 武器不存在! ")
	}

	if res.Armor != nil {
		pubAttr = pubAttr.Add(res.Armor.Attr)
	}

	if res.Ornament != nil {
		pubAttr = pubAttr.Add(res.Ornament.Attr)
	}

	res.Role = &AttrObj{
		Obj:  p.Role,
		Attr: pubAttr.Add(res.Main.Attr),
	}

	if p.Role.SkinId != "" {
		res.Role.Attr = res.Role.Attr.Add(excel.RoleSkinDataConf[p.Role.SkinId].Attribute)
	}

	if res.Sub != nil {
		p.roleAttr = res.Role.Attr.Add(res.Sub.Attr).AttributeConst()
		p.subAttr = res.Sub.Attr.Add(pubAttr).AttributeConst()
	} else {
		p.roleAttr = res.Role.Attr.AttributeConst()
	}

	p.mainAttr = res.Main.Attr.Add(pubAttr).AttributeConst()
	return nil
}

// 添加道具
//1 根据类型进行详细操作
//2 金币,钻石,体力 直接增加
//3 卡片 查背包种是否存在相同的卡片,存在就直接增加数量,不存在就直接创建卡片类型并给予获得的数量,写入数据库
//4 主武,副武,护甲,饰品,英雄装备 类型的就检测背包种是否存在这种装备
// 4.1 存在: 根据获得和当前用有的装备的等级大小,将等级低的转为碎片,等级高的装备保留
// 4.2 不存在: 使用一个数量创建这种类型的武器并入库,若还剩下其他数量直接转卡片
func (p *Player) AddAttachments(Attachment []*global.Item) (*response.AttachmentsResp, []bool, error) {

	var rs = &response.AttachmentsResp{
		Replace: make(map[string]*response.ReplaceColl),
	}
	var (
		flush = make([]bool, 5)
		err   error
		where = bson.M{"uid": p.Uid}
	)
	for _, attachment := range Attachment {
		if attachment.Count <= 0 {
			continue
		}
		switch attachment.ItemType {
		case define.ItemExpType:
			_, ok := p.game.playAddExp(p, attachment.Count, true)
			if ok {
				attachment.ItemID = "exp"
				rs.Attachments = append(rs.Attachments, attachment)
				flush[4] = true
			} else {
				err = fmt.Errorf("经验添加失败!!!")
			}
		case define.ItemGoldType:
			p.Account.Gold += attachment.Count
			_, err = model.GetUserCollection().UpdateOne(nil, where, bson.M{"gold": p.Account.Gold})
			attachment.ItemID = "gold"
			rs.Attachments = append(rs.Attachments, attachment)
			flush[0] = true

		case define.ItemDiamondType:
			p.Account.Diamond += attachment.Count
			_, err = model.GetUserCollection().UpdateOne(nil, where, bson.M{"diamond": p.Account.Diamond})
			attachment.ItemID = "diamond"
			rs.Attachments = append(rs.Attachments, attachment)
			flush[0] = true

		case define.ItemStrengthType:
			p.Account.FlushStrength(attachment.Count)
			attachment.ItemID = "strength"
			rs.Attachments = append(rs.Attachments, attachment)
			flush[1] = true

		case define.ItemMainWeaponType:
			equipConf, ok := excel.RoleMainWeaponConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("主武器配置表中数据不存在! ID: %s", attachment.ItemID)
			}
			flush[2], err = p.GotEquip(equipConf, attachment, rs, define.ItemMainWeaponType)

		case define.ItemSubWeaponType:
			equipConf, ok := excel.RoleSubWeaponConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("副武器配置表中数据不存在! ID: %s", attachment.ItemID)
			}
			flush[2], err = p.GotEquip(equipConf, attachment, rs, define.ItemSubWeaponType)

		case define.ItemArmorType:
			equipConf, ok := excel.RoleArmorDataConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("护甲配置表种数据不存在! ID: %s", attachment.ItemID)
			}
			flush[2], err = p.GotEquip(equipConf, attachment, rs, define.ItemArmorType)

		case define.ItemOrnamentsType:
			equipConf, ok := excel.RoleOrnamentsDataConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("饰品配置表中数据不存在! ID: %s", attachment.ItemID)
			}
			flush[2], err = p.GotEquip(equipConf, attachment, rs, define.ItemOrnamentsType)

		case define.ItemHeroEquipType:
			equipConf, ok := excel.HeroEquipmentDataConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("英雄装备配置表中数据不存在! ID: %s", attachment.ItemID)
			}
			flush[3], err = p.GotEquip(equipConf, attachment, rs, define.ItemHeroEquipType)

		case define.ItemCardType:
			if ipack, ok := p.packIndex[attachment.ItemID]; ok {
				pack := ipack.(*model.UpgradeCard)
				pack.Num += attachment.Count
				_, err = model.GetBackPackColl().UpdateOne(nil, bson.M{"_id": pack.Id}, bson.M{"num": pack.Num})
			} else {
				pack := &model.UpgradeCard{
					Id:    xid.New().String(),
					Uid:   p.Uid,
					TabId: attachment.ItemID,
					Num:   attachment.Count,
				}
				err = model.GetBackPackColl().InsertOne(nil, pack)
				if err == nil {
					p.packIndex[attachment.ItemID] = pack
					p.Package[pack.Id] = pack
				}
			}

			rs.Attachments = append(rs.Attachments, attachment)

		case define.ItemHeroType:
			c, ok := excel.HeroDataConf[attachment.ItemID]
			if !ok {
				return nil, nil, fmt.Errorf("英雄配置表中数据不存在! ID: %s", attachment.ItemID)
			}
			flush[3], err = p.GotHero(c, attachment, rs)
		}
		if err != nil {
			fmt.Println("添加道具失败：", err)
			return nil, nil, err
		}
	}

	return rs, flush, nil
}

//4 主武,副武,护甲,饰品,英雄装备 类型的就检测背包种是否存在这种装备
// 4.1 存在: 根据获得和当前用有的装备的等级大小,将等级低的转为碎片,等级高的装备保留
// 4.2 不存在: 使用一个数量创建这种类型的武器并入库,若还剩下其他数量直接转卡片
func (p *Player) GotEquip(tableConf *excel.EquipagePublicConfig, attach *global.Item, rs *response.AttachmentsResp, chestType define.ItemType) (flush bool, err error) {
	ownedEquip, ok := p.equipIndex[tableConf.RelationId]
	if !ok {
		var equipTy define.EquipageType
		switch chestType {
		case define.ItemMainWeaponType:
			equipTy = define.MainWeaponType
		case define.ItemSubWeaponType:
			equipTy = define.SubWeaponType
		case define.ItemArmorType:
			equipTy = define.ArmorType
		case define.ItemOrnamentsType:
			equipTy = define.OrnamentType
		case define.ItemHeroEquipType:
			equipTy = define.HeroEquipType
		}

		equip := model.GetEquipmentCollection().CreateEquipment(p.Uid, tableConf.RelationId, equipTy, tableConf.Index)
		err = model.GetEquipmentCollection().InsertOne(nil, equip)
		if err != nil {
			log.Error("玩家获得装备失败:", err)
			return
		}
		p.Equip[equip.Id] = equip
		p.equipIndex[equip.RelationId] = equip

		rs.Attachments = append(rs.Attachments, &global.Item{
			ItemID:   attach.ItemID,
			Count:    1,
			ItemType: attach.ItemType,
		})

		if attach.Count == 1 {
			return
		}
		// 多个就继续往下处理, 将多的转为碎片
		ownedEquip = nil
		attach.Count--
	}

	var (
		count       int
		cardId, num = utils.UnmarshalItemsKV(tableConf.SplitStarDeplete)
	)
	if num == 0 {
		return
	}

	count = num * attach.Count

	if ownedEquip != nil && ownedEquip.Level < tableConf.Level { // 武器已经存在了,根据等级大小,将等级小的转换为碎片
		oldTableConf := excel.RoleMainWeaponConf[excel.GetConfigId(ownedEquip.RelationId, ownedEquip.Level)] // 旧的装备配置保存起来,准备拆解
		flush = ownedEquip.Use

		ownedEquip.Level = tableConf.Level
		ownedEquip.Quality = tableConf.StarLevel
		_, err = model.GetEquipmentCollection().UpdateOne(nil, bson.M{"_id": ownedEquip.Id}, bson.M{"level": tableConf.Level, "quality": tableConf.StarLevel})
		if err != nil {
			return
		}

		if replace, ok := rs.Replace[ownedEquip.RelationId]; ok {
			replace.NewLv = tableConf.Level
		} else {
			rs.Replace[ownedEquip.RelationId] = &response.ReplaceColl{
				NewLv:      tableConf.Level,
				OldLv:      oldTableConf.Level,
				RelationId: tableConf.RelationId,
				ItemType:   chestType,
			}
		}

		count = num * (attach.Count - 1) // 有一个已经转为低等级的装备,不能与高等级的统一计算卡张数
		_, num = utils.UnmarshalItemsKV(oldTableConf.SplitStarDeplete)
		count += num
	}

	if pack, ok := p.packIndex[cardId]; ok {
		met := pack.(*model.UpgradeCard)
		met.Num += count
		_, err = model.GetBackPackColl().UpdateOne(nil, bson.M{"_id": met.Id}, bson.M{"num": met.Num})
		if err != nil {
			return
		}
	} else {
		card := &model.UpgradeCard{
			Id:    xid.New().String(),
			Uid:   p.Uid,
			TabId: cardId,
			Num:   count,
		}
		err = model.GetBackPackColl().InsertOne(nil, card)
		if err != nil {
			return
		}
		p.packIndex[cardId] = card
		p.Package[card.Id] = card
	}
	rs.Attachments = append(rs.Attachments, &global.Item{
		ItemID:   cardId,
		Count:    count,
		ItemType: define.ItemCardType,
	})

	return
}

func (p *Player) GotHero(tableConf *excel.HeroDataConfig, attach *global.Item, rs *response.AttachmentsResp) (flush bool, err error) {
	ownedHero, ok := p.heroIndex[tableConf.RelationId]

	if !ok { // 为玩家创建这个角色
		psIds := []string{excel.HeroPassiveSkillConf[excel.GetConfigId(tableConf.PSId1, 1)].Id}

		if tableConf.PSId2 != "0" && tableConf.PSId2 != "" {
			psIds = append(psIds, excel.HeroPassiveSkillConf[excel.GetConfigId(tableConf.PSId2, 1)].Id)
		}
		if tableConf.PSId3 != "0" && tableConf.PSId3 != "" {
			psIds = append(psIds, excel.HeroPassiveSkillConf[excel.GetConfigId(tableConf.PSId3, 1)].Id)
		}
		if tableConf.PSId4 != "0" && tableConf.PSId4 != "" {
			psIds = append(psIds, excel.HeroPassiveSkillConf[excel.GetConfigId(tableConf.PSId4, 1)].Id)
		}

		hero := &model.HeroModel{
			Id:          xid.New().String(),
			Uid:         p.Uid,
			RelationId:  tableConf.RelationId,
			Use:         false,
			Rarity:      tableConf.Rarity,
			Lv:          1,
			StarLv:      1,
			Index:       tableConf.Index,
			TotalPower:  0,
			ASId:        excel.HeroActiveSkillConf[excel.GetConfigId(tableConf.ASId, 1)].Id,
			PsIds:       psIds,
			Strengthen:  0,
			EquipmentId: "",
		}

		err = model.GetHeroCollection().InsertOne(nil, hero)
		if err != nil {
			return
		}

		p.Hero[hero.Id] = hero
		p.heroIndex[hero.RelationId] = hero

		rs.Attachments = append(rs.Attachments, &global.Item{
			ItemID:   attach.ItemID,
			Count:    1,
			ItemType: define.ItemHeroType,
		})

		if attach.Count == 1 {
			return
		}
		// 继续将多余的英雄转为 升级,升星 卡
		attach.Count--

	}
	var (
		LvCardId, startCardId                    string
		lvCardCount, startCardCount, lvN, StartN int
	)
	LvCardId, lvN = utils.UnmarshalItemsKV(tableConf.SplitUpDeplete)
	lvCardCount = attach.Count * lvN
	startCardId, StartN = utils.UnmarshalItemsKV(tableConf.SplitStarDeplete)
	startCardCount = attach.Count * StartN

	if ownedHero != nil && tableConf.Level > ownedHero.Lv {
		_, err = model.GetHeroCollection().UpdateOne(nil, bson.M{"_id": ownedHero.Id}, bson.M{"level": tableConf.Level, "starLv": tableConf.StarLevel})
		if err != nil {
			return
		}
		oldConfig := excel.HeroDataConf[excel.GetConfigId(ownedHero.RelationId, ownedHero.Lv)] // 被替换的英雄配置缓存起来准备拆解

		ownedHero.Lv = tableConf.Level
		ownedHero.StarLv = tableConf.StarLevel
		flush = ownedHero.Use

		// 扣除一个的数量
		lvCardCount -= lvN
		startCardCount -= StartN

		// 加上一个替换下来的英雄拆解的数量
		_, n := utils.UnmarshalItemsKV(oldConfig.SplitUpDeplete)
		lvCardCount += n
		_, n = utils.UnmarshalItemsKV(oldConfig.SplitStarDeplete)
		startCardCount += n

		if replace, ok := rs.Replace[oldConfig.RelationId]; ok {
			replace.NewLv = tableConf.Level
		} else {
			rs.Replace[oldConfig.RelationId] = &response.ReplaceColl{
				NewLv:      oldConfig.Level,
				OldLv:      tableConf.Level,
				RelationId: tableConf.RelationId,
				ItemType:   define.ItemHeroEquipType,
			}
		}
	}

AddCard:
	if icard, ok := p.packIndex[LvCardId]; ok { // 卡处理
		card := icard.(*model.UpgradeCard)
		card.Num += lvCardCount
		if n, e := model.GetBackPackColl().UpdateOne(nil, bson.M{"_id": card.Id}, bson.M{"num": card.Num}); e != nil || n == 0 {
			err = e
			return
		}
		rs.Attachments = append(rs.Attachments, &global.Item{
			ItemID:   LvCardId,
			Count:    lvCardCount,
			ItemType: define.ItemCardType,
		})
	} else {
		card := &model.UpgradeCard{
			Id:    xid.New().String(),
			Uid:   p.Uid,
			TabId: LvCardId,
			Num:   lvCardCount,
		}
		err = model.GetBackPackColl().InsertOne(nil, card)
		if err != nil {
			return
		}
		p.Package[card.Id] = card
		p.packIndex[card.TabId] = card
		rs.Attachments = append(rs.Attachments, &global.Item{
			ItemID:   LvCardId,
			Count:    lvCardCount,
			ItemType: define.ItemCardType,
		})
	}
	if startCardCount != 0 && LvCardId != startCardId {
		LvCardId = startCardId
		lvCardCount = startCardCount

		goto AddCard
	}

	return
}
