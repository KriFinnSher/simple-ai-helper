package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"support/internal/config"
	"support/internal/usecase/knowledge"
)

type KnowledgeUseCase struct {
	repo knowledge.Repo
}

func NewKnowledgeInstance(repo knowledge.Repo) *KnowledgeUseCase {
	return &KnowledgeUseCase{repo: repo}
}

func (k *KnowledgeUseCase) ExtractIntent(ctx context.Context, message string) (string, error) {
	pythonServiceURL := fmt.Sprintf("http://%s:%s/predict", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	data := map[string]string{"message": message}
	req, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(pythonServiceURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var result map[string]string
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	intent, ok := result["intent"]
	if !ok {
		return "", err
	}

	if !k.repo.Exist(ctx, intent) {
		return "", errors.New("unable to find intent")
	}

	return intent, nil
}

func (k *KnowledgeUseCase) GetAnswer(ctx context.Context, intent string) (string, error) {
	if !k.repo.Exist(ctx, intent) {
		return "", errors.New("unable to find intent")
	}
	knowledgeInstance, err := k.repo.Get(ctx, intent)
	if err != nil {
		return "", err
	}
	return knowledgeInstance.Answer, nil
}
