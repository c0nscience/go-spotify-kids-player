package ha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type playRequest struct {
	EntityId         []string `json:"entity_id"`
	MediaContentId   string   `json:"media_content_id"`
	MediaContentType string   `json:"media_content_type"`
	Enqueue          string   `json:"enqueue"`
}

type joinRequest struct {
	EntityId     []string `json:"entity_id"`
	GroupMembers []string `json:"group_members"`
}

type basicRequest struct {
	EntityId []string `json:"entity_id"`
}

type Service string

const (
	PlayService   Service = "play_media"
	StopService   Service = "media_stop"
	JoinService   Service = "join"
	UnjoinService Service = "unjoin"
)

var accessToken = os.Getenv("HA_TOKEN")
var host = os.Getenv("HA_HOST")

func entityIds(rooms []string) []string {
	var res []string
	for _, room := range rooms {
		res = append(res, fmt.Sprintf("media_player.%s", room))
	}

	return res
}

func urlFor(service Service) string {
	return fmt.Sprintf("%s/api/services/media_player/%s", host, service)
}

func sendRequest(method string, url string, reqBody interface{}) error {
	b, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 200 {
		return errors.Join(fmt.Errorf("%s request to %s", method, url), err)
	}

	return nil
}

func Play(contentId string, rooms []string) error {
	reqBody := playRequest{
		EntityId:         entityIds(rooms),
		MediaContentId:   contentId,
		MediaContentType: "playlist",
		Enqueue:          "replace",
	}

	return sendRequest(http.MethodPost, urlFor(PlayService), &reqBody)
}

func Stop(rooms []string) error {
	reqBody := basicRequest{
		EntityId: entityIds(rooms),
	}

	return sendRequest(http.MethodPost, urlFor(StopService), &reqBody)
}

func Join(rooms []string) error {
	ids := entityIds(rooms)

	reqBody := joinRequest{
		EntityId:     ids,
		GroupMembers: ids,
	}

	return sendRequest(http.MethodPost, urlFor(JoinService), &reqBody)
}

func Unjoin(rooms []string) error {
	reqBody := basicRequest{
		EntityId: entityIds(rooms),
	}

	return sendRequest(http.MethodPost, urlFor(UnjoinService), &reqBody)
}
