package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AlexEagle1535/market-rent-bot/db"
	"github.com/AlexEagle1535/market-rent-bot/handlers"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	_ = godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	// Чтение и преобразование строки в слайс
	var admins = strings.Split(os.Getenv("ADMINS"), ",")
	db.InitDB()
	if db.DB == nil {
		log.Fatal("Ошибка: соединение с БД не установлено")
	}
	defer db.DB.Close()
	for _, admin := range admins {
		admin = strings.TrimSpace(admin) // Удаляем лишние пробелы
		if admin == "" {
			log.Fatal("Пустой username в списке админов")
		}
		err := db.SetUserRole(0, admin, "admin")
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
	bh.HandleMessage(handlers.Start, th.CommandEqual("start"))

	// Обработка callback-кнопок
	bh.HandleCallbackQuery(handlers.CallbackQuery, th.AnyCallbackQueryWithMessage())

	// обработчик тексовых событий
	bh.HandleMessage(handlers.TextMessage)
	// Запускаем обработчик
	go func() {
		if err := bh.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Ждём сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⛔ Завершение работы бота. Отключаем соединения.")
}
