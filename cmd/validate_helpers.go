package cmd

import (
	"context"

	"github.com/govinda-attal/winning11/feeds"
)

func ruleProvider(ctx context.Context) (feeds.RulesProvider, error) {

	if cfg.Validator.Local {
		return feeds.NewLocalRulesProvider(ctx, cfg.Validator.RulesPath)
	}
	db := cfg.Validator.DB
	return feeds.NewDBRulesProvider(ctx, db.URI, db.Name)
}
