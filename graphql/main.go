package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	if err := envconfig.Process(context.TODO(), &cfg); err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(s.ToExecutableSchema())
	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("Graphql playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
