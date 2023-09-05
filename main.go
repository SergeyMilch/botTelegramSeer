package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var userStates = make(map[int]UserState)

type UserState struct {
	PageNumber int
}

func setupTelegramBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func getInput(prompt string, min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)
	var value int
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		num, err := strconv.Atoi(input)
		if err == nil && num >= min && num <= max {
			value = num
			break
		} else {
			err = fmt.Errorf("Неправильный ввод [%v]. Пожалуйста, попробуйте снова.", num)
		}
	}
	return value
}

func extractSentence(pageContent []string, pageNumber int, lineNumber int) (string, error) {
	if pageNumber < 1 || pageNumber > len(pageContent) {
		return "", fmt.Errorf("Неверный номер страницы: [%v]", pageNumber)
	}

	lines := strings.Split(pageContent[pageNumber-1], "\n")

	if lineNumber < 1 || lineNumber > len(lines) {
		return "", fmt.Errorf("Неверный номер строки: [%v]", lineNumber)
	}

	line := lines[lineNumber-1]

	// Разделение по точкам.
	sentences := strings.Split(line, ".")

	// Возврат первого предложения, либо ошибки, если предложения не найдены.
	if len(sentences) > 0 {
		return sentences[0], nil
	}
	return "", fmt.Errorf("Предложение не найдено: [%v]", sentences)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Путь к книге
	bookPath := "updated_Azazel.txt"

	content, err := os.ReadFile(bookPath)
	if err != nil {
		log.Fatal(err)
	}

	pageContent := strings.Split(string(content), "===============") // Разделитель страниц

	bot, err := setupTelegramBot()
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// updates := bot.ListenForWebhook("/webhook")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Введите номер страницы:")
			bot.Send(msg)
			continue
		}

		if update.Message.Text == "/cancel" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выход из режима ввода. Введите /start, чтобы начать заново.")
			bot.Send(msg)
			continue
		}

		if _, ok := userStates[update.Message.From.ID]; !ok {
			// Новый пользователь начинает ввод номера страницы
			pageNumberInput, err := strconv.Atoi(update.Message.Text)
			if err != nil || pageNumberInput < 1 || pageNumberInput > len(pageContent) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный номер страницы. Пожалуйста, введите корректный номер страницы или введите /cancel для отмены.")
				bot.Send(msg)
			} else {
				// Номер страницы введен корректно, сохраняем его и переходим к вводу номера строки
				userStates[update.Message.From.ID] = UserState{PageNumber: pageNumberInput}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите номер строки:")
				bot.Send(msg)
			}
		} else {
			// Пользователь уже ввел номер страницы, ожидаем ввод номера строки
			lineNumberInput, err := strconv.Atoi(update.Message.Text)
			if err != nil || lineNumberInput < 1 || lineNumberInput > len(strings.Split(pageContent[userStates[update.Message.From.ID].PageNumber-1], "\n")) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный номер строки. Пожалуйста, введите корректный номер строки или введите /cancel для отмены.")
				bot.Send(msg)
			} else {
				// Номер строки введен корректно, извлекаем предложение и отправляем его
				sentence, _ := extractSentence(pageContent, userStates[update.Message.From.ID].PageNumber, lineNumberInput)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извлеченное предложение: "+sentence)
				bot.Send(msg)
				// Удаляем состояние пользователя после завершения запроса
				delete(userStates, update.Message.From.ID)
			}
		}
	}

	// pageNumber := getInput("Введите номер страницы (от 1 до 55): ", 1, len(pageContent))
	// lineNumber := getInput("Введите номер строки (от 1 до 30): ", 1, len(strings.Split(pageContent[pageNumber-1], "\n")))

	// sentence, _ := extractSentence(pageContent, pageNumber, lineNumber)
	// fmt.Println("Извлеченное предложение:", sentence)
}
