package functional

import (
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/mohammadne/fesghel/internal/config"
	"github.com/mohammadne/fesghel/internal/entities"
	"github.com/mohammadne/fesghel/internal/urls"
)

var urlsService urls.Service

func TestMain(m *testing.M) {
	cfg := struct {
		URLs *urls.Config `required:"true"`
	}{}

	entities.LoadEnvironment(string(entities.EnvironmentLocal))
	if err := config.Load(&cfg); err != nil {
		log.Panicf("failed to load config: \n%v", err)
	}

	var err error
	urlsService, err = urls.NewService(cfg.URLs, zap.NewNop())
	if err != nil {
		log.Fatalf("failed to initialize urls: \n%v", err)
	}

	m.Run()
}
