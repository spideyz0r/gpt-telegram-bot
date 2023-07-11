package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pborman/getopt"
	openai "github.com/spideyz0r/openai-go"
)

const (
	command_prefix      = "/"
	default_system_role = "You're a Guru that combines your technical expertise, clear communication style, and didactic approach to share your knowledge and answer questions."
)

func main() {

	help := getopt.BoolLong("help", 'h', "display this help")
	model := getopt.StringLong("model", 'm', "gpt-3.5-turbo", "model. default: gpt-3.5-turbo")
	whitelist_path := getopt.StringLong("whitelist file", 'w', "", "path to file with whitelisted users")
	debug := getopt.BoolLong("debug", 'd', "enable debug mode")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	var whitelisted = map[int64]bool{}
	if *whitelist_path != "" {
		readWhiteList(&whitelisted, *whitelist_path)
	}

	var conversations = make(map[int64]chan string)
	var err error

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_KEY"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Creating bot with name %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Starting to listen for updates")

	for update := range updates {
		if *whitelist_path != "" && !whitelisted[update.Message.Chat.ID] {
			fmt.Printf("User %d is not in the whitelist\n", update.Message.Chat.ID)
			continue
		}
		conversation, ok := conversations[update.Message.Chat.ID]
		if !ok {
			conversation = make(chan string)
			conversations[update.Message.Chat.ID] = conversation
			go runConversation(update.Message.Chat.ID, bot, conversation, *debug, *model)
		}
		conversation <- update.Message.Text
	}
}

func runConversation(userID int64, telegramBot *tgbotapi.BotAPI, conversation chan string, debug bool, model string) {
	var openai_client *openai.OpenAIClient
	var botMessage string

	apiKey := os.Getenv("OPENAI_API_KEY")
	openai_client = openai.NewOpenAIClient(apiKey)

	messages := []openai.Message{
		{
			Role:    "system",
			Content: default_system_role,
		},
	}

	for {
		userMessage := <-conversation
		log.Printf("New message from user %i: %s", userID, userMessage)
		if debug {
			log.Printf("Messages m1: %v", messages)
		}

		if strings.HasPrefix(userMessage, command_prefix) {
			re := regexp.MustCompile(`^/[^ ]+`)
			command := re.FindString(userMessage)
			args := re.ReplaceAllString(userMessage, "")
			log.Printf("New command received %s: command: %s args: %s", userID, command, args)
			botMessage, messages = runCommands(command, args, messages)

		} else {
			messages = append(messages, openai.Message{
				Role:    "user",
				Content: string(userMessage),
			})
			completion, err := openai_client.GetCompletion(model, messages, float32(0.5))
			if err != nil {
				log.Printf("Error creating OpenAI completion: %s", err)
				continue
			}
			if debug {
				log.Printf("Completion: %v", completion)
			}
			botMessage = completion.Choices[0].Message.Content
		}

		if debug {
			log.Printf("Messages m2: %v", messages)
		}
		log.Printf("New Answer from bot %s: %s", userID, botMessage)
		msg := tgbotapi.NewMessage(userID, botMessage)
		_, err := telegramBot.Send(msg)
		if err != nil {
			log.Printf("Error sending message to user: %s", err)
		}
	}
}

func runCommands(command string, args string, messages []openai.Message) (string, []openai.Message) {
	var botMsg string
	var commands = map[string]string{
		"/help":  "show this help",
		"/reset": "restart the conversation",
		"/role":  "set the system role",
	}

	switch command {
	case "/role":
		if args == "" {
			botMsg = "Syntax is /role <role>"
			return botMsg, messages
		}
		messages[0] = openai.Message{
			Role:    "system",
			Content: args,
		}
		botMsg = fmt.Sprintf("Role set to %s", strings.TrimSpace(args))
	case "/reset":
		messages = []openai.Message{
			{
				Role:    "system",
				Content: default_system_role,
			},
		}
		botMsg = "Conversation reset"
	default:
		botMsg = "Available commands:\n"
		for k, v := range commands {
			botMsg += fmt.Sprintf("%s: %s\n", k, v)
		}
	}
	return botMsg, messages
}

func readWhiteList(whitelist *map[int64]bool, file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := scanner.Text()
		id_64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		(*whitelist)[id_64] = true
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
