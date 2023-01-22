package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type Talk struct {
	ID        string
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Room      string
	Backup    bool
}

func (t Talk) IsQuickie() bool {
	return t.EndTime.Sub(t.StartTime) == 15*time.Minute
}

func (t Talk) IsLab() bool {
	return t.EndTime.Sub(t.StartTime) >= 100*time.Minute
}

type OpenFeedback struct {
	Sessions map[string]Session `json:"sessions"`
	Speakers map[string]Speaker `json:"speakers"`
}

type Session struct {
	IDSpeakers []string `json:"speakers"`
	Title      string   `json:"title"`
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`
	TrackTitle string   `json:"trackTitle"`
}

type Speaker struct {
	ID string `json:"id"`
}

// ReadTalks returns the list of talks for a speaker in a map indexed by speaker id.
func ReadTalks(loc *time.Location) (map[string][]Talk, error) {
	data, err := os.ReadFile("data/openfeedback.json")
	if err != nil {
		return nil, err
	}

	var openFeedback OpenFeedback
	if err := json.Unmarshal(data, &openFeedback); err != nil {
		return nil, err
	}

	talks := make(map[string][]Talk)
	for id, session := range openFeedback.Sessions {
		startTime, err := time.ParseInLocation("2006-01-02T15:04:05", session.StartTime, loc)
		if err != nil {
			return nil, err
		}

		endTime, err := time.ParseInLocation("2006-01-02T15:04:05", session.EndTime, loc)
		if err != nil {
			return nil, err
		}

		talk := Talk{
			ID:        id,
			Title:     strings.TrimSpace(session.Title),
			StartTime: startTime,
			EndTime:   endTime,
			Room:      session.TrackTitle,
		}

		for _, idSpeaker := range session.IDSpeakers {
			speakerTalks, ok := talks[idSpeaker]
			if ok {
				speakerTalks = append(speakerTalks, talk)
				continue
			}
			speakerTalks = make([]Talk, 0)
			speakerTalks = append(speakerTalks, talk)
			talks[idSpeaker] = speakerTalks
		}
	}

	return talks, err
}
