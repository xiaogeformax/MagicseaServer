package db


type BattleState int32

const (
	BattleStateFree  BattleState = 0
	BattleStateQueue BattleState = 1
	BattleStateFight BattleState = 2
	BattleStateEnd   BattleState = 3
)

//获取玩家战斗信息
func GetPlayerBattleInfo(id uint64) *PlayerBattleInfo {
	info := &PlayerBattleInfo{}
	if b, _ := GetRedisObject(info, id, GetRedisBattle()); b {
		return info
	}
	return nil
}

//设置玩家战斗信息
func SetPlayerBattleInfo(id uint64, info *PlayerBattleInfo) {
	SetRedisObject(info, id, GetRedisBattle())
}

//设置玩家战斗结束
func SetPlayerBattleState(id uint64, state int32) {
	info := &PlayerBattleInfo{BattleState: state}
	UpdateRedisObjectFields(info, id, GetRedisBattle(), "BattleState")
}
