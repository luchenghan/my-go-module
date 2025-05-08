package vegeta

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var host = "https://egame-uat.idc.pstdsf.com"

var method = "POST"

var url1 = fmt.Sprintf("%s/gpa/api/non/v1/auth/login", host)
var url2 = fmt.Sprintf("%s/gpa/api/non/v1/auth/register", host)
var url3 = fmt.Sprintf("%v/gpa/api/v1", host)

var contentType []string = []string{
	"application/json",
}

var secret1 = "pc|<+^UE;Y5B$p-!$/KS$B&nQ5@F<wV#"
var secret2 = "7FAE582C79E3BBC94D974519C7963"

func GPA() {
	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Minute

	//targets := []vegeta.Target{}
	//targeter := vegeta.NewStaticTargeter(targets...)

	targeter := func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		b := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "erictest",
			Password: "erictest",
		}

		body, err := json.Marshal(b)
		if err != nil {
			return err
		}

		//x := xid.New()
		//
		//us := fmt.Sprintf("username%s", x.String())
		//body2 := map[string]any{
		//	"Request": map[string]any{
		//		"username":    us,
		//		"rawEcId":     "1",
		//		"rawEcName":   "testEcSite",
		//		"rawEcUserId": "testEcSite/" + us,
		//		"isTrial":     true,
		//	},
		//	"Action": "CreateMember",
		//}
		//
		//b2, _ := Encode(body2, secret)

		//b := struct {
		//	Username string `json:"username"`
		//	Password string `json:"password"`
		//}{
		//	Username: fmt.Sprintf("user%d", time.Now().UnixNano()),
		//	Password: "password",
		//}

		tgt.Method = method
		tgt.URL = url1
		tgt.Body = body
		tgt.Header = map[string][]string{
			"Content-Type": contentType,
			//"Content-Type":  {"text/plain"},
			//"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlY1NpdGVJZCI6IjEiLCJpYXQiOjE3NDUzNzY3MTEsImlzcyI6ImdvLXB1YmxpYy1hcGkifQ.9N1InUqKQZitwY1fpMTqUCOKOVa_ZO0OaepqsLZOtsw"},
		}

		return nil
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for result := range attacker.Attack(targeter, rate, duration, "test") {
		metrics.Add(result)
	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latency (mean): %s\n", metrics.Latencies.Mean)
	fmt.Printf("Throughput: %.2f req/s\n", metrics.Throughput)
}

func Encode(payloads map[string]interface{}, secret string) ([]byte, error) {
	builder := jwt.NewBuilder()
	for k, v := range payloads {
		builder.Claim(k, v)
	}

	tok, err := builder.IssuedAt(time.Now()).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build jwt: %w", err)
	}

	// Sign a JWT!
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, []byte(secret)))
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %w", err)
	}

	return signed, nil
}
