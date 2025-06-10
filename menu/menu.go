package menu

import (
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
			tu.InlineKeyboardButton("📋 Список пользователей").WithCallbackData("list_tenants"),
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

func BackButton() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}
