package ha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func url(service Service) string {
	return fmt.Sprintf("%s/api/services/media_player/%s", host, service)
}

func sendRequest(method string, url string, reqBody interface{}) (io.ReadCloser, error) {
	var req *http.Request
	if reqBody != nil {
		b, err := json.Marshal(&reqBody)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, url, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
	} else {
		r, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		req = r
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 200 {
		return nil, errors.Join(fmt.Errorf("%s request to %s", method, url), err)
	}

	return resp.Body, nil
}

func Play(contentId string, rooms []string) error {
	reqBody := playRequest{
		EntityId:         entityIds(rooms),
		MediaContentId:   contentId,
		MediaContentType: "playlist",
		Enqueue:          "replace",
	}

	_, err := sendRequest(http.MethodPost, url(PlayService), &reqBody)
	return err
}

func Stop(rooms []string) error {
	reqBody := basicRequest{
		EntityId: entityIds(rooms),
	}

	_, err := sendRequest(http.MethodPost, url(StopService), &reqBody)
	return err
}

func AvailableRooms() ([]string, error) {

	bdy, err := sendRequest(http.MethodGet, fmt.Sprintf("%s/api/states", host), nil)
	respBdy := []map[string]interface{}{}
	err = json.NewDecoder(bdy).Decode(&respBdy)
	if err != nil {
		return nil, err
	}

	res := []string{}

	for _, m := range respBdy {
		id, ok := m["entity_id"].(string)
		if !ok || !strings.HasPrefix(id, "media_player") {
			continue
		}

		state, ok := m["state"].(string)
		if !ok || state == "unavailable" {
			continue
		}

		res = append(res, strings.TrimPrefix(id, "media_player."))
	}

	return res, err
}

func Join(rooms []string) error {
	ids := entityIds(rooms)

	reqBody := joinRequest{
		EntityId:     ids,
		GroupMembers: ids,
	}

	_, err := sendRequest(http.MethodPost, url(JoinService), &reqBody)
	return err
}

func Unjoin(rooms []string) error {
	reqBody := basicRequest{
		EntityId: entityIds(rooms),
	}

	_, err := sendRequest(http.MethodPost, url(UnjoinService), &reqBody)
	return err
}
