package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pavel/workout_service/config"
)

func InitElastic(cfg *config.Config) (error, *elasticsearch.Client) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("%s://%s%s", cfg.Elastic.Network, cfg.Elastic.Host, cfg.Elastic.Port)},
		Username:  cfg.Elastic.Username,
		Password:  cfg.Elastic.Password,
	})
	if err != nil {
		return err, nil
	}
	_, err = es.Info()
	if err != nil {
		return err, nil
	}

	return nil, es
}
