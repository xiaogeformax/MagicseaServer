package session

type SessionManager struct{
	players map[uint64]*PlayerSession
	name2players map[string]*PlayerSession //名字索引
}
func NewSessionManager() *SessionManager {
	sm:= new (SessionManager)
	sm.players = make(map[uint64]*PlayerSession)
	sm.name2players = make(map[string]*PlayerSession)
	return sm
}