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

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
)

var (
	is_waiting bool = false
)

func init() {
	log.SetOutput(os.Stdout)
}

func dice_roll() int {
	roll := rand.Intn(2) + 1
	log.Printf("Rolled a %d", roll)
	return roll
}

func api_call(prompt string) string {
	role := "You are a discord moderator named Nutter that is sarcastic and ironic. You don't like your users. You also really like to use emojis."
	token := os.Getenv("OPENAI_TOKEN")
	client := openai.NewClient(token)
	is_waiting = true
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: role,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion  error: %v\n", err)
		return ""
	}
	message := fmt.Sprintf("%v", response.Choices[0].Message.Content)
	log.Printf("Message: %s", message)
	is_waiting = false
	return message
}

func main() {
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
	if dice_roll() == 1 {
		return
	}
	if is_waiting == true {
		return
	}
	session.ChannelMessageSend(message.ChannelID, api_call(message.Content))
}
