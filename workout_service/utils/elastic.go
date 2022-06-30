package utils

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pavel/workout_service/pkg/errors"
	"github.com/pavel/workout_service/pkg/logger"
)

type ElasticSearchResults struct {
	Total int    `json:"total"`
	Hits  []*any `json:"hits"`
}

func ValidateElasticResponse(res *esapi.Response, logger *logger.Logger) error {
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			logger.Error(fmt.Sprintf("json decoder error: ERROR %s", err.Error()))
			return errors.UnprocessableEntity
		}
		logger.Error(fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"]))
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}
