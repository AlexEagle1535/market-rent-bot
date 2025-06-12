package menu

import (
	"fmt"

	"github.com/AlexEagle1535/market-rent-bot/db"
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
			tu.InlineKeyboardButton("👤 Пользователи системы").WithCallbackData("admin_users"),
		),
	)
}

func AdminTetants() *telego.InlineKeyboardMarkup {
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
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Добавить арендатора").WithCallbackData("add_tenant"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("admin_users"),
		),
	)
}

func AdminUserList(users []db.User) *telego.InlineKeyboardMarkup {

	rows := make([][]telego.InlineKeyboardButton, 0)
	for _, u := range users {
		username := "null"
		if u.Username.Valid {
			username = u.Username.String
		}
		telegramID := "0"
		if u.TelegramID.Valid {
			telegramID = fmt.Sprintf("%d", u.TelegramID.Int64)
		}
		var label string
		if username != "null" {
			label = fmt.Sprintf("%s - %s", username, u.Role)
		} else {
			label = fmt.Sprintf("%s - %s", telegramID, u.Role)
		}
		userBtn := tu.InlineKeyboardButton(label).WithCallbackData("noop")
		delBtn := tu.InlineKeyboardButton("❌").WithCallbackData(fmt.Sprintf("confirm_delete:%s:%s", telegramID, username))
		rows = append(rows, []telego.InlineKeyboardButton{userBtn, delBtn})
	}
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("admin_users"),
	))
	keyboard := tu.InlineKeyboard(rows...)
	return keyboard
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

func BackButton() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}
