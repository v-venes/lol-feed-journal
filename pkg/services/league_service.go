package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	leaguemodel "github.com/v-venes/feed_journal/pkg/models/league"
)

type LeagueService struct {
	key      string
	basePath string
	client   *http.Client
}

type NewLeagueServiceParams struct {
	Key      string
	BasePath string
}

type GetAccountIDParams struct {
	Name string
	Tag  string
}

type GetMatchesIDsParams struct {
	AccountID string
	From      uint32
	To        uint32
}

func NewLeagueService(params NewLeagueServiceParams) *LeagueService {
	return &LeagueService{
		key:      params.Key,
		basePath: params.BasePath,
		client:   &http.Client{},
	}
}

func (l *LeagueService) GetAccountID(params GetAccountIDParams) (*leaguemodel.Account, error) {
	path := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", l.basePath, params.Name, params.Tag, l.key)

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	accountInfo := &leaguemodel.Account{}

	err = json.NewDecoder(resp.Body).Decode(accountInfo)

	if err != nil {
		return nil, err
	}

	return accountInfo, nil
}

func (l *LeagueService) GetMatchesIDs(params GetMatchesIDsParams) ([]string, error) {
	path := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?startTime=%d&endTime=%d&api_key=%s", l.basePath, params.AccountID, params.From, params.To, l.key)

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var matchesIDs []string

	err = json.NewDecoder(resp.Body).Decode(&matchesIDs)

	if err != nil {
		return nil, err
	}

	return matchesIDs, nil
}

func (l *LeagueService) GetMatchDetails(matchId string) (*leaguemodel.Match, error) {
	path := fmt.Sprintf("%s/lol/match/v5/matches/%s?api_key=%s", l.basePath, matchId, l.key)

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	matchDetails := &leaguemodel.Match{}

	err = json.NewDecoder(resp.Body).Decode(matchDetails)

	if err != nil {
		return nil, err
	}

	return matchDetails, nil
}

func (l *LeagueService) DownloadChampionSquareImage(championName string) (*string, error) {
	// TODO: Refatorar essa parte
	path := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/15.2.1/img/champion/%s.png", championName)

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	filePath := fmt.Sprintf("/tmp/%s.png", championName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return &filePath, nil
}
