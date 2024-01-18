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

	// 更新がなければ終了
	if no == int64(latestNo) {
		slog.Info("no update")
		return nil
	}

	// 更新があれば、メールで通知
	subject := "なろう更新通知"
	body := fmt.Sprintf("最新話: %d", latestNo)
	if err := SendEmail(ctx, subject, body); err != nil {
		slog.Error("failed to send email", err)
		return err
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
