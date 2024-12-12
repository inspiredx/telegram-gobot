package main

import (
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Определение структуры для вопроса квиза
type QuizQuestion struct {
	Question string
	Options  []string
	Answer   int // Индекс правильного ответа в массиве options (начиная с 0)
}

// Вопросы квиза с вариантами ответов и правильными ответами
var quizQuestions = []QuizQuestion{
	{
		Question: "Когда был основан город Ломоносов?",
		Options:  []string{"1715", "1765", "1812", "1900"},
		Answer:   1, // правильный ответ — 1765
	},
	{
		Question: "Какое название носил город до 1948 года?",
		Options:  []string{"Ломоносов", "Озерки", "Орехово", "Северная звезда"},
		Answer:   0, // правильный ответ — Ломоносов
	},
	{
		Question: "Какой крупнейший памятник города?",
		Options:  []string{"Памятник Петру I", "Памятник В.И. Ленину", "Памятник Николаю I", "Памятник А.С. Пушкину"},
		Answer:   2, // правильный ответ — Памятник Николаю I
	},
	{
		Question: "Какой реки протекает через город Ломоносов?",
		Options:  []string{"Невы", "Охта", "Славянка", "Нарва"},
		Answer:   2, // правильный ответ — Славянка
	},
	{
		Question: "Как называется главный университет в Ломоносове?",
		Options:  []string{"Ломоносовский университет", "Северный университет", "Ленинградский университет", "Московский университет"},
		Answer:   0, // правильный ответ — Ломоносовский университет
	},
}

type UserSession struct {
	QuestionIndex int
	Score         int
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
		c.Send("Привет, " + username + "! Давай начнем квиз по городу Ломоносов. Я буду задавать вопросы по одному.")

		// Начинаем квиз
		startQuiz(b, c)

		return nil
	})

	// Обработка нажатия на кнопки
	b.Handle(tele.OnCallback, func(c tele.Context) error {
		// Получаем сессии пользователя
		session := getSession(c.Sender().ID)

		// Получаем текущий вопрос
		q := quizQuestions[session.QuestionIndex]

		// Проверяем правильность ответа
		userAnswer := c.Callback().Data
		answerIndex, err := stringToInt(userAnswer)
		if err != nil {
			return err
		}

		var correctAnswer string
		if answerIndex == q.Answer {
			session.Score++
			correctAnswer = "Правильно!"
		} else {
			correctAnswer = "Неправильно!"
		}

		// Отправляем сообщение о правильности ответа
		b.Send(c.Sender(), correctAnswer)

		// Переходим к следующему вопросу или завершаем квиз
		session.QuestionIndex++
		if session.QuestionIndex < len(quizQuestions) {
			askQuestion(b, c, session)
		} else {
			b.Send(c.Sender(), fmt.Sprintf("Квиз завершен! Ваш результат: %d/%d", session.Score, len(quizQuestions)))
		}

		return nil
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для начала квиза
func startQuiz(b *tele.Bot, c tele.Context) {
	// Создаем сессию пользователя
	session := &UserSession{
		QuestionIndex: 0,
		Score:         0,
	}

	// Сохраняем сессию
	setSession(c.Sender().ID, session)

	// Задаем первый вопрос
	askQuestion(b, c, session)
}

// Функция для отправки вопроса
func askQuestion(b *tele.Bot, c tele.Context, session *UserSession) {
	q := quizQuestions[session.QuestionIndex]

	// Формируем inline кнопки с вариантами ответов
	buttons := []tele.InlineButton{}
	for i, option := range q.Options {
		buttons = append(buttons, tele.InlineButton{
			Text: option,
			Data: fmt.Sprintf("%d", i), // Сохраняем индекс ответа в Data
		})
	}

	// Формируем сообщение с вопросом
	msg := q.Question
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			buttons,
		},
	}

	// Отправляем вопрос с кнопками
	b.Send(c.Sender(), msg, keyboard)
}

// Функции для работы с сессиями
var sessions = make(map[int64]*UserSession)

func getSession(userID int64) *UserSession {
	session, exists := sessions[userID]
	if !exists {
		return nil
	}
	return session
}

func setSession(userID int64, session *UserSession) {
	sessions[userID] = session
}

func stringToInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
