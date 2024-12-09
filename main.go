package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Конфигурация OpenAI API
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"
const openAIAPIKey = "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"

// Карта городов и часовых поясов
var userCityMap = map[string]string{
	"Москва":   "Europe/Moscow",
	"Лондон":   "Europe/London",
	"Нью-Йорк": "America/New_York",
	"Токио":    "Asia/Tokyo",
	"Сидней":   "Australia/Sydney",
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

	// Запуск фоновой горутины, которая будет отправлять мотивационные фразы в 9:00, 13:00, 21:00, 23:59, 00:10 и 00:15
	go func() {
		for {
			// Каждую минуту проверяем текущее время
			time.Sleep(time.Minute)
			checkAndSendMotivation(b)
		}
	}()

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
	case currentTime.Hour() == 0 && currentTime.Minute() == 10:
		phrase = "10 минут после полуночи, но не беспокойся — новый день уже начался! Время двигаться вперед!"
	case currentTime.Hour() == 0 && currentTime.Minute() == 15:
		phrase = "15 минут после полуночи! Начни новый день с вдохновения и уверенности!"
	default:
		phrase = "Продолжай двигаться вперед! Ты на верном пути!"
	}

	return phrase
}

// Функция для проверки времени и отправки мотивационной фразы
func checkAndSendMotivation(b *tele.Bot) {
	// Список всех пользователей, которым нужно отправить фразу (предположим, что пользователи сохранены где-то в мапе)
	// Для простоты примера будем использовать фиктивную карту пользователей
	users := []tele.User{
		{Username: "testuser1", ID: 12345}, // Пример пользователя
	}

	// Для каждого пользователя проверяем его город и отправляем мотивационную фразу
	for _, user := range users {
		city := "Москва" // Здесь нужно получить город пользователя, если он уже указан
		motivationPhrase := getMotivationalPhrase(city)
		b.Send(&user, motivationPhrase)
	}
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
