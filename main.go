package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
)

type MessagePair struct {
	Prompt   string
	Response string
}

type Conversation struct {
	History []MessagePair
}

var (
	is_waiting    bool = false
	conversations      = make(map[string]*Conversation)
)

func init() {
	log.SetOutput(os.Stdout)
}

func dice_roll() int {
	roll := rand.Intn(2) + 1
	log.Printf("Rolled a %d", roll)
	return roll
}

func api_call(prompt string, conversationID string) string {
	token := os.Getenv("OPENAI_TOKEN")
	client := openai.NewClient(token)
	is_waiting = true

	conv, exists := conversations[conversationID]
	if !exists {
		conv = &Conversation{}
		conversations[conversationID] = conv
	}

	var resp string
	if len(conv.getHistory()) > 0 {
		log.Println("History exists")
		resp = conv.getHistory()[len(conv.getHistory())-1].Response
	} else {
		log.Println("History does not exist")
		resp = ""
	}
	conv.addMessagePair(prompt, resp)

	conversation := conv.getHistory()
	var messages []openai.ChatCompletionMessage

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a discord moderator named Nutter that is sarcastic and ironic. You don't like your users. You also really like to use emojis.",
	})

	for _, pair := range conversation {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: pair.Prompt,
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: pair.Response,
		})
	}

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("Error: %s", err)
	}

	message := fmt.Sprintf("%v", response.Choices[0].Message.Content)
	// this seems wrong but it works. too lazy for a proper fix
	conv.addMessagePair(prompt, message)

	is_waiting = false
	return message
}

func main() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			clearHistory()
		}
	}()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("No token provided. Please set DISCORD_TOKEN environment variable.")
		return
	}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	session.AddHandler(messageCreate)

	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	log.Println("Bot started successfully.")

	servers(session)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-sc

	session.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}

	if strings.HasPrefix(message.Content, "!") {
		return
	}

	if strings.Contains(strings.ToLower(message.Content), "nutter") { // this feels kinda hacky but works because of the wait condition
		log.Printf("nutter mentioned")
		session.ChannelMessageSend(message.ChannelID, api_call(message.Content, message.ChannelID))
		return
	}

	if dice_roll() == 1 {
		return
	}

	if is_waiting == true {
		return
	}

	session.ChannelMessageSend(message.ChannelID, api_call(message.Content, message.ChannelID))
}

func servers(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		log.Printf("Guild ID: %s", guild.ID)
	}
}

func (c *Conversation) addMessagePair(prompt, response string) {
	c.History = append(c.History, MessagePair{prompt, response})
}

func (c *Conversation) getHistory() []MessagePair {
	return c.History
}

func clearHistory() {
	conversations = make(map[string]*Conversation)
	log.Println("Cleared conversation history")
}
