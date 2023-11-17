package functions

import (
	"context"
	"log"
	"log/slog"

	api "github.com/aokabi/narou-update-notify/api"

	"cloud.google.com/go/firestore"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func NotifyPubSub(ctx context.Context, _ PubSubMessage) error {
	// 確認済みの最新話を取得
	client, err := firestore.NewClient(
		ctx,
		"main-349812",
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	doc, err := client.Collection("state").Doc("latest").Get(ctx)
	if err != nil {
		slog.Error("failed to get document", err)
		return err
	}
	noAny, err := doc.DataAt("no")
	if err != nil {
		slog.Error("failed to get number", err)
		return err
	}
	no, ok := noAny.(int64)
	if !ok {
		slog.Error("failed to convert number", err, "no", noAny)
		return err
	}

	// 最新話を取得
	novelInfo, err := api.GetNovelInfo(ctx)
	if err != nil {
		slog.Error("failed to get novel info", err)
		return err
	}

	latestNo := novelInfo[1].GeneralAllNo

	// 更新があれば、Slackに通知
	if no != int64(latestNo) {
		slog.Info("novel updated", "latest", latestNo)
		//TODO: 通知
	}

	// 確認済みの最新話を更新
	_, err = client.Collection("state").Doc("latest").Update(ctx, []firestore.Update{
		{
			Path:  "no",
			Value: latestNo,
		},
	})
	if err != nil {
		slog.Error("failed to update latest", err)
		return err
	}

	return nil
}
