package main

import (
	"log"
	"math/rand"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Массив мотивационных фраз
var motivationPhrases = []string{
	"Ты можешь больше, чем думаешь!",
	"Каждый шаг приближает тебя к цели!",
	"Никогда не сдавайся, даже если трудно!",
	"Ты сильнее, чем кажется!",
	"Верь в себя — ты уже на правильном пути!",
	"Ты заслуживаешь успеха!",
	"Сегодня лучший день для того, чтобы начать!",
	"Твоя упорность и труд принесут плоды!",
	"Не останавливайся — твоя цель уже близка!",
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
		// Приветствие пользователю
		c.Send("Привет, " + username + "! Я буду отправлять тебе мотивационные фразы раз в минуту.")
		
		// Запускаем фоновую горутину, которая будет отправлять фразы каждую минуту
		go sendRandomMotivationalPhrases(b, c.Sender().ID)

		return nil
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для отправки случайной мотивационной фразы раз в минуту
func sendRandomMotivationalPhrases(b *tele.Bot, userID int64) {
	for {
		// Спим 1 минуту
		time.Sleep(time.Minute)

		// Генерируем случайную фразу
		randomPhrase := motivationPhrases[rand.Intn(len(motivationPhrases))]

		// Отправляем фразу пользователю
		_, err := b.Send(&tele.User{ID: userID}, randomPhrase)
		if err != nil {
			log.Printf("Ошибка при отправке фразы: %v", err)
		}
	}
}
