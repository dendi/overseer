//
// This is the slack bridge, which should be built like so:
//
//     go build .
//
// Once built launch it as follows:
//
//     $ ./slack-bridge -slack=slack-webhook-url
//
// When a test fails an slack will sent via SMTP
//
// Eka
// --
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

type SlackRequestBody struct {
	Text string 					`json:"text,omitempty"`
	IconEmoji string 				`json:"icon_emoji,omitempty"`
	Channel string					`json:"channel"`
	Blocks []SlackBlock				`json:"blocks"`
}

type SlackBlock struct {
	Type string			`json:"type"`
	Text SlackText	`json:"text,omitempty"`
}

type SlackText struct {
	Text string			`json:"text"`
	Type string			`json:"type"`
}

type SlackBrige struct {
	slackWebhook string
	slackChannel string

	SendTestSuccess   bool
	SendTestRecovered bool
}

//
// Given a JSON string decode it and post it via slack if it describes
// a test-failure.
//
func (bridge *SlackBridge) process(msg []byte) {
	testResult, err := test.ResultFromJSON(msg)
	if err != nil {
		panic(err)
	}

	// If the test passed then we don't care, unless otherwise defined
	shouldSend := true
	if testResult.Error == nil {
		shouldSend = false

		if bridge.SendTestSuccess {
			shouldSend = true
		}

		if bridge.SendTestRecovered && testResult.Recovered {
			shouldSend = true
		}
	}

	if !shouldSend {
		return
	}

	fmt.Printf("Processing result: %+v\n", testResult)

	// Define Title
	titleText := SlackText{
		Text: fmt.Sprintf(":warning: *%s %s*", "Error:", testResult.Error),
		Type: "mrkdwn",
	}

	if testResult.IsDedup {
		titleText.Text = fmt.Sprintf(":warning: *%s %s*", "Error (deduplicated):", testResult.Error)
	}

	if testResult.Recovered {
		titleText.Text = ":white_check_mark: *Error Recovered*"
	}

	title := SlackBlock{
		Type: "section",
		Text: titleText,
	}

	tagText := SlackText{
		Text: "",
		Type: "mrkdwn",
	}

	// Define Tag
	if testResult.Tag != "" {
		tagText.Text = fmt.Sprintf("Tag : %s", testResult.Tag)
	} else {
		tagText.Text = "Tag : None"
	}

	tag := SlackBlock{
		Type: "context",
		Text: tagText,
	}

	divider := SlackBlock{
		Type: "divider",
	}

	body := SlackRequestBody{
		Channel: slackChannel,
		Blocks: []SlackBlock{
			title,
			tag,
			divider,
		}
	}

	if testResult.Details != "" {
		detail = SlackBlock{
			Type: "section",
			Text: SlackText{
				Text: testResult.Details,
				Type: "mrkdwn",
			}
		}
		body.Blocks = append(body.Blocks, detail)
	}

	info := SlackBlock{
		Type: "section",
		Text: SlackText{
			Type: "mrkdwn",
			Text: fmt.Sprintf("Input: %s\nTarget: %s\nType: %s", testResult.Input, testResult.Target, testResult.Type)
		},
	}
	body.Blocks = append(body.Blocks, info)

	date := SlackBlock{
		Type: "context",
		Text: SlackText{
			Type: "mrkdwn",
			Text: time.Now().UTC().String()
		},
	}
	body.Blocks = append(body.Blocks, date)

	slackBody, _ := json.Marshal(body)
    req, err := http.NewRequest(http.MethodPost, slackWebhook, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
}

//
// Entry Point
//
func main() {

	//
	// Parse our flags
	//
	redisHost := flag.String("redis-host", "127.0.0.1:6379", "Specify the address of the redis queue.")
	redisPass := flag.String("redis-pass", "", "Specify the password of the redis queue.")
	redisQueueKey := flag.String("redis-queue-key", "overseer.results", "Specify the redis queue key to use.")

	slackWebhook := flag.String("slack-webhook", "https://hooks.slack.com/services/T1234/xxxx/xxx", "Slack Webhook URL")
	slackChannel := flag.String("slack-channel", "#my-channel", "Slack Channel Name")

	sendTestSuccess := flag.Bool("send-test-success", false, "Send also test results when successful")
	sendTestRecovered := flag.Bool("send-test-recovered", false, "Send also test results when a test recovers from failure (valid only when used together with deduplication rules)")

	flag.Parse()

	//
	// Create the redis client
	//
	r := redis.NewClient(&redis.Options{
		Addr:     *redisHost,
		Password: *redisPass,
		DB:       0, // use default DB
	})

	//
	// And run a ping, just to make sure it worked.
	//
	_, err := r.Ping().Result()
	if err != nil {
		fmt.Printf("Redis connection failed: %s\n", err.Error())
		os.Exit(1)
	}

	bridge := SlackBridge{
		slackWebhook:      slackWebhook,
		slackChannel:	   slackChannel,
		SendTestRecovered: *sendTestRecovered,
		SendTestSuccess:   *sendTestSuccess,
	}

	for {

		//
		// Get test-results
		//
		msg, _ := r.BLPop(0, *redisQueueKey).Result()

		//
		// If they were non-empty, process them.
		//
		//   msg[0] will be "overseer.results"
		//
		//   msg[1] will be the value removed from the list.
		//
		if len(msg) >= 1 {
			bridge.Process([]byte(msg[1]))
		}
	}
}
