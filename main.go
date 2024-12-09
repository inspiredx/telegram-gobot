package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v4"
	"github.com/kelvins/slugify"
)

// Конфигурация OpenAI API
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"
const openAIAPIKey = "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"

var userInfo = make(map[string]map[string]string)

func main() {
	// Токен Telegram бота
	token := "7824776293:AAHIfprFjFTWYBFA05KaHs6cRPN-_xOoe1Q"

	// Настройки бота
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// Создание бота
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %s", err)
	}

	// Логирование всех входящих обновлений
	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			log.Printf("Получено сообщение от %s: %s", c.Sender().Username, c.Text())
			return next(c)
		}
	})

	// Обработка команды /start
	b.Handle("/start", func(c tele.Context) error {
		// Запрос имени пользователя
		return c.Send("Привет! Как тебя зовут?")
	})

	// Обработка текстовых сообщений (имя)
	b.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Sender().Username

		// Если имя еще не введено, сохраняем имя
		if _, exists := userInfo[userID]; !exists {
			userInfo[userID] = make(map[string]string)
			userInfo[userID]["name"] = c.Text()

			// Запрашиваем город
			return c.Send("Приятно познакомиться, "+c.Text()+"! Теперь скажи, в каком городе ты живешь?")
		}

		// Если город еще не введен, сохраняем город
		if userInfo[userID]["city"] == "" {
			userInfo[userID]["city"] = c.Text()

			// Сохраняем время и запускаем рассылку мотивационных фраз
			go startMotivationScheduler(userID)

			return c.Send("Спасибо! Мы будем отправлять тебе мотивационные фразы по твоему времени в течение дня.")
		}

		// Если все уже введено, просто отвечаем
		return c.Send("Данные уже сохранены. Введенные данные: " + "Имя: " + userInfo[userID]["name"] + ", Город: " + userInfo[userID]["city"])
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Мотивационные фразы для разных времен дня
var motivationPhrases = map[int]string{
	9:   "Доброе утро! Сегодня будет замечательный день, не забывай о своих целях!",
	13:  "Полдень! Не останавливайся, продолжай двигаться к своей мечте!",
	21:  "Вечер! Ты молодец, все что было сделано сегодня — это шаг вперед!",
	23:  "Ночь. Помни, что каждый день — это шанс стать лучше! Завтра новый день для новых достижений.",
}

// Функция для рассылки мотивационных фраз по времени пользователя
func startMotivationScheduler(userID string) {
	// Пример часового пояса для города пользователя, можно интегрировать с API для получения точного часового пояса
	// Для простоты используем UTC+3 (можно заменить на реальное значение)
	location, err := time.LoadLocation("Europe/Moscow") // Используем Москву как пример
	if err != nil {
		log.Printf("Ошибка при загрузке часового пояса: %s", err)
		return
	}

	for _, hour := range []int{9, 13, 21, 23} {
		// Получаем время следующего мотивационного сообщения
		nextTime := time.Now().In(location).Truncate(24 * time.Hour).Add(time.Duration(hour) * time.Hour)

		// Если текущее время позже заданного (например, 9:00), ставим на следующий день
		if time.Now().After(nextTime) {
			nextTime = nextTime.Add(24 * time.Hour)
		}

		// Ожидаем до нужного времени и отправляем мотивационную фразу
		time.Sleep(time.Until(nextTime))

		// Отправляем мотивационное сообщение
		sendMotivationMessage(userID, motivationPhrases[hour])
	}
}

// Функция для отправки мотивационного сообщения пользователю
func sendMotivationMessage(userID, message string) {
	// Создаем Telegram бота для отправки сообщений
	token := "7824776293:AAHIfprFjFTWYBFA05KaHs6cRPN-_xOoe1Q"
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Printf("Ошибка при создании бота: %s", err)
		return
	}

	// Отправка сообщения
	user, err := b.ChatByID(userID)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения пользователю %s: %s", userID, err)
		return
	}

	err = b.Send(user, message)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %s", err)
	}
}

// Функция для отправки запроса к OpenAI (если нужно для других целей)
func getOpenAIResponse(userMessage string) (string, error) {
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": userMessage},
		},
	}
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIAPIURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+openAIAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("OpenAI API error: %s", body)
		return "", err
	}

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return "", err
	}

	choices := responseBody["choices"].([]interface{})
	if len(choices) > 0 {
		message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		return message, nil
	}

	return "Пустой ответ от OpenAI.", nil
}
