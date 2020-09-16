package confs3

import (
	"fmt"
	"time"

	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

const (
	TIME_LAYOUT = `2006-01-02 15:04:05`
)

type Lifecycle struct {
}

func NewLifeCycle() *Lifecycle {
	return &Lifecycle{}
}

func (p *Lifecycle) ExpiresIn(prefix string, days int) (*lifecycle.Configuration, error) {
	config := lifecycle.NewConfiguration()
	rule, err := p.ExpireRule(prefix, days, "")
	if err != nil {
		return nil, err
	}
	config.Rules = []lifecycle.Rule{
		rule,
	}
	return config, nil
}

// ExpiresAt date , example 2006-01-02 15:04:05
func (p *Lifecycle) ExpiresAt(prefix string, date string) (*lifecycle.Configuration, error) {
	config := lifecycle.NewConfiguration()
	rule, err := p.ExpireRule(prefix, 0, date)
	if err != nil {
		return nil, err
	}
	config.Rules = []lifecycle.Rule{
		rule,
	}
	return config, nil
}

func (p *Lifecycle) ExpireRule(prefix string, days int, date string) (lifecycle.Rule, error) {

	if days < 0 {
		return lifecycle.Rule{}, fmt.Errorf("days (%d) < 0 ", days)
	}

	exp := lifecycle.Expiration{}
	id := ""
	if days >= 0 && date == "" {
		exp.Days = lifecycle.ExpirationDays(days)
		id = fmt.Sprintf("expired-in-%d", days)
	}

	if date != "" {
		d, err := time.Parse(TIME_LAYOUT, date)
		if err != nil {
			return lifecycle.Rule{}, err
		}
		exp.Date = lifecycle.ExpirationDate{Time: d}
		id = fmt.Sprintf("expired-at-%s", date)
	}

	rule := lifecycle.Rule{
		ID:         id,
		Status:     "Enabled",
		Expiration: exp,
		Prefix:     prefix,
	}

	return rule, nil
}
