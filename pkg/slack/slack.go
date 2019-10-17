package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ms-choudhary/slackup/pkg/api"
)

type SlackMessage struct {
	User     string `json:"user"`
	Text     string `json:"text"`
	Ts       string `json:"ts"`
	ThreadTs string `json:"thread_ts"`
}

type SlackChannelsHistoryResponse struct {
	Ok       bool           `json:"ok"`
	Messages []SlackMessage `json:"messages"`
	Error    string         `json:"error"`
}

const SlackChannelHistoryApi = "https://slack.com/api/channels.history"

func getSlackMessages(filters map[string]string) ([]SlackMessage, error) {
	queryParams := "?"
	for k, v := range filters {
		queryParams = fmt.Sprintf("%s%s=%s&", queryParams, k, v)
	}
	//remove trailing &
	queryParams = queryParams[:len(queryParams)-1]

	log.Printf("slack history api: %s", SlackChannelHistoryApi+queryParams)

	resp, err := http.Get(SlackChannelHistoryApi + queryParams)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var slackResponse SlackChannelsHistoryResponse
	err = json.Unmarshal(body, &slackResponse)
	if err != nil {
		return nil, err
	}

	if !slackResponse.Ok {
		return nil, fmt.Errorf("slack channel history api failed: %v", slackResponse.Error)
	}

	return slackResponse.Messages, nil
}

func (s SlackMessage) NewMessage() *api.Message {
	return &api.Message{
		User: s.User,
		Text: s.Text,
		Ts:   s.Ts,
	}
}

// parent thread has same ts as thread_ts or empty thread_ts (in case of no threads)
// for comment thread_ts is it's parent's thread_ts
func (s SlackMessage) isParentThread() bool {
	return s.ThreadTs == "" || s.ThreadTs == s.Ts
}

func convertMessages(slackMessages []SlackMessage) []*api.Message {
	res := []*api.Message{}
	lookup := map[string]*api.Message{}
	for _, slackmsg := range slackMessages {
		if slackmsg.isParentThread() {
			if msg, ok := lookup[slackmsg.ThreadTs]; ok {
				msg.UpdateMessage(slackmsg.User, slackmsg.Text, slackmsg.Ts)
			} else {
				lookup[slackmsg.ThreadTs] = slackmsg.NewMessage()
				res = append(res, lookup[slackmsg.ThreadTs])
			}
		} else {
			if _, ok := lookup[slackmsg.ThreadTs]; !ok {
				lookup[slackmsg.ThreadTs] = &api.Message{}
				res = append(res, lookup[slackmsg.ThreadTs])
			}
			lookup[slackmsg.ThreadTs].AddComment(*(slackmsg.NewMessage()))
		}
	}
	return res
}

func debugMessage(m *api.Message) {
	log.Printf("user: %s, text: %s, ts: %s", m.User, m.Text, m.Ts)
	for _, c := range m.Comments {
		debugMessage(&c)
	}
}
