package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Конфигурация OpenAI API
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"
const openAIAPIKey = "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"

var userCityMap = map[string]string{
	"Москва":   "Europe/Moscow",
	"Лондон":   "Europe/London",
	"Нью-Йорк": "America/New_York",
	"Токио":    "Asia/Tokyo",
	"Сидней":   "Australia/Sydney",
	// Добавьте другие города и часовые пояса
}

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
		username := c.Sender().Username
		if username == "" {
			username = c.Sender().FirstName
		}
		return c.Send("Привет, " + username + "! Напиши мне свой город, чтобы я мог учесть твой часовой пояс.")
	})

	// Запрос имени и города
	b.Handle(tele.OnText, func(c tele.Context) error {
		userMessage := c.Text()
		if _, exists := userCityMap[strings.Title(userMessage)]; exists {
			// Город распознан, сохраняем и сообщаем
			c.Set("city", userMessage)
			return c.Send("Спасибо, я запомнил твой город! Теперь я буду отправлять тебе мотивационные фразы по твоему местному времени.")
		} else if c.Get("city") == nil {
			// Если город не был установлен, просим ввести город
			return c.Send("Пожалуйста, напиши свой город проживания.")
		}

		// Мотивационные фразы в зависимости от времени
		city := c.Get("city").(string)
		motivationPhrase := getMotivationalPhrase(city)
		return c.Send(motivationPhrase)
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для получения мотивационной фразы в зависимости от времени
func getMotivationalPhrase(city string) string {
	// Получаем часовой пояс для города
	location, err := time.LoadLocation(userCityMap[strings.Title(city)])
	if err != nil {
		log.Printf("Ошибка загрузки часового пояса: %v", err)
		return "Не удалось определить часовой пояс для вашего города."
	}

	// Текущее время в выбранном часовом поясе
	currentTime := time.Now().In(location)

	// Определяем фразу в зависимости от времени
	var phrase string
	switch {
	case currentTime.Hour() == 9:
		phrase = "Доброе утро! Сегодня отличный день для новых начинаний!"
	case currentTime.Hour() == 13:
		phrase = "Уже полдень! Время для новых достижений!"
	case currentTime.Hour() == 21:
		phrase = "Вечер наступил. Расслабься и отпразднуй все, что ты достиг сегодня!"
	case currentTime.Hour() == 23 && currentTime.Minute() == 59:
		phrase = "Перед сном подытожь день. Ты сделал много! Завтра будет еще лучше!"
	default:
		phrase = "Продолжай двигаться вперед! Ты на верном пути!"
	}

	return phrase
}

// Функция для отправки запроса к OpenAI (можно оставить, если потребуется)
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
