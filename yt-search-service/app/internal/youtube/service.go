package youtube

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/serjic/yt-search/pkg/client/youtube"
	"github.com/serjic/yt-search/pkg/logger"
)

type service struct {
	logger *logger.Logger
	client youtube.Client
}

type Service interface {
	FindTrackByName(ctx context.Context, trackName string) (string, error)
}

func NewService(logger *logger.Logger, client youtube.Client) Service {
	return &service{
		logger: logger,
		client: client,
	}
}

func (s *service) FindTrackByName(ctx context.Context, trackName string) (string, error) {
	response, err := s.client.SearchTrack(ctx, trackName)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		s.logger.Error(responseData["error"].(map[string]interface{})["message"].(string))
		return "", fmt.Errorf("failed")
	}

	/* TODO */
	a := responseData["items"].([]interface{})
	b := a[0].(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string)

	return fmt.Sprintf("https://music.youtube.com/watch?v=%s", b), nil
}
