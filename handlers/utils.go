package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AlexEagle1535/market-rent-bot/db"
	"github.com/AlexEagle1535/market-rent-bot/menu"
	"github.com/AlexEagle1535/market-rent-bot/states"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendMenu(ctx *th.Context, msg telego.Message) error {
	role, err := db.GetUserRole(msg.From.ID, msg.From.Username)
	if err != nil {
		log.Printf("Ошибка получения роли пользователя %d: %v", msg.From.ID, err)
		return err
	}
	switch role {
	case "admin":
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"🔐 Админ меню",
		).WithReplyMarkup(menu.Admin()))
	case "tenant":
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"👤 Меню арендатора",
		).WithReplyMarkup(menu.Tenant()))
	default:
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"🚫 У вас нет доступа к боту.",
		))
	}
	return nil
}

func buildUserList(userID int64) (string, *telego.InlineKeyboardMarkup, error) {
	state := states.GetListState(userID)

	var users []db.User
	var err error

	// Применяем фильтры
	if state.Search != "" {
		users, err = db.SearchUsers(state.Search, state.Filter)
	} else {
		switch state.Filter {
		case "admin":
			users, err = db.GetUsersByRole("admin")
		case "tenant":
			users, err = db.GetUsersByRole("tenant")
		default:
			users, err = db.GetAllUsers()
		}
	}
	if err != nil {
		return "", nil, err
	}

	text := "👥 Список пользователей"
	if state.Search != "" {
		text += fmt.Sprintf("\nРезультаты поиска по: %s", state.Search)
	}
	markup := menu.AdminUserList(users, state)
	return text, markup, nil
}

func toggleActivitySelection(userID int64, selectedID int) []int {
	// Получаем текущий список
	raw, _ := states.GetTemp(userID, "selected_activity_ids")
	selected := make(map[int]bool)

	for _, idStr := range strings.Split(raw, ",") {
		if idStr == "" {
			continue
		}
		id, _ := strconv.Atoi(idStr)
		selected[id] = true
	}

	if selected[selectedID] {
		delete(selected, selectedID)
	} else {
		selected[selectedID] = true
	}

	// Сохраняем обновлённый список
	var newList []string
	for id := range selected {
		newList = append(newList, strconv.Itoa(id))
	}
	states.SetTemp(userID, "selected_activity_ids", strings.Join(newList, ","))

	// Возвращаем для использования в UI
	var result []int
	for id := range selected {
		result = append(result, id)
	}
	return result
}
