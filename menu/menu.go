package menu

import (
	"fmt"

	"github.com/AlexEagle1535/market-rent-bot/db"
	"github.com/AlexEagle1535/market-rent-bot/states"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Admin() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📢 Рассылка").WithCallbackData("admin_broadcast"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("✅ Подтверждения").WithCallbackData("admin_approvals"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("👨‍💼 Арендаторы").WithCallbackData("admin_tenants"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🧺 Рынок").WithCallbackData("admin_market"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("👤 Пользователи системы").WithCallbackData("admin_users"),
		),
	)
}

func AdminTenants() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📥 Импорт арендаторов из файла").WithCallbackData("import_csv"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📋 Список арендаторов").WithCallbackData("list_tenants"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}

func AdminTenantsList(tanants []db.Tenant, state *states.ListState) *telego.InlineKeyboardMarkup {
	const pageSize = 10
	page := state.Page
	start := page * pageSize
	if start >= len(tanants) {
		page = 0
		state.Page = 0
		start = 0
	}
	end := start + pageSize
	if end > len(tanants) {
		end = len(tanants)
	}
	slice := tanants[start:end]

	rows := make([][]telego.InlineKeyboardButton, 0)

	// Кнопка "Добавить арендатора" и поиск
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить").WithCallbackData("add_tenant"),
		tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search_tenant"),
	))

	// Вывод списка арендаторов
	for _, t := range slice {
		label := fmt.Sprintf("%s", t.FullName)
		viewBtn := tu.InlineKeyboardButton(label).WithCallbackData(fmt.Sprintf("view_tenant:%d", t.ID))
		rows = append(rows, tu.InlineKeyboardRow(viewBtn))
	}

	// Постраничный вывод
	pageRow := []telego.InlineKeyboardButton{}
	if page > 0 {
		pageRow = append(pageRow, tu.InlineKeyboardButton("⬅️ Предыдущая страница").WithCallbackData("page_prev"))
	}
	if end < len(tanants) {
		pageRow = append(pageRow, tu.InlineKeyboardButton("➡️ Следующая страница").WithCallbackData("page_next"))
	}
	if len(pageRow) > 0 {
		rows = append(rows, pageRow)
	}

	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 В меню").WithCallbackData("go_back"),
	))

	return tu.InlineKeyboard(rows...)
}

func AdminMarket() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🏪 Павильоны").WithCallbackData("pavilions"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🧑‍🌾 Виды деятельности").WithCallbackData("activity_types"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}

func AdminPavilionList(pavilions []db.Pavilion, state *states.ListState) *telego.InlineKeyboardMarkup {
	const pageSize = 10
	page := state.Page
	start := page * pageSize
	if start >= len(pavilions) {
		page = 0
		state.Page = 0
		start = 0
	}
	end := start + pageSize
	if end > len(pavilions) {
		end = len(pavilions)
	}
	slice := pavilions[start:end]

	rows := make([][]telego.InlineKeyboardButton, 0)

	// Кнопка "Добавить павильон" и поиск
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить").WithCallbackData("add_pavilion"),
		tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search_pavilion"),
	))

	// Вывод списка павильонов
	for _, p := range slice {
		label := fmt.Sprintf("№ %s", p.Number)
		viewBtn := tu.InlineKeyboardButton(label).WithCallbackData(fmt.Sprintf("view_pavilion:%d", p.ID))
		rows = append(rows, tu.InlineKeyboardRow(viewBtn))
	}

	// Постраничный вывод
	pageRow := []telego.InlineKeyboardButton{}
	if page > 0 {
		pageRow = append(pageRow, tu.InlineKeyboardButton("⬅️ Предыдущая страница").WithCallbackData("page_prev"))
	}
	if end < len(pavilions) {
		pageRow = append(pageRow, tu.InlineKeyboardButton("➡️ Следующая страница").WithCallbackData("page_next"))
	}
	if len(pageRow) > 0 {
		rows = append(rows, pageRow)
	}

	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 В меню").WithCallbackData("go_back"),
	))

	return tu.InlineKeyboard(rows...)
}

func AdminActivityTypesList(types []db.ActivityType, state *states.ListState) *telego.InlineKeyboardMarkup {
	const pageSize = 10
	page := state.Page
	start := page * pageSize
	if start >= len(types) {
		page = 0
		state.Page = 0
		start = 0
	}
	end := start + pageSize
	if end > len(types) {
		end = len(types)
	}
	slice := types[start:end]

	rows := make([][]telego.InlineKeyboardButton, 0)

	// Кнопка "Добавить вид деятельности" и поиск
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить").WithCallbackData("add_activity_type"),
		tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search_activity_type"),
	))

	// Вывод списка видов деятельности
	for _, t := range slice {
		label := t.Name
		viewBtn := tu.InlineKeyboardButton(label).WithCallbackData("noop")
		rows = append(rows, tu.InlineKeyboardRow(viewBtn))
	}

	// Постраничный вывод
	pageRow := []telego.InlineKeyboardButton{}
	if page > 0 {
		pageRow = append(pageRow, tu.InlineKeyboardButton("⬅️ Предыдущая страница").WithCallbackData("page_prev"))
	}
	if end < len(types) {
		pageRow = append(pageRow, tu.InlineKeyboardButton("➡️ Следующая страница").WithCallbackData("page_next"))
	}
	if len(pageRow) > 0 {
		rows = append(rows, pageRow)
	}

	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 В меню").WithCallbackData("go_back"),
	))

	return tu.InlineKeyboard(rows...)
}

