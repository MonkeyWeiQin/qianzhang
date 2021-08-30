package model

import (
	"battle_rabbit/define"
	"battle_rabbit/excel"
	"battle_rabbit/global"
)

func RewardToItems(reward *excel.RewardDataConfig) (items []*global.Item) {
	items = append(items, &global.Item{
		Count:    reward.Gold,
		ItemType: define.ItemGoldType,
	})

	if reward.Diamond > 0 {
		items = append(items, &global.Item{
			Count:    reward.Diamond,
			ItemType: define.ItemDiamondType,
		})
	}

	if reward.PhysicalPower > 0 {
		items = append(items, &global.Item{
			Count:    reward.PhysicalPower,
			ItemType: define.ItemStrengthType,
		})
	}

	if reward.Exp > 0 {
		items = append(items, &global.Item{
			Count:    reward.Exp,
			ItemType: define.ItemExpType,
		})

	}

	if reward.Hero != "0" {
		items = append(items, global.DisassembleToItem(reward.Hero, define.ItemHeroType)...)
	}

	if reward.Card != "0" {
		items = append(items, global.DisassembleToItem(reward.Card, define.ItemCardType)...)
	}
	if reward.MainWeapon != "0" {
		items = append(items, global.DisassembleToItem(reward.MainWeapon, define.ItemMainWeaponType)...)
	}

	if reward.SubWeapon != "0" {
		items = append(items, global.DisassembleToItem(reward.SubWeapon, define.ItemSubWeaponType)...)
	}

	if reward.Armor != "0" {
		items = append(items, global.DisassembleToItem(reward.Armor, define.ItemArmorType)...)
	}

	if reward.Ornament != "0" {
		items = append(items, global.DisassembleToItem(reward.Ornament, define.ItemOrnamentsType)...)
	}

	if reward.HeroEquip != "0" {
		items = append(items, global.DisassembleToItem(reward.HeroEquip, define.ItemHeroEquipType)...)
	}
	return
}