package main

import "github.com/mohammadne/fesghel/pkg/databases/postgres"

type Config struct {
	Postgres *postgres.Config `required:"true"`
}
