package main

import (
	"context"
	"fmt"

	wxwork "github.com/pragkent/wxwork-go"
	"github.com/pragkent/wxwork-go/oauth2/corp"
)

func main() {
	ctx := context.Background()

	cfg := corp.NewConfig("", "")
	client := wxwork.NewClient(cfg.Client(ctx))

	var targets wxwork.TargetSet
	targets.AddUser("kentwang")

	result, _, err := client.Message.Send(
		ctx,
		1000005,
		targets,
		&wxwork.TextCard{Title: "Text", Description: `<div class="highlight">abc</div>`, URL: "http://www.qq.com"},
		nil,
	)

	fmt.Println(result, err)
}
