package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botToken    = "7784381609:AAGUppghBFCRJKrAzCNRp-SrRzj1ebA2gLI"
	logFileName = "bot_logs.txt"
)

var (
	adminChatID int64 = 6223799614
	groupChatID int64 = -1002591133699
)

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "BOT_LOG: ", log.LstdFlags|log.Lshortfile)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panicf("Ошибка создания бота: %v", err)
	}
	bot.Debug = true
	log.Printf("Бот запущен: @%s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		userName := update.Message.From.UserName
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		log.Printf("ChatID: %d", chatID)

		// --- Ответ на /start с инструкцией ---
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
			intro := `✨ *Привет! Я — бот поддержки Skillsrock Academy.*

		Чтобы мы быстрее и точнее помогли вам, пожалуйста, опишите ваш вопрос в сообщении.
		Это может быть *технический* или *организационный* вопрос — мы всё передадим нужному специалисту 💬
		
		📌 *Пожалуйста, укажите:*
		1. 🔍 *Краткое описание проблемы*
		2. 📚 *Тема или направление*
		3. 🕒 *Примерные сроки решения (если есть дедлайн)*
		
		*Пример:*
		Мне нужна помощь с установкой зависимостей".
		Хотелось бы решить в течение дня.
		
		Мы обязательно свяжемся с вами в ближайшее время! Спасибо за обращение 🙌`
		
				msg := tgbotapi.NewMessage(chatID, intro)
				msg.ParseMode = "Markdown"
				_, _ = bot.Send(msg)
				continue
			}
		}

		// Формируем лог-сообщение
		logText := "👤 Пользователь: @" + userName +
			" (ID: " + strconv.FormatInt(int64(userID), 10) + ")\n" +
			"✉️ Сообщение: " + text + "\n---"

		logger.Println(logText) // Лог в файл

		// Отправка админу
		adminMsg := tgbotapi.NewMessage(adminChatID, logText)
		_, _ = bot.Send(adminMsg)

		// Отправка в группу
		groupMsg := tgbotapi.NewMessage(groupChatID, logText)
		_, err := bot.Send(groupMsg)
		if err != nil {
			log.Printf("Ошибка отправки в группу: %v", err)
		}

		// Ответ пользователю
		reply := "✅ *Ваш запрос принят! В ближайшее время с вами свяжется технический специалист.* \n\n🔹После предоставления поддержки, пожалуйста, напишите в этот чат, удалось ли решить вашу задачу, а также укажите примерное время, которое вы затратили на её решение.💻⏳"
		userMsg := tgbotapi.NewMessage(chatID, reply)
		userMsg.ParseMode = "Markdown"
		_, _ = bot.Send(userMsg)
	}
}
