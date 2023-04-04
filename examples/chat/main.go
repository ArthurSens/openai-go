package main

import (
	"context"
	"log"
	"os"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
	"github.com/rakyll/openai-go/pprof"
)

func main() {
	ctx := context.Background()
	s := openai.NewSession(os.Getenv("OPENAI_API_KEY"))

	client := pprof.NewClientWithProfile(chat.NewClient(s, "gpt-3.5-turbo"))
	defer client.Close()
	defer func() {
		err := client.WriteToFile("profile.pb.gz")
		if err != nil {
			log.Fatalf("Failed to write profile: %v", err)
		}
	}()
	resp, err := client.CreateCompletion(ctx, &chat.CreateCompletionParams{
		Messages: []*chat.Message{
			{Role: "user", Content: "hello"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to complete: %v", err)
	}

	for _, choice := range resp.Choices {
		msg := choice.Message
		log.Printf("role=%q, content=%q", msg.Role, msg.Content)
	}
}
