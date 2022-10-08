package model

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/db"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

type TalkType int32

const (
	TalkType_OnlineSession TalkType = iota + 1
	TalkType_RecordingSession
	TalkType_Opening
	TalkType_Ending
	TalkType_Commercial

	dateLayout = "2006-01-02"
)

type Talks []Talk

func (t Talk) convertTalkType(title string, presentationMethod *string) (TalkType, error) {
	switch {
	case presentationMethod == nil:
		switch title {
		case "Opening":
			return TalkType_Opening, nil
		case "休憩":
			return TalkType_Commercial, nil
		case "Closing":
			return TalkType_Ending, nil
		}
	case *presentationMethod == "オンライン登壇":
		return TalkType_OnlineSession, nil
	case *presentationMethod == "事前収録":
		return TalkType_RecordingSession, nil
	}
	return 0, fmt.Errorf("model.convertTalkType not found. title: %s, presentationMethod: %s", title, *presentationMethod)
}

func (t Talk) GetTalkType(title string, presentationMethod *string) (TalkType, error) {
	return t.convertTalkType(title, presentationMethod)
}

func (ts Talks) isStartNextTalkSoon() bool {
	now := nowFunc()
	for _, talk := range ts {
		if now.After(talk.StartAt) {
			diffTime := time.Duration(talk.EndAt.Sub(now))
			if 0 < diffTime && diffTime <= utils.HowManyMinutesUntilNotify {
				return true
			}
		}
	}
	return false
}

func (ts Talks) HasNotify(ctx context.Context, rc *db.RedisClient) (bool, error) {
	// 次のtalkがもうすぐ始まるか判定し,まだ通知が行われていない場合は通知を行う.
	if ts.isStartNextTalkSoon() {
		result := rc.Client.Get(ctx, db.NextTalkNotificationKey)
		if result.Err() != nil {
			return false, result.Err()
		}
		if result != nil {
			return false, nil
		}
	}
	return true, nil
}

func (ts Talks) GetCurrentTalk() (*Talk, error) {
	now := nowFunc()
	for _, talk := range ts {
		if now.After(talk.StartAt) && now.Before(talk.EndAt) {
			return &talk, nil
		}
	}
	return nil, fmt.Errorf("Current talk not found")
}

func (ts Talks) GetNextTalk(currentTalk *Talk) (*Talk, error) {
	for i, talk := range ts {
		if talk.Id == currentTalk.Id {
			if i+1 == len(ts) {
				return nil, fmt.Errorf("This talk is last")
			}
			return &ts[i+1], nil
		}
	}
	return nil, fmt.Errorf("Next talk not found")
}

func (t Talk) GetActualStartAtAndEndAt(conferenceDayDate string, startAt, endAt time.Time) (time.Time, time.Time, error) {
	cDate, err := time.Parse(dateLayout, conferenceDayDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return cDate.Add(time.Duration(startAt.Hour()*int(time.Hour) + startAt.Minute()*int(time.Minute) + startAt.Second()*int(time.Second))),
		cDate.Add(time.Duration(endAt.Hour()*int(time.Hour) + endAt.Minute()*int(time.Minute) + endAt.Second()*int(time.Second))),
		nil
}

type Talk struct {
	Id           int32
	TalkName     string
	TrackId      int32
	TrackName    string
	EventAbbr    string
	SpeakerNames []string
	Type         TalkType
	StartAt      time.Time
	EndAt        time.Time
}
