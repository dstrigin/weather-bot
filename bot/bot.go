package bot

import (
	"echo-bot/weather"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Init() {
	godotenv.Load()

	botToken := os.Getenv("BOT_API_KEY")
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())

	// checking if token is invalid
	if err != nil {
		log.Fatalf("Can't create bot: %s", err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	msgHandler, _ := th.NewBotHandler(bot, updates)
	msgHandler.Handle(StartHandler, th.CommandEqual("start"))
	msgHandler.Handle(WeatherHandler, th.TextEqual("Узнать погоду"))
	msgHandler.Handle(ResponseHandler, th.AnyMessage())

	defer msgHandler.Stop()
	defer bot.StopLongPolling()

	msgHandler.Start()

}

func StartHandler(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Узнать погоду"),
		),
	)

	_, _ = bot.SendSticker(
		tu.Sticker(
			chatID,
			tu.FileFromID("CAACAgIAAxkBAAEMpIxmu5uZZNjmlqHjTy7uPyyNyHaZ5wACNQEAAjDUnRG0uDX9ZqC2fDUE"),
		),
	)
	message := tu.Message(chatID,
		"Привет! По твоему запросу я смогу подсказать тебе погоду в Ростове-на-Дону на данный момент.\nВоспользуйся опцией ниже.").WithReplyMarkup(keyboard)

	_, _ = bot.SendMessage(message)
}

func WeatherHandler(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)
	message := tu.Message(chatID, weather.GetInfoForBot())
	_, _ = bot.SendMessage(message)
}

func ResponseHandler(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)
	message := tu.Message(chatID,
		"Я не могу ответить на это сообщение. Единственная моя функция - подсказывать погоду.")
	_, _ = bot.SendMessage(message)
}
