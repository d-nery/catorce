package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type OverallStats map[int]ChatStats

type ChatStats struct {
	Group   GroupStats
	Players map[int]PlayerStats
}

type GroupStats struct {
	GamesPlayed int
	P2Sequence  int
	P4Played    int
}

type PlayerStats struct {
	GamesWon       int
	GamesPlayed    int
	Points         int
	CatorcesCalled int
	CatorcesMissed int
}

func ReadStatsFromFile(file string) (OverallStats, error) {
	os := OverallStats{}
	body, err := ioutil.ReadFile(file)

	if err != nil {
		return os, err
	}

	err = json.Unmarshal(body, &os)
	return os, err
}

func (os *OverallStats) Persist(file string) error {
	body, err := json.Marshal(os)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, body, 0644)
}

func (os *OverallStats) Dump() error {
	body, err := json.MarshalIndent(os, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
