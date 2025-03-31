package functions

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	api "github.com/aokabi/narou-update-notify/api"

	"cloud.google.com/go/firestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type ncode = string

type latestEpisodes struct {
	Episodes map[ncode]episodeInfo `firestore:"episodes"`
}

type episodeInfo struct {
	LatestEpisodeNo int    `firestore:"latest_episode_no"`
	Title           string `firestore:"title"`
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

	doc, err := client.Collection("latestEpisodes").Doc("latest").Get(ctx)
	if err != nil {
		slog.Error("failed to get document", err)
		return err
	}
	var latests latestEpisodes
	if err := doc.DataTo(&latests); err != nil {
		slog.Error("failed to unmarshal document", err)
		return err
	}

	for c, episode := range latests.Episodes {
		novelInfo, err := api.GetNovelInfo(ctx, c)
		if err != nil {
			slog.Error("failed to get novel info", err)
			return err
		}

		latestNo := novelInfo[1].GeneralAllNo
		// 更新がなければ次の作品へ
		if int64(episode.LatestEpisodeNo) == int64(latestNo) {
			slog.Info("no update for ncode:", c)
			continue
		}

		// 更新があれば、メールで通知
		subject := fmt.Sprintf("なろう更新通知(%s)", episode.title)
		body := fmt.Sprintf("最新話: %d", latestNo)
		if err := SendEmail(ctx, subject, body); err != nil {
			slog.Error("failed to send email", err)
			return err
		}

		// 確認済みの最新話を更新
		if _, err := client.Collection("latestEpisodes").Doc("latest").Update(ctx, []firestore.Update{
			{
				Path:  fmt.Sprintf("episodes.%s.latest_episode_no", c),
				Value: latestNo,
			},
		}); err != nil {
			slog.Error("failed to update document", err)
			return err
		}
	}

	return nil
}

func SendEmail(ctx context.Context, subject, body string) error {
	// Create a new session using your AWS credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"), // Replace with your desired AWS region
	})
	if err != nil {
		return err
	}

	// Create a new SES service client
	svc := ses.New(sess)

	// Specify the email parameters
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(os.Getenv("TO_ADDRESS")), // Replace with the recipient's email address
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body), // Replace with the email body content
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject), // Replace with the email subject
			},
		},
		Source: aws.String(os.Getenv("FROM_ADDRESS")), // Replace with the sender's email address
	}

	// Send the email
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}
