package pprof

import (
	"context"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/rakyll/openai-go/chat"
)

var tokenUsageProfile = pprof.NewProfile("openai.token.usage")

type ClientWithProfile struct {
	chat *chat.Client
	// completion completion.Client
}

func NewClientWithProfile(chat *chat.Client) *ClientWithProfile {
	runtime.SetBlockProfileRate(1)
	return &ClientWithProfile{chat: chat}
}

func (c *ClientWithProfile) CreateCompletion(ctx context.Context, params *chat.CreateCompletionParams) (*chat.CreateCompletionResponse, error) {
	r, err := c.chat.CreateCompletion(ctx, params)
	if err != nil {
		return r, err
	}
	tokenUsageProfile.Add(c.chat, r.Usage.TotalTokens)
	return r, err
}

func (c *ClientWithProfile) Close() {
	tokenUsageProfile.Remove(c.chat)
}

func (c *ClientWithProfile) WriteToFile(profileOutPath string) error {
	out, err := os.Create(profileOutPath)
	if err != nil {
		return err
	}
	if err := tokenUsageProfile.WriteTo(out, 0); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
