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
			rs <- &result.Result{Ip: "1.2.3.4", Port: 80}
		}

		close(rs)
	}()

	t.Run("test", func(t *testing.T) {
		i := 0
		go verifier.Verify(rs, availableRs)
		for r := range availableRs {
			if r.Ip != "" {
				i++
			}
		}

		if i != 0 {
			t.Errorf("test failed, want 0, but got %d", i)
		}
	})
}

func BenchmarkVerifier_Verify(b *testing.B) {
	verifier := NewVerifier()
	rs := make(chan *result.Result, 1000)
	availableRs := make(chan *result.Result, 1)
	go func() {
		for i := 0; i < 1000; i++ {
			rs <- &result.Result{Ip: "1.2.3.4", Port: 80}
		}

		close(rs)
	}()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		verifier.Verify(rs, availableRs)
	}

	for _ = range availableRs {

	}
}
