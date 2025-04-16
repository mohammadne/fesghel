package main

import (
	"github.com/mohammadne/fesghel/internal/urls"
	"github.com/mohammadne/fesghel/pkg/observability/logger"
)

type Config struct {
	URLs   *urls.Config   `required:"true"`
	Logger *logger.Config `required:"true"`
}
