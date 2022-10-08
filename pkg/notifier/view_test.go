package notifier

import (
	"testing"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/google/go-cmp/cmp"
)

func Test_viewSession(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		expectedStr := `
{
	"blocks": [
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Current Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "10:00 - 11:00",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: オンライン登壇",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: kanata",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10001|ものすごい発表>"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "plain_text",
					"text": "Next Talk",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "plain_text",
					"text": "Track A",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "11:00 - 12:30",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Type: 事前収録",
					"emoji": true
				},
				{
					"type": "plain_text",
					"text": "Speaker: hoge, fuga",
					"emoji": true
				}
			]
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Title: <https://event.cloudnativedays.jp/cndt2101/talks/10002|さらにものすごい発表>"
			},
			"accessory": {
				"type": "multi_static_select",
				"placeholder": {
					"type": "plain_text",
					"text": "switching",
					"emoji": true
				},
				"options": [
					{
						"text": {
							"type": "plain_text",
							"text": "シーンを切り替える",
							"emoji": true
						},
						"value": "1__A"
					}
				],
				"action_id": "broadcast_scenenext"
			}
		},
		{
			"type": "divider"
		}
	]
}
`
		expected, err := castFromStringToMsg(expectedStr)
		if err != nil {
			t.Fatal(err)
		}

		got, err := viewSession(model.CurrentAndNextTalk{
			Current: model.Talk{
				Id:           10001,
				TalkName:     "ものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 10, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_OnlineSession,
				SpeakerNames: []string{"kanata"},
				EventAbbr:    "cndt2101",
			},
			Next: model.Talk{
				Id:           10002,
				TalkName:     "さらにものすごい発表",
				TrackId:      1,
				TrackName:    "A",
				StartAt:      time.Date(2022, 10, 1, 11, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				EndAt:        time.Date(2022, 10, 1, 12, 30, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
				Type:         model.TalkType_RecordingSession,
				SpeakerNames: []string{"hoge", "fuga"},
				EventAbbr:    "cndt2101",
			},
		})
		if err != nil {
			t.Errorf("error = %v", err)
			return
		}
		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf(diff)
		}
	})
}
