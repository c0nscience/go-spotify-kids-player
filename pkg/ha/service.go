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

type unjoinRequest struct {
	EntityId []string `json:"entity_id"`
}

var accessToken = os.Getenv("HA_TOKEN")
var host = os.Getenv("HA_HOST")

func toEntityIds(rooms []string) []string {
	var entityIds []string
	for _, room := range rooms {
		entityIds = append(entityIds, fmt.Sprintf("media_player.%s", room))
	}

	return entityIds
}

func Play(contentId string, rooms []string) error {
	entityIds := toEntityIds(rooms)

	reqBody := playRequest{
		EntityId:         entityIds,
		MediaContentId:   contentId,
		MediaContentType: "playlist",
		Enqueue:          "replace",
	}

	b, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/services/media_player/play_media", host), bytes.NewReader(b))
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
		return errors.New("could start playing")
	}

	return nil
}

func Join(rooms []string) (string, error) {
	entityIds := toEntityIds(rooms)

	reqBody := joinRequest{
		EntityId:     entityIds,
		GroupMembers: entityIds,
	}

	b, err := json.Marshal(&reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/services/media_player/join", host), bytes.NewReader(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 200 {
		return "", errors.New("could not join the rooms")
	}

	return rooms[0], nil
}

func Unjoin(rooms []string) error {
	entityIds := toEntityIds(rooms)

	reqBody := unjoinRequest{
		EntityId: entityIds,
	}

	b, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/services/media_player/unjoin", host), bytes.NewReader(b))
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
		return errors.New("could not unjoin the rooms")
	}

	return nil
}
