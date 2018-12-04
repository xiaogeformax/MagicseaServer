package session

type SessionManager struct{
	players map[uint64]*PlayerSession
	name2players map[string]*PlayerSession //名字索引
}