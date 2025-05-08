package vegeta

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func AttackTargeter(name string, r vegeta.Rate, d time.Duration, t vegeta.Targeter) {
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for result := range attacker.Attack(t, r, d, name) {
		metrics.Add(result)
	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latency (mean): %s\n", metrics.Latencies.Mean)
	fmt.Printf("Throughput: %.2f req/s\n", metrics.Throughput)
}

func AttackTargets(name string, r vegeta.Rate, d time.Duration, ts ...vegeta.Target) {
	AttackTargeter(name, r, d, vegeta.NewStaticTargeter(ts...))
}
