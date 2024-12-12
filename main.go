package main

import (
	"log"
	"math/rand"
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

func main() {
	// Токен Telegram бота
	token := "YOUR_BOT_TOKEN"

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
		c.Send("Привет, " + username + "! Давай начнем квиз по городу Ломоносов. Я буду задавать вопросы, а ты выбирай правильный ответ!")

		// Запуск квиза
		startQuiz(b, c)

		return nil
	})

	// Запуск бота
	log.Println("Бот запущен...")
	b.Start()
}

// Функция для начала квиза
func startQuiz(b *tele.Bot, c tele.Context) {
	score := 0
	totalQuestions := len(quizQuestions)

	// Перемешиваем вопросы
	rand.Seed(time.Now().UnixNano())
	quizQuestionsCopy := append([]QuizQuestion{}, quizQuestions...)
	rand.Shuffle(len(quizQuestionsCopy), func(i, j int) {
		quizQuestionsCopy[i], quizQuestionsCopy[j] = quizQuestionsCopy[j], quizQuestionsCopy[i]
	})

	// Задаем вопросы
	for _, q := range quizQuestionsCopy { // Используем _ вместо i
		// Формируем сообщение с вопросом и вариантами ответов
		options := ""
		for j, option := range q.Options {
			options += string('A'+j) + ". " + option + "\n"
		}
		msg := q.Question + "\n" + options

		// Отправляем вопрос
		b.Send(c.Sender(), msg)

		// Ожидаем ответа пользователя
		answerCh := make(chan string)
		b.Handle("/answer", func(c tele.Context) error {
			answer := c.Text()
			answerCh <- answer
			return nil
		})

		// Получаем ответ от пользователя
		answer := <-answerCh

		// Проверяем, правильный ли ответ
		var correctAnswer string
		if answer == string('A'+q.Answer) {
			score++
			correctAnswer = "Правильно!"
		} else {
			correctAnswer = "Неправильно!"
		}

		// Отправляем ответ пользователю
		b.Send(c.Sender(), correctAnswer)
	}

	// Подведение итогов
	resultMessage := "Квиз завершен! Ваш результат: \n"
	resultMessage += "Правильных ответов: " + string(score) + "/" + string(totalQuestions)
	b.Send(c.Sender(), resultMessage)
}
