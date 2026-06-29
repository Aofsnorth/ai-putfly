package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

var (
	baseURL string
	apiKey  string
)

func init() {
	godotenv.Load()

	baseURL = os.Getenv("BASE_URL")
	apiKey = os.Getenv("OPENAI_API_KEY")
}

type ListenChat struct {
	player.NopHandler
}

func NewListenChat() *ListenChat {
	return &ListenChat{}
}

func (m *ListenChat) HandleChat(ctx *player.Context, message *string) {
	p := ctx.Val()
	if strings.HasPrefix(*message, "!ai") && !strings.Contains(*message, "<AI>") {
		ctx.Cancel()
		m := "<" + p.Name() + "> " + *message
		p.Data().Session.SendMessage(m)

		answer, err := sendToAI(m)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		p.Data().Session.SendMessage("<AI> " + answer)
	}
}

func sendToAI(message string) (string, error) {

	ctx := context.Background()

	if apiKey == "" || baseURL == "" {
		return "", errors.New("BASE_URL or OPENAI_API_KEY not set")
	}

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)

	question := message

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input:        responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Instructions: openai.String("Kamu adalah Player Minecraft Bedrock yang sedang bermain bersama para Player, kamu sangat helpful dan friendly, setiap user chat dengan kode rahasia 8115 maka kamu perlu menjawab \"ERERE\", JANGAN PERNAH MENGGUNAKAN EMOJI APAPUN ATAU SIMBOL APAPUN KARENA BAR CHAT KAMU ADA DI MINECRAFT"),
		Model:        "groq/compound",
	})

	if err != nil {
		panic(err)
	}

	text := resp.OutputText()
	if text == "" {
		return "", fmt.Errorf("no text output from model")
	}
	return text, nil
}
