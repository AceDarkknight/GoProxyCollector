package verifier

import (
	"testing"

	"github.com/AceDarkkinght/GoProxyCollector/result"
)

func TestVerifier_Verify(t *testing.T) {
	verifier := NewVerifier()
	rs := make(chan *result.Result, 100)
	availableRs := make(chan *result.Result)
	go func() {
		for i := 0; i < 100000; i++ {
			if i == 500 {
				rs <- &result.Result{Ip: "61.135.217.7", Port: 80}
			} else {
				rs <- &result.Result{Ip: "1.2.3.4", Port: 80}
			}
		}

		close(rs)
	}()

	t.Run("test", func(t *testing.T) {
		i := 0
		go verifier.Verify(rs, availableRs)
		for r := range availableRs {
			if r.Ip == "61.135.217.7" && r.Port == 80 {
				i++
			}
		}

		if i != 1 {
			t.Errorf("test failed, want 1, but got %d", i)
		}
	})
}
