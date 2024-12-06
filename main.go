package main

import (
	"log"
	"math/rand"
	"time"

	tele "gopkg.in/telebot.v4"
)

func main() {
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
	sideDishes := []string{"Жареная картошка", "Рис", "Овощи", "Гречка", "Макароны"}
	mainDishes := []string{"Стейк", "Курица", "Рыба", "Говядина", "Паста", "Котлеты", "Ребрышки"}

	// Обработка команды /start
	b.Handle("/start", func(c tele.Context) error {
		// Получаем имя пользователя или юзернейм
		username := c.Sender().Username
		if username == "" {
			// Если юзернейм пустой, используем имя пользователя
			username = c.Sender().FirstName
		}
		return c.Send("Привет, " + username + "!")
	})

	// Обработка команды /food
	b.Handle("/food", func(c tele.Context) error {
		// Генерация случайных индексов для выбора гарнира и основного блюда
		rand.Seed(time.Now().UnixNano())
		sideDish := sideDishes[rand.Intn(len(sideDishes))]
		mainDish := mainDishes[rand.Intn(len(mainDishes))]

		// Отправка результата пользователю
		return c.Send("Гарнир: " + sideDish + ", Основное блюдо: " + mainDish)
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}
