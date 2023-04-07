package feeds

import (
	"context"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type localRulesPrv struct {
	rules ArticleValidationRules
}

func NewLocalRulesProvider(ctx context.Context, localPath string) (RulesProvider, error) {
	localPath = path.Clean(localPath)

	bb, err := os.ReadFile(localPath)
	if err != nil {
		return nil, fmt.Errorf("rules read error: %w", err)
	}

	var rules = make(ArticleValidationRules)
	if err := yaml.Unmarshal(bb, &rules); err != nil {
		return nil, fmt.Errorf("rules read error: %w", err)
	}
	return &localRulesPrv{
		rules: rules,
	}, nil
}

func (lp *localRulesPrv) Rules(ctx context.Context, topic string) (*ValidationRules, error) {
	rules, ok := lp.rules[topic]
	if !ok {
		return nil, fmt.Errorf("topic %s not found", topic)
	}
	return &rules, nil
}

func (lp *localRulesPrv) Shutdown(ctx context.Context) error {
	return nil
}
