package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// debug
func checkErr(err error, context string) {
	if err != nil {
		log.Printf("[%s] %s", context, err)
	}
}

// some base URLs
const (
	NA_SPEC_HOST         = "http://spectator.na.lol.riotgames.com:8088"
	NA_SPEC_OBS_CONSUMER = NA_SPEC_HOST + "/observer-mode/rest/consumer"
	NA_SPEC_FEATURED     = NA_SPEC_HOST + "/observer-mode/rest/featured"
)

// Holds info that can be used to build/update API URLs
type SpectatorURL struct {
	Base             string
	Consumer         string
	Featured         string
	GetGameMetaData  string
	GetLastChunkInfo string
	GetGameDataChunk string
	Token            int
	Platform         string
	GameId           int
}

func NewSpectatorURL() *SpectatorURL {
	s := &SpectatorURL{}
	s.Base = NA_SPEC_HOST
	s.Consumer = NA_SPEC_OBS_CONSUMER
	s.Featured = NA_SPEC_FEATURED
	s.Platform = "NA1"
	s.Token = 0
	s.GameId = 0
	s.Update(s.GameId, s.Token)
	return s
}

// Update/build the API URLs with new gameId or tokens
func (s *SpectatorURL) Update(game int, tok int) {
	s.GameId = game
	s.Token = tok
	gameIdStr := strconv.Itoa(game)
	tokStr := strconv.Itoa(tok)
	s.GetGameMetaData = s.Consumer + "/getGameMetaData/" +
		s.Platform + "/" + gameIdStr + "/0/token"
	s.GetLastChunkInfo = s.Consumer + "/getLastChunkInfo/" +
		s.Platform + "/" + gameIdStr + "/0/token"
	s.GetGameDataChunk = s.Consumer + "/getGameDataChunk/" +
		s.Platform + "/" + gameIdStr + "/" + tokStr + "/token/"
}

// json structs for unmarshalling responses
/////////////////////////////////////////////

// Featured game list
type Featured struct {
	GameList              []Game `"json:gameList"`
	ClientRefreshInterval int    `"json:clientRefreshInterval"`
}

type Game struct {
	Participants      []Participant    `"json:participants"`
	GameStartTime     int              `"json:gameStartTime"`
	GameQueueConfigId int              `"json:gameQueueConfigId"`
	GameType          string           `"json:gameType"`
	GameId            int              `"json:gameId"`
	Observers         Observer         `"json:observers"`
	BannedChampions   []BannedChampion `"json:bannedChampions"`
	GameTypeConfigId  int              `"json:gameTypeConfigId"`
	GameMode          string           `"json:gameMode"`
	GameLength        int              `"json:gameLength"`
	MapId             int              `"json:mapId"`
	PlatformId        string           `"json:platformId"`
}

type Participant struct {
	TeamId        int    `"json:teamId"`
	Bot           bool   `"json:bot"`
	SummonerName  string `"json:summonerName"`
	SkinIndex     int    `"json:skinIndex"`
	Spell1Id      int    `"json:spell1Id"`
	Spell2Id      int    `"json:spell2Id"`
	ChampionId    int    `"json:championId"`
	ProfileIconId int    `"json:profileIconId"`
}

type Observer struct {
	EncryptionKey string `"json:encryptionKey"`
}

type BannedChampion struct {
	PickTurn   int `"json:pickTurn"`
	TeamId     int `"json:teamId"`
	ChampionId int `"json:championId"`
}

// getGameMetaData response
type GameMetaData struct {
	LastKeyFrameId               int               `"json:lastKeyFrameId"`
	KeyFrameTimeInterval         int               `"json:keyFrameTimeInterval"`
	Port                         int               `"json:port"`
	EncryptionKey                string            `"json:encryptionKey"`
	PendingAvailableChunkInfo    []PendingChunk    `"json:pendingAvailableChunkInfo"`
	LastChunkId                  int               `"json:lastChunkId"`
	CreateTime                   string            `"json:createTime"`
	StartGameChunkId             int               `"json:startGameChunkId"`
	FeaturedGame                 bool              `"json:featuredGame"`
	DecodedEncryptionKey         string            `"json:decodedEncryptionKey"`
	GameEnded                    bool              `"json:gameEnded"`
	DelayTime                    int               `"json:delayTime"`
	GameLength                   int               `"json:gameLength"`
	StartTime                    string            `"json:startTime"`
	InterestScore                int               `"json:interestScore"`
	EndStartupChunkId            int               `"json:endStartupChunkId"`
	ClientBackFetchingEnabled    bool              `"json:clientBackFetchingEnabled"`
	PendingAvailableKeyFrameInfo []PendingKeyFrame `"json:pendingAvailableKeyFrameInfo"`
	ClientAddedLag               int               `"json:clientAddedLag"`
	ChunkTimeInterval            int               `"json:chunkTimeInterval"`
	GameKey                      GameKeyInfo       `"json:gameKey"`
	GameServerAddress            string            `"json:gameServerAddress"`
	ClientBackFetchingFreq       int               `"json:clientBackFetchingFreq"`
}

type PendingChunk struct {
	ReceivedTime string `"json:receivedTime"`
	Duration     int    `"json:duration"`
	Id           int    `"json:id"`
}

type PendingKeyFrame struct {
	ReceivedTime string `"json:receivedTime"`
	NextChunkId  int    `"json:nextChunkId"`
	Id           int    `"json:id"`
}

type GameKeyInfo struct {
	PlatformId string `"json:platformId"`
	GameId     int    `"json:gameId"`
}

// getLastChunkInfo response
type LastChunkInfo struct {
	KeyFrameId         int `"json:keyFrameId"`
	Duration           int `"json:duration"`
	NextChunkId        int `"json:nextChunkId"`
	NextAvailableChunk int `"json:nextAvailableChunk"`
	ChunkId            int `"json:chunkId"`
	AvailableSince     int `"json:availableSince"`
	EndGameChunkId     int `"json:endGameChunkId"`
	EndStartupChunkId  int `"json:endStartupChunkId"`
	StartGameChunkId   int `"json:startGameChunkId"`
}

func main() {
	client := &http.Client{}
	resp, err := client.Get(NA_SPEC_FEATURED)
	checkErr(err, "get")
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, "read")

	var featuredGames Featured
	json.Unmarshal(body, &featuredGames)
	log.Println(featuredGames.GameList[0].GameId)

	s := NewSpectatorURL()
	s.Update(featuredGames.GameList[0].GameId, 0)
	log.Println(s.GetGameMetaData)
	resp, err = client.Get(s.GetGameMetaData)

	checkErr(err, "get metadata")
	log.Printf("%#v", resp)

	body, err = ioutil.ReadAll(resp.Body)
	checkErr(err, "read metadata")
	var gameMeta GameMetaData
	json.Unmarshal(body, &gameMeta)

	log.Printf("%#v", gameMeta)

	resp, err = client.Get(s.GetLastChunkInfo)
	checkErr(err, "get last chunk info")
	body, err = ioutil.ReadAll(resp.Body)
	checkErr(err, "read last chunk info")

	var last LastChunkInfo
	json.Unmarshal(body, &last)
	log.Printf("%#v", last)

	s.Update(s.GameId, last.ChunkId)
	log.Println(s.GetGameDataChunk)

	resp, err = client.Get(s.GetGameDataChunk)
	checkErr(err, "get data chunk")

	body, err = ioutil.ReadAll(resp.Body)
	checkErr(err, "read data chunk")

	log.Printf("%#v", body)
}
