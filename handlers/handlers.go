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

func Start(ctx *th.Context, msg telego.Message) error {
	err := sendMenu(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func CallbackQuery(ctx *th.Context, query telego.CallbackQuery) error {
	if query.Message == nil {
		log.Println("⚠️ CallbackQuery.Message is nil — возможно это от inline-кнопки без сообщения")
		return nil
	}
	message := query.Message.Message()
	if message == nil {
		return nil
	}

	var newText string
	var newMarkup *telego.InlineKeyboardMarkup

	switch {
	// === Общее ===
	case query.Data == "go_back":
		role, err := db.GetUserRole(query.From.ID, query.From.Username)
		if err != nil {
			log.Printf("Ошибка получения роли пользователя %d: %v", query.From.ID, err)
			return err
		}
		states.Set(query.From.ID, "main_menu")
		if role == "admin" {
			newText = "🔐 Админ меню"
			newMarkup = menu.Admin()
		} else {
			newText = "👤 Меню арендатора"
			newMarkup = menu.Tenant()
		}

	// === Админ ===
	case query.Data == "admin_tenants":
		newText = "🧑‍💼 Раздел: Арендаторы"
		newMarkup = menu.AdminTetants()

	case query.Data == "admin_broadcast":
		newText = "📢 Введите текст рассылки (заглушка)."
		newMarkup = menu.BackButton()

	case query.Data == "admin_approvals":
		newText = "✅ Задачи на подтверждение (заглушка)."
		newMarkup = menu.BackButton()

	// === Существующие ===
	case query.Data == "import_csv":
		newText = "📥 Загрузка арендаторов из CSV (заглушка)."
		newMarkup = menu.BackButton()

	case query.Data == "list_tenants":
		newText = "📋 Список арендаторов:\n1. ИП Иванов\n2. ООО Рынок\n(заглушка)"
		newMarkup = menu.BackButton()

	case query.Data == "admin_users":
		newText = "👤 Пользователи системы"
		newMarkup = menu.AdminUsers()

	case query.Data == "add_user":
		newText = "➕ Добавление пользователя"
		newMarkup = menu.AddUser()

	case query.Data == "add_admin":
		newText = "Введите username нового администратора:"
		states.Set(query.From.ID, "awaiting_admin_data")
		newMarkup = menu.BackButton()

	case query.Data == "add_tenant":
		newText = "Введите username нового арендатора:"
		states.Set(query.From.ID, "awaiting_tenant_data")
		newMarkup = menu.BackButton()

	case query.Data == "list_users":
		users, err := db.GetAllUsers()
		if err != nil {
			return err
		}
		newText = "👥 Список пользователей"
		newMarkup = menu.AdminUserList(users)

	case strings.HasPrefix(query.Data, "confirm_delete:"):
		data := strings.Split(query.Data, ":")
		var msgOutput string
		if data[1] == "0" {
			msgOutput = data[2]
		} else {
			msgOutput = data[1]
		}
		newText = fmt.Sprintf("Вы уверены, что хотите удалить пользователя %s?", msgOutput)
		newMarkup = menu.ConfirmDeleteUser(data[1], data[2])

	case strings.HasPrefix(query.Data, "delete_user:"):
		data := strings.Split(query.Data, ":")
		telegramID, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			log.Printf("Ошибка преобразования ID пользователя %s: %v", data[1], err)
			return err
		}
		username := data[2]
		err = db.DeleteUser(telegramID, username)
		var msg string
		if err != nil {
			msg = fmt.Sprintf("Ошибка удаления пользователя %s: %v", username, err)
			log.Printf(msg)
		} else {
			msg = fmt.Sprintf("Пользователь %s успешно удалён", username)
		}
		newText = msg
		newMarkup = menu.OkButton("list_users")

	default:
		log.Printf("Неизвестный callback: %s", query.Data)
		// Ответ на callback (чтобы убрать "часики" у пользователя)
		_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
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
}

func TextMessage(ctx *th.Context, msg telego.Message) error {
	userID := msg.From.ID
	state := states.Get(userID)

	if state == "awaiting_admin_data" {
		username := msg.Text
		err := db.SetUserRole(0, username, "admin")
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении админа."))

		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Админ добавлен!"))
			states.Set(userID, "main_menu")
		}
	}
	if state == "awaiting_tenant_data" {
		username := msg.Text
		err := db.SetUserRole(0, username, "tenant")
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении арендатора."))
		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Арендатор добавлен!"))
			states.Set(userID, "main_menu")
		}
	}
	sendMenu(ctx, msg)
	return nil
}
