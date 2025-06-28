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
	var err error

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
		newMarkup = menu.AdminTenants()

	case query.Data == "admin_broadcast":
		newText = "📢 Введите текст рассылки (заглушка)."
		newMarkup = menu.BackButton("go_back")

	case query.Data == "admin_approvals":
		newText = "✅ Задачи на подтверждение (заглушка)."
		newMarkup = menu.BackButton("go_back")

	// === Существующие ===
	case query.Data == "import_csv":
		newText = "📥 Загрузка арендаторов из CSV (заглушка)."
		newMarkup = menu.BackButton("go_back")

	// case query.Data == "list_tenants":
	// 	newText = "📋 Список арендаторов:\n1. ИП Иванов\n2. ООО Рынок\n(заглушка)"
	// 	newMarkup = menu.BackButton("go_back")

	case query.Data == "admin_users":
		newText = "👤 Пользователи системы"
		newMarkup = menu.AdminUsers()

	case query.Data == "add_user":
		newText = "➕ Добавление пользователя"
		newMarkup = menu.AddUser()

	case query.Data == "add_admin":
		newText = "Введите username нового администратора:"
		states.Set(query.From.ID, "awaiting_admin_data")
		newMarkup = menu.BackButton("go_back")

	case query.Data == "add_tenant":
		newText = "Введите username нового арендатора:"
		states.Set(query.From.ID, "awaiting_tenant_data")
		newMarkup = menu.BackButton("go_back")

	case query.Data == "list_users":
		states.UpdateListState(query.From.ID, func(s *states.ListState) {
			s.Scope = "users"
			s.Page = 0
		})
		newText, newMarkup, err = buildUserList(query.From.ID)
		if err != nil {
			return err
		}

	case query.Data == "admin_market":
		newText = "🧺 Рынок"
		newMarkup = menu.AdminMarket()

	case query.Data == "pavilions":
		var pavilions []db.Pavilion
		pavilions, err = db.GetAllPavilions()
		if err != nil {
			log.Printf("Ошибка получения списка павильонов: %v", err)
		}
		newText = "🏪 Павильоны"
		state := states.GetListState(query.From.ID) // Инициализируем состояние списка, если оно не было создано
		state.Scope = "pavilions"
		state.Page = 0
		newMarkup = menu.AdminPavilionList(pavilions, state)

	case query.Data == "add_pavilion":
		newText = "Введите номер нового павильона:"
		states.Set(query.From.ID, "adding_pavilion_number")
		newMarkup = menu.BackButton("pavilions")

	case query.Data == "list_tenants":
		var tenants []db.Tenant
		tenants, err = db.GetAllTenants()
		if err != nil {
			log.Printf("Ошибка получения списка арендаторов: %v", err)
			return err
		}
		newText = "📋 Список арендаторов"
		state := states.GetListState(query.From.ID) // Инициализируем состояние списка, если оно не было создано
		state.Scope = "tenants"
		state.Page = 0
		newMarkup = menu.AdminTenantsList(tenants, state)

	case query.Data == "add_tenant":
		newText = "Введите username нового арендатора:"
		states.Set(query.From.ID, "awaiting_tenant_data")
		newMarkup = menu.BackButton("list_tenants")

	case strings.HasPrefix(query.Data, "select_activity_type:"):
		idStr := strings.TrimPrefix(query.Data, "select_activity_type:")
		selectedID, _ := strconv.Atoi(idStr)

		selectedIDs := toggleActivitySelection(query.From.ID, selectedID)

		activityTypes, err := db.GetAllActivityTypes()
		if err != nil {
			log.Printf("Ошибка получения видов деятельности: %v", err)
			return err
		}

		// Преобразуем в map[int]bool
		selectedMap := make(map[int]bool)
		for _, id := range selectedIDs {
			selectedMap[id] = true
		}
		newText = "Выберите виды деятельности арендатора:"
		newMarkup = menu.AdminActivityTypeSelect(activityTypes, selectedMap)

	case query.Data == "finish_activity_selection":
		rawIDs, ok := states.GetTemp(query.From.ID, "selected_activity_ids")
		tenantIDStr, ok2 := states.GetTemp(query.From.ID, "tenant_id")
		if !ok || !ok2 || rawIDs == "" {
			newText = "❌ Ошибка: нет выбранных данных."
			newMarkup = menu.BackButton("list_tenants")
			break
		}

		tenantID, _ := strconv.Atoi(tenantIDStr)
		var activityIDs []int
		for _, idStr := range strings.Split(rawIDs, ",") {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				activityIDs = append(activityIDs, id)
			}
		}

		err := db.SaveTenantActivityTypes(tenantID, activityIDs)
		if err != nil {
			log.Printf("Ошибка сохранения видов деятельности арендатора: %v", err)
			newText = "❌ Ошибка при сохранении видов деятельности."
			newMarkup = menu.BackButton("list_tenants")
		} else {
			newText = "✅ Виды деятельности успешно сохранены."
			newMarkup = menu.OkButton("list_tenants")
		}

		states.ClearTemp(query.From.ID)
		states.Set(query.From.ID, "main_menu")

	case query.Data == "page_next":
		states.UpdateListState(query.From.ID, func(s *states.ListState) {
			s.Page++
		})
		state := states.GetListState(query.From.ID)
		switch state.Scope {
		case "users":
			newText, newMarkup, err = buildUserList(query.From.ID)
		case "pavilions":
			pavs, err := db.GetAllPavilions()
			if err != nil {
				return err
			}
			newText = "🏪 Павильоны"
			newMarkup = menu.AdminPavilionList(pavs, state)
		case "activity_types":
			activityTypes, err := db.GetAllActivityTypes()
			if err != nil {
				log.Printf("Ошибка получения видов деятельности: %v", err)
				return err
			}
			newText = "🧑‍🌾 Виды деятельности"
			newMarkup = menu.AdminActivityTypesList(activityTypes, state)
		case "tenants":
			tenants, err := db.GetAllTenants()
			if err != nil {
				log.Printf("Ошибка получения списка арендаторов: %v", err)
				return err
			}
			newText = "📋 Список арендаторов"
			newMarkup = menu.AdminTenantsList(tenants, state)
		}

	case query.Data == "page_prev":
		states.UpdateListState(query.From.ID, func(s *states.ListState) {
			if s.Page > 0 {
				s.Page--
			}
		})
		state := states.GetListState(query.From.ID)
		switch state.Scope {
		case "users":
			newText, newMarkup, err = buildUserList(query.From.ID)
		case "pavilions":
			pavs, err := db.GetAllPavilions()
			if err != nil {
				return err
			}
			newText = "🏪 Павильоны"
			newMarkup = menu.AdminPavilionList(pavs, state)
		case "activity_types":
			activityTypes, err := db.GetAllActivityTypes()
			if err != nil {
				log.Printf("Ошибка получения видов деятельности: %v", err)
				return err
			}
			newText = "🧑‍🌾 Виды деятельности"
			newMarkup = menu.AdminActivityTypesList(activityTypes, state)
		case "tenants":
			tenants, err := db.GetAllTenants()
			if err != nil {
				log.Printf("Ошибка получения списка арендаторов: %v", err)
				return err
			}
			newText = "📋 Список арендаторов"
			newMarkup = menu.AdminTenantsList(tenants, state)
		}

	case strings.HasPrefix(query.Data, "view_pavilion:"):
		id := strings.Split(query.Data, ":")
		if len(id) < 2 {
			log.Printf("Некорректный формат данных для просмотра павильона: %s", query.Data)
			return nil
		}
		pavilionID, err := strconv.Atoi(id[1])
		if err != nil {
			log.Printf("Ошибка преобразования ID павильона %s: %v", id[1], err)
			return err
		}
		pavilion, err := db.GetPavilionByID(pavilionID)
		if err != nil {
			log.Printf("Ошибка получения павильона с ID %d: %v", pavilionID, err)
			return err
		}
		if pavilion == nil {
			log.Printf("Павильон с ID %d не найден", pavilionID)
			return nil
		}
		newText = fmt.Sprintf("🏪 Павильон №%s\nПлощадь: %f\n", pavilion.Number, pavilion.Area)
		newMarkup = menu.BackButton("pavilions")

	case query.Data == "activity_types":
		activityTypes, err := db.GetAllActivityTypes()
		if err != nil {
			log.Printf("Ошибка получения видов деятельности: %v", err)
			return err
		}
		newText = "🧑‍🌾 Виды деятельности"
		state := states.GetListState(query.From.ID) // Инициализируем состояние списка
		state.Scope = "activity_types"
		state.Page = 0
		newMarkup = menu.AdminActivityTypesList(activityTypes, state)

	case query.Data == "add_activity_type":
		newText = "Введите название нового вида деятельности:"
		states.Set(query.From.ID, "awaiting_activity_type_data")
		newMarkup = menu.BackButton("activity_types")

	case strings.HasPrefix(query.Data, "filter:"):
		filter := strings.TrimPrefix(query.Data, "filter:")
		states.UpdateListState(query.From.ID, func(s *states.ListState) {
			s.Filter = filter
			s.Page = 0
		})
		newText, newMarkup, err = buildUserList(query.From.ID)
		if err != nil {
			return err
		}

	case query.Data == "search_user":
		states.Set(query.From.ID, "awaiting_user_search_input")
		newText = "🔍 Введите username или Telegram ID для поиска:"
		newMarkup = menu.BackButton("list_users")

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

	case query.Data == "reset_search":
		states.UpdateListState(query.From.ID, func(s *states.ListState) {
			s.Search = ""
			s.Page = 0
		})
		newText, newMarkup, err = buildUserList(query.From.ID)
		if err != nil {
			return err
		}

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
			log.Printf("%s", msg)
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
	//////////////////////// Добавление администратора //////////////////////////////
	if state == "awaiting_admin_data" {
		username := msg.Text
		err := db.SetUserRole(0, username, "admin")
		if err != nil {
			log.Printf("Ошибка при добавлении админа: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении админа."))

		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Админ добавлен!"))
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
		}
	}
	//////////////////////////////// Добавление арендатора //////////////////////////////
	if state == "awaiting_tenant_data" {
		username := strings.TrimSpace(msg.Text)
		if username == "" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Username не может быть пустым."))
			return nil
		}
		states.SetTemp(userID, "tenant_username", username)
		states.Set(userID, "awaiting_tenant_fio")
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Введите ФИО арендатора:"))
		return nil
	}

	if state == "awaiting_tenant_fio" {
		fio := strings.TrimSpace(msg.Text)
		if fio == "" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ ФИО не может быть пустым."))
			return nil
		}
		// Получаем username из временного хранилища
		states.SetTemp(userID, "tenant_fio", fio)
		states.Set(userID, "awaiting_registration_type")
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Введите номер регистрации арендатора:"))
		return nil
	}
	if state == "awaiting_registration_type" {
		registrationType := strings.TrimSpace(msg.Text)
		if registrationType == "" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Номер регистрации не может быть пустым."))
			return nil
		}
		states.SetTemp(userID, "tenant_registration_type", registrationType)
		states.Set(userID, "awaiting_cash_register")
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "У арендатора есть кассовый аппарат? (да/нет)"))
		return nil
	}
	if state == "awaiting_cash_register" {
		hasCashRegister := strings.ToLower(strings.TrimSpace(msg.Text))
		if hasCashRegister != "да" && hasCashRegister != "нет" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ответ должен быть 'да' или 'нет'."))
			return nil
		}
		// Получаем данные из временного хранилища
		username, ok := states.GetTemp(userID, "tenant_username")
		if !ok {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Не удалось найти username арендатора. Попробуйте заново."))
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		}
		fio, ok := states.GetTemp(userID, "tenant_fio")
		if !ok {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Не удалось найти ФИО арендатора. Попробуйте заново."))
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		}
		registrationType, ok := states.GetTemp(userID, "tenant_registration_type")
		if !ok {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Не удалось найти номер регистрации арендатора. Попробуйте заново."))
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		}
		var hasCashRegisterBool bool

		if hasCashRegister == "да" {
			hasCashRegisterBool = true
		} else {
			hasCashRegisterBool = false
		}

		// Сохраняем арендатора в БД
		tenantId, err := db.AddTenant(username, fio, registrationType, hasCashRegisterBool)
		if err != nil {
			log.Printf("Ошибка при добавлении арендатора: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении арендатора."))

			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		} else {
			states.SetTemp(userID, "tenant_id", strconv.Itoa(int(tenantId)))
			states.Set(userID, "awaiting_tenant_extended_data")
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Арендатор успешно добавлен! Хотите добавить дополнительные данные? (да/нет)"))
		}
		return nil
	}

	if state == "awaiting_tenant_extended_data" {
		answer := strings.ToLower(strings.TrimSpace(msg.Text))
		if answer != "да" && answer != "нет" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ответ должен быть 'да' или 'нет'."))
			return nil
		}
		if answer == "да" {
			states.Set(userID, "awaiting_activity_type_select")
		} else {
			// Завершаем процесс добавления арендатора
			states.ClearTemp(userID)
			states.Set(userID, "main_menu")
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Арендатор успешно добавлен!"))
			sendMenu(ctx, msg)
		}
		return nil
	}

	if state == "awaiting_activity_type_select" {
		ActivityTypes, err := db.GetAllActivityTypes()
		if err != nil {
			log.Printf("Ошибка получения видов деятельности: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка получения видов деятельности."))
			states.Set(userID, "main_menu")
			states.ClearTemp(userID)
			sendMenu(ctx, msg)
			return nil
		}
		selectedMap := make(map[int]bool)
		states.SetTemp(userID, "activity_selection_process", "")
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Выберите тип деятельности арендатора из существующих, или добавьте новый (можно несколько вариантов)").WithReplyMarkup(menu.AdminActivityTypeSelect(ActivityTypes, selectedMap)))
	}
	///////////////////////////////////////// Добавление павильона //////////////////////////////
	if state == "adding_pavilion_number" {
		number := strings.TrimSpace(msg.Text)
		states.SetTemp(userID, "pavilion_number", number)
		states.Set(userID, "adding_pavilion_area")

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Введите площадь павильона в м² (пример 15.5):"))
		return nil
	}
	if state == "adding_pavilion_area" {
		input := strings.TrimSpace(msg.Text)
		area, err := strconv.ParseFloat(input, 64) // Проверка на корректность числа
		if err != nil {
			log.Printf("Ошибка преобразования площади: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Площадь указана некорректно, повторите ввод."))
			return nil
		}
		// Получаем номер из временного хранилища
		number, ok := states.GetTemp(userID, "pavilion_number")
		if !ok {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Не удалось найти номер. Попробуйте заново."))
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		}

		// Сохраняем в БД
		err = db.AddPavilion(number, area)
		if err != nil {
			log.Printf("Ошибка при добавлении павильона: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении павильона."))
		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Павильон успешно добавлен!"))
		}

		states.ClearTemp(userID)
		states.Set(userID, "main_menu")
		sendMenu(ctx, msg)
		return nil
	}
	////////////////////////////////////////////// Добавление вида деятельности //////////////////////////////
	if state == "awaiting_activity_type_data" {
		name := strings.TrimSpace(msg.Text)
		if name == "" {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Название не может быть пустым."))
			return nil
		}
		err := db.AddActivityType(name)
		if err != nil {
			log.Printf("Ошибка при добавлении вида деятельности: %v", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "❌ Ошибка при добавлении вида деятельности."))
		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "✅ Вид деятельности успешно добавлен!"))
		}
		_, ok := states.GetTemp(userID, "activity_selection_process")
		if ok {
			states.Set(userID, "awaiting_activity_type_select")
			return nil
		} else {
			states.Set(userID, "main_menu")
			sendMenu(ctx, msg)
			return nil
		}
	}
	///////////////////// Поиск пользователя //////////////////////////////
	if state == "awaiting_user_search_input" {
		search := strings.TrimSpace(msg.Text)
		states.UpdateListState(userID, func(s *states.ListState) {
			s.Search = search
			s.Page = 0
		})
		states.Set(userID, "main_menu")

		text, markup, err := buildUserList(userID)
		if err != nil {
			fmt.Printf("Ошибка при поиске пользователя: %v\n", err)
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), "Ошибка при поиске пользователя."))
			sendMenu(ctx, msg)
		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(msg.Chat.ID), text).WithReplyMarkup(markup))
		}
		return nil
	}
	if state == "main_menu" {
		sendMenu(ctx, msg)
	}

	return nil
}
