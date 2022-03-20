package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/serj1c/tbot/pkg/client/imgur"
	"github.com/serj1c/tbot/pkg/logger"
)

type imgurService struct {
	logger *logger.Logger
	client imgur.Client
}

type ImgurService interface {
	ShareImage(ctx context.Context, image []byte) (string, error)
}

func NewImgurService(logger *logger.Logger, client imgur.Client) ImgurService {
	return &imgurService{
		logger: logger,
		client: client,
	}
}

func (i *imgurService) ShareImage(ctx context.Context, image []byte) (string, error) {
	response, err := i.client.UploadImage(ctx, image)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()
	var responseData map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", nil
	}

	if response.StatusCode != 200 {
		i.logger.Error(responseData)
		return "", fmt.Errorf("failed to upload an image")
	}

	return responseData["data"].(map[string]interface{})["link"].(string), nil
}
