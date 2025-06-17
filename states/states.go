package states

import (
	"sync"
)

// --- Простые состояния пользователей (одна строка) ---
var (
	userStates   = make(map[int64]string)   // состояние пользователя
	userData     = make(map[int64][]string) // данные пользователей
	stateDefault = "main_menu"
)

var stateMutex sync.RWMutex

func Set(userID int64, state string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	userStates[userID] = state
}

func Get(userID int64) string {
	stateMutex.RLock()
	defer stateMutex.RUnlock()
	state, ok := userStates[userID]
	if !ok {
		return stateDefault
	}
	return state
}

func SetData(userID int64, data []string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	userData[userID] = data
}

func GetData(userID int64) []string {
	stateMutex.RLock()
	defer stateMutex.RUnlock()
	return userData[userID]
}

func Clear(userID int64) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	delete(userStates, userID)
	delete(userData, userID)
}

// --- Состояние листа пользователей (UserListState) ---

type UserListState struct {
	Page   int
	Filter string // all, admin, tenant
	Search string
}

var (
	userListStates = make(map[int64]*UserListState)
	userListMutex  sync.RWMutex
)

func GetUserListState(userID int64) *UserListState {
	userListMutex.RLock()
	state, ok := userListStates[userID]
	userListMutex.RUnlock()

	if ok {
		return state
	}

	userListMutex.Lock()
	defer userListMutex.Unlock()
	state = &UserListState{
		Page:   0,
		Filter: "all",
		Search: "",
	}
	userListStates[userID] = state
	return state
}

func UpdateUserListState(userID int64, update func(state *UserListState)) {
	userListMutex.Lock()
	defer userListMutex.Unlock()
	if state, ok := userListStates[userID]; ok {
		update(state)
	}
}

func ClearUserListState(userID int64) {
	userListMutex.Lock()
	defer userListMutex.Unlock()
	delete(userListStates, userID)
}
