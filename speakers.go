package main

import (
	"encoding/csv"
	"os"
	"strings"
)

type Sex string

const (
	SexFemale = "F"
	SexMale   = "M"
)

func (s Sex) isFemale() bool {
	return s == SexFemale
}

func (s Sex) isMale() bool {
	return s == SexMale
}

type Information struct {
	Firstname     string
	Lastname      string
	Ticket        string
	Email         string
	UID           string
	Backup        bool
	Accommodation bool
	Sex           Sex
}

func (i Information) IsFemale() bool {
	return i.Sex.isFemale()
}

func (i Information) isMale() bool {
	return i.Sex.isMale()
}

func ReadSpeakers() (map[string]Information, error) {
	r, err := os.Open("data/speakers.csv")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	information := make(map[string]Information, len(records))
	for i, record := range records {
		if i == 0 {
			continue
		}

		info := Information{
			Firstname:     strings.TrimSpace(record[0]),
			Lastname:      strings.TrimSpace(record[1]),
			Ticket:        record[2],
			Email:         strings.TrimSpace(record[3]),
			UID:           record[4],
			Backup:        record[5] == "true",
			Accommodation: record[6] == "true",
			Sex:           Sex(strings.TrimSpace(record[7])),
		}

		if info.UID == "" {
			info.UID = info.Ticket
		}

		information[info.UID] = info
	}

	return information, nil
}
