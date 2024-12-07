package main

import (
	"bytes"
	"encoding/json"
	"io" // заменили ioutil на io
	"log"
	"math/rand"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Конфигурация OpenAI API
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"
const openAIAPIKey = "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"

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

	// Массивы для гарниров и основных блюд
	sideDishes := []string{"Картошка", "Рис", "Овощи", "Гречка", "Макароны"}
	mainDishes := []string{"Стейк", "Курица", "Рыба", "Говядина", "Паста"}

	// Обработка команды /start
	b.Handle("/start", func(c tele.Context) error {
		username := c.Sender().Username
		if username == "" {
			username = c.Sender().FirstName
		}
		return c.Send("Привет, " + username + "! Используй /food, чтобы получить идею для обеда.")
	})

	// Обработка команды /food
	b.Handle("/food", func(c tele.Context) error {
		rand.Seed(time.Now().UnixNano())
		sideDish := sideDishes[rand.Intn(len(sideDishes))]
		mainDish := mainDishes[rand.Intn(len(mainDishes))]
		return c.Send("Гарнир: " + sideDish + ", Основное блюдо: " + mainDish)
	})

	// Обработка сообщений, пересылаемых в OpenAI GPT
	b.Handle(tele.OnText, func(c tele.Context) error {
		userMessage := c.Text()
		log.Printf("Запрос в OpenAI: %s", userMessage)
		response, err := getOpenAIResponse(userMessage)
		if err != nil {
			log.Printf("Ошибка при запросе к OpenAI: %s", err)
			return c.Send("Не удалось получить ответ от OpenAI.")
		}
		return c.Send(response)
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для отправки запроса к OpenAI
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
		body, _ := io.ReadAll(resp.Body) // заменили ioutil.ReadAll на io.ReadAll
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
