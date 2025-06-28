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

type ListState struct {
	Page   int
	Filter string
	Search string
	Scope  string // users, pavilions, etc.
}

var (
	listStates   = make(map[int64]*ListState)
	listStateMux sync.RWMutex
)

func GetListState(userID int64) *ListState {
	listStateMux.RLock()
	state, ok := listStates[userID]
	listStateMux.RUnlock()

	if ok {
		return state
	}

	listStateMux.Lock()
	defer listStateMux.Unlock()
	state = &ListState{
		Page:   0,
		Filter: "all",
		Search: "",
		Scope:  "users",
	}
	listStates[userID] = state
	return state
}

func UpdateListState(userID int64, update func(state *ListState)) {
	listStateMux.Lock()
	defer listStateMux.Unlock()
	if state, ok := listStates[userID]; ok {
		update(state)
	}
}

var tempStorage = make(map[int64]map[string]string)

func SetTemp(userID int64, key, value string) {
	if _, ok := tempStorage[userID]; !ok {
		tempStorage[userID] = make(map[string]string)
	}
	tempStorage[userID][key] = value
}

func GetTemp(userID int64, key string) (string, bool) {
	if userData, ok := tempStorage[userID]; ok {
		val, ok := userData[key]
		return val, ok
	}
	return "", false
}

func ClearTemp(userID int64) {
	delete(tempStorage, userID)
}