func AdminActivityTypeSelect(types []db.ActivityType, selected map[int]bool) *telego.InlineKeyboardMarkup {
	rows := make([][]telego.InlineKeyboardButton, 0)

	for _, t := range types {
		label := t.Name
		if selected[t.ID] {
			label = "✅ " + label
		}
		btn := tu.InlineKeyboardButton(label).
			WithCallbackData(fmt.Sprintf("select_activity_type:%d", t.ID))
		rows = append(rows, tu.InlineKeyboardRow(btn))
	}

	// Кнопка завершения выбора
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить").WithCallbackData("add_activity_type"),
		tu.InlineKeyboardButton("✅ Завершить").WithCallbackData("finish_activity_selection"),
	))

	return tu.InlineKeyboard(rows...)
}

func AdminPavilionSelect(pavilions []db.Pavilion) *telego.InlineKeyboardMarkup {
	rows := make([][]telego.InlineKeyboardButton, 0)

	for _, p := range pavilions {
		label := fmt.Sprintf("№ %s", p.Number)
		btn := tu.InlineKeyboardButton(label).
			WithCallbackData(fmt.Sprintf("select_pavilion:%d", p.ID))
		rows = append(rows, tu.InlineKeyboardRow(btn))
	}

	// Кнопка завершения выбора
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить").WithCallbackData("add_pavilion"),
	))

	return tu.InlineKeyboard(rows...)
}

func AdminUsers() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📥 Импорт пользователей из файла").WithCallbackData("import_csv"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("👥 Список пользователей").WithCallbackData("list_users"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Добавить пользователя").WithCallbackData("add_user"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}

func AddUser() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Добавить админа").WithCallbackData("add_admin"),
		),
		// tu.InlineKeyboardRow(
		// 	tu.InlineKeyboardButton("Добавить арендатора").WithCallbackData("add_tenant"),
		// ),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("admin_users"),
		),
	)
}

func AdminUserList(users []db.User, state *states.ListState) *telego.InlineKeyboardMarkup {
	const pageSize = 10
	page := state.Page
	search := state.Search
	start := page * pageSize
	if start >= len(users) {
		page = 0
		state.Page = 0
		start = 0
	}
	end := start + pageSize
	if end > len(users) {
		end = len(users)
	}
	slice := users[start:end]

	rows := make([][]telego.InlineKeyboardButton, 0)

	// Фильтры
	filterRow := tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("👥 Все").WithCallbackData("filter:all"),
		tu.InlineKeyboardButton("🛡 Админы").WithCallbackData("filter:admin"),
		tu.InlineKeyboardButton("🏠 Арендаторы").WithCallbackData("filter:tenant"),
	)
	rows = append(rows, filterRow)

	// Поиск
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔍 Поиск").WithCallbackData("search_user"),
	))

	// Список пользователей
	for _, u := range slice {
		username := "null"
		if u.Username.Valid {
			username = u.Username.String
		}
		telegramID := "0"
		if u.TelegramID.Valid {
			telegramID = fmt.Sprintf("%d", u.TelegramID.Int64)
		}

		label := fmt.Sprintf("%s - %s", username, u.Role)
		userBtn := tu.InlineKeyboardButton(label).WithCallbackData("noop")
		delBtn := tu.InlineKeyboardButton("❌").WithCallbackData(fmt.Sprintf("confirm_delete:%s:%s", telegramID, username))
		rows = append(rows, tu.InlineKeyboardRow(userBtn, delBtn))
	}

	// Пагинация
	pagination := []telego.InlineKeyboardButton{}
	if page > 0 {
		pagination = append(pagination, tu.InlineKeyboardButton("⬅️ Предыдущая страница").WithCallbackData("page_prev"))
	}
	if end < len(users) {
		pagination = append(pagination, tu.InlineKeyboardButton("➡️ Следующая страница").WithCallbackData("page_next"))
	}
	if len(pagination) > 0 {
		rows = append(rows, pagination)
	}

	if search != "" {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔄 Сбросить поиск").WithCallbackData("reset_search"),
		))
	}

	// Назад
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 В меню").WithCallbackData("admin_users"),
	))

	return tu.InlineKeyboard(rows...)
}

func ConfirmDeleteUser(telegramID, username string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("✅ Подтвердить").WithCallbackData(fmt.Sprintf("delete_user:%s:%s", telegramID, username)),
			tu.InlineKeyboardButton("❌ Отмена").WithCallbackData("list_users"),
		),
	)
}

func Tenant() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("💳 Платежи").WithCallbackData("tenant_payments"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📄 Мой договор").WithCallbackData("tenant_contract"),
		),
	)
}

func OkButton(data string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("✅ ОК").WithCallbackData(data),
		),
	)
}

func BackButton(data string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData(data),
		),
	)
}
