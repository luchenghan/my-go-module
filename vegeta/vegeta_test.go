package vegeta

import (
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestAttackTargets(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 10 * time.Second

	ts := []vegeta.Target{}

	AttackTargets("test", rate, duration, ts...)
}
func TestAttackTargeter(t *testing.T) {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 10 * time.Second

	tgtr := func(tgt *vegeta.Target) error {
		return nil
	}

	AttackTargeter("test", rate, duration, tgtr)
}
