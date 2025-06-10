package states

import "sync"

var (
	userStates   = make(map[int64]string)   // состояние пользователя
	userData     = make(map[int64][]string) // данные пользователей
	stateMutex   = &sync.Mutex{}            // мьютекс для безопасного доступа
	stateDefault = "main_menu"              // дефолтное состояние
)

func Set(userID int64, state string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	userStates[userID] = state
}

func Get(userID int64) string {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	state, ok := userStates[userID]
	if !ok {
		return stateDefault
	}
	return state
}
