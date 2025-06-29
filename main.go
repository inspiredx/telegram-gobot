package main

import (
	"log"
	"math/rand"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Массив мотивационных фраз
var motivationPhrases = []string{
	"Ты можешь больше, чем думаешь!!",
	"Каждый шаг приближает тебя к цели!",
	"Никогда не сдавайся, даже если трудно!",
	"Ты сильнее, чем кажется!",
	"Верь в себя — ты уже на правильном пути!",
	"Ты заслуживаешь успеха!",
	"Сегодня лучший день для того, чтобы начать!",
	"Твоя упорность и труд принесут плоды!",
	"Не останавливайся — твоя цель уже близка!",
	"Не бойся ошибаться, ведь ошибки — это опыт!",
	"Каждый день — это шанс стать лучше, чем вчера!",
	"Если ты хочешь, чтобы всё изменилось, начни с себя!",
	"Ты способен на большее, чем можешь себе представить!",
	"Ты уже делаешь шаги к своей мечте — не останавливайся!",
	"Будь настойчивым, и успех не заставит себя ждать!",
	"Не останавливайся на достигнутом — впереди ещё много побед!",
	"Не бойся неудач — они приводят к успеху!",
	"Каждый новый день — это возможность для роста!",
	"Самые лучшие вещи происходят за пределами твоей зоны комфорта!",
	"Ты — автор своей жизни, и она будет успешной!",
	"Смело иди к своей мечте, и она сбудется!",
	"Ты заслуживаешь счастья и успеха!",
	"Если ты хочешь изменений — начни с действий!",
	"Сегодня — идеальный день для начала нового пути!",
	"Ты обладаешь всеми силами, чтобы изменить свою жизнь!",
	"Будь лучшей версией себя и продолжай двигаться вперед!",
	"Верь в себя, даже если мир сомневается!",
	"Ты всегда можешь сделать первый шаг — и всё изменится!",
	"Твои усилия обязательно приведут к результату!",
	"Жизнь — это не проблема, которую нужно решить, а приключение, которое нужно пережить!",
	"Только ты решаешь, как прожить этот день!",
	"Каждый день — это новая возможность улучшить свою жизнь!",
	"Никогда не теряй веру в себя, даже когда всё идет не по плану!",
}

func main() {
	// Токен Telegram бота
	token := ${{ secrets.TELEGRAM_BOT_TOKEN }}

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
		c.Send("Привет, " + username + "! Я буду отправлять тебе мотивационные фразы раз в 10 минут.")
		
		// Запускаем фоновую горутину, которая будет отправлять фразы каждые 10 минут
		go sendRandomMotivationalPhrases(b, c.Sender().ID)

		return nil
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для отправки случайной мотивационной фразы каждые 10 минут
func sendRandomMotivationalPhrases(b *tele.Bot, userID int64) {
	for {
		// Спим 10 минут
		time.Sleep(10 * time.Minute)

		// Генерируем случайную фразу
		randomPhrase := motivationPhrases[rand.Intn(len(motivationPhrases))]

		// Отправляем фразу пользователю
		_, err := b.Send(&tele.User{ID: userID}, randomPhrase)
		if err != nil {
			log.Printf("Ошибка при отправке фразы: %v", err)
		}
	}
}
