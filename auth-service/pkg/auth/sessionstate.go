package auth

type LoginSessionState struct {
	Modules     []LoginSessionStateModuleInfo
	SharedState map[string]string
	UserId      string
	SessionId   string
}

type LoginSessionStateModuleInfo struct {
	Id          string
	Type        string
	Properties  map[string]interface{}
	State       ModuleState
	SharedState map[string]string
}

func (l *LoginSessionState) UpdateModuleInfo(mIndex int, mInfo LoginSessionStateModuleInfo) {
	l.Modules[mIndex] = mInfo
}

type ModuleState int

const (
	Fail ModuleState = -1 + iota
	Start
	InProgress //callbacks requested
	Pass
)

const (
	AuthCookieName    = "BlazewallAuthSession"
	SessionCookieName = "BlazewallSession"
)
