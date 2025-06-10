package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/AlexEagle1535/market-rent-bot/db"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	userStates   = make(map[int64]string)   // состояние пользователя
	userData     = make(map[int64][]string) // данные пользователей
	stateMutex   = &sync.Mutex{}            // мьютекс для безопасного доступа
	stateDefault = "main_menu"              // дефолтное состояние
)

func setState(userID int64, state string) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	userStates[userID] = state
}

func getState(userID int64) string {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	state, ok := userStates[userID]
	if !ok {
		return stateDefault
	}
	return state
}

func sendMenu(conn *sql.DB, ctx *th.Context, msg telego.Message) error {
	role, err := db.GetUserRole(conn, msg.From.ID, msg.From.Username)
	if err != nil {
		log.Printf("Ошибка получения роли пользователя %d: %v", msg.From.ID, err)
		return err
	}
	switch role {
	case "admin":
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"🔐 Админ меню",
		).WithReplyMarkup(adminMenu()))
	case "tenant":
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"👤 Меню арендатора",
		).WithReplyMarkup(tenantMenu()))
	default:
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(msg.Chat.ID),
			"🚫 У вас нет доступа к боту.",
		))
	}
	return nil
}

func main() {
	_ = godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	// Чтение и преобразование строки в слайс
	var admins = strings.Split(os.Getenv("ADMINS"), ",")
	conn := db.InitDB()
	defer conn.Close()
	for _, admin := range admins {
		admin = strings.TrimSpace(admin) // Удаляем лишние пробелы
		if admin == "" {
			log.Fatal("Пустой username в списке админов")
		}
		err := db.SetUserRole(conn, 0, admin, "admin")
		if err != nil {
			log.Fatal("Ошибка добавления пользователя:", err)
		}
	}

	// Проверка, является ли username админом
	ctx := context.Background()
	bot, err := telego.NewBot(token, telego.WithDefaultLogger(true, true))
	if err != nil {
		log.Fatal(err)
	}

	// Отключаем Webhook (если вдруг)
	_ = bot.DeleteWebhook(ctx, nil)

	// Получаем канал обновлений
	updates, err := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{Timeout: 60})
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем обработчик
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	defer bh.Stop()

	// Обработка /start
	bh.HandleMessage(func(ctx *th.Context, msg telego.Message) error {
		err := sendMenu(conn, ctx, msg)
		if err != nil {
			return err
		}
		return nil
	}, th.CommandEqual("start"))

	// Обработка callback-кнопок
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		message := query.Message.Message()
		if message == nil {
			return nil
		}

		var newText string
		var newMarkup *telego.InlineKeyboardMarkup

		switch query.Data {
		// === Общее ===
		case "go_back":
			role, err := db.GetUserRole(conn, query.From.ID, query.From.Username)
			if err != nil {
				log.Printf("Ошибка получения роли пользователя %d: %v", query.From.ID, err)
				return err
			}
			setState(query.From.ID, "main_menu")
			if role == "admin" {
				newText = "🔐 Админ меню"
				newMarkup = adminMenu()
			} else {
				newText = "👤 Меню арендатора"
				newMarkup = tenantMenu()
			}

		// === Админ ===
		case "admin_tenants":
			newText = "🧑‍💼 Раздел: Арендаторы"
			newMarkup = tenantAdminMenu()

		case "admin_broadcast":
			newText = "📢 Введите текст рассылки (заглушка)."
			newMarkup = backButton()

		case "admin_approvals":
			newText = "✅ Задачи на подтверждение (заглушка)."
			newMarkup = backButton()

		// === Существующие ===
		case "import_csv":
			newText = "📥 Загрузка арендаторов из CSV (заглушка)."
			newMarkup = backButton()

		case "list_tenants":
			newText = "📋 Список арендаторов:\n1. ИП Иванов\n2. ООО Рынок\n(заглушка)"
			newMarkup = backButton()

		case "admin_users":
			newText = "👤 Пользователи системы"
			newMarkup = usersAdminMenu()

		case "add_user":
			newText = "➕ Добавление пользователя"
			newMarkup = addUserMenu()

		case "add_admin":
			newText = "Введите username нового администратора:"
			setState(query.From.ID, "awaiting_admin_data")
			newMarkup = backButton()

		case "add_tenant":
			newText = "Введите username нового арендатора:"
			setState(query.From.ID, "awaiting_tenant_data")
			newMarkup = backButton()
		}
		// Редактируем сообщение
		_, _ = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: message.Chat.ID},
			MessageID: message.MessageID,
			Text:      newText,
			ReplyMarkup: &telego.InlineKeyboardMarkup{
				InlineKeyboard: newMarkup.InlineKeyboard,
			},
		})

		// Ответ на callback (чтобы убрать "часики" у пользователя)
		_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))

		return nil
	}, th.AnyCallbackQueryWithMessage())

	// обработчик тексовых событий
	bh.HandleMessage(func(ctx *th.Context, msg telego.Message) error {
		userID := msg.From.ID
		state := getState(userID)

		if state == "awaiting_admin_data" {
			username := msg.Text
			err := db.SetUserRole(conn, 0, username, "admin")
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении админа."))

			} else {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Админ добавлен!"))
				setState(userID, "main_menu")
			}
		}
		if state == "awaiting_tenant_data" {
			username := msg.Text
			err := db.SetUserRole(conn, 0, username, "tenant")
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении арендатора."))
			} else {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Арендатор добавлен!"))
				setState(userID, "main_menu")
			}
		}
		sendMenu(conn, ctx, msg)
		return nil
	})
	// Запускаем обработчик
	go func() {
		if err := bh.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Блокируем main
	select {}
}

func adminMenu() *telego.InlineKeyboardMarkup {
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

func tenantAdminMenu() *telego.InlineKeyboardMarkup {
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

func usersAdminMenu() *telego.InlineKeyboardMarkup {
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

func addUserMenu() *telego.InlineKeyboardMarkup {
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

func tenantMenu() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("💳 Платежи").WithCallbackData("tenant_payments"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📄 Мой договор").WithCallbackData("tenant_contract"),
		),
	)
}

func backButton() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("go_back"),
		),
	)
}
