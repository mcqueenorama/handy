package breaker

import (
	"testing"
	"time"
)

func TestErrorRateUnderThreshold(t *testing.T) {
	c := metric{}

	c.Success(0)
	c.Success(0)
	c.Success(0)
	c.Failure(0)
	c.Success(0)
	c.Failure(0)
	c.Success(0)
	c.Success(0)

	if r := c.Summary().rate; r == 0.0 {
		t.Errorf("expected error rate to be greater than zero,  got: %f in %+v", r, c.Summary())
	}
}

func TestErrorRateOverThreshold(t *testing.T) {
	c := metric{}

	c.Failure(0)
	c.Failure(0)
	c.Failure(0)
	c.Failure(0)
	c.Failure(0)
	c.Failure(0)
	c.Failure(0)
	c.Success(0)
	c.Success(0)

	if ex, s := 0.70, c.Summary(); s.rate < ex {
		t.Errorf("expected error rate to be over %d%%, got: %f in %+v", ex*100, s, s.rate)
	}
}

func TestErrorRateCalculatedFromLast5Seconds(t *testing.T) {
	defer afterTest()

	c := metric{}

	fakenow := time.Now()
	now = func() time.Time { return fakenow }

	// 77% error for 5 seconds
	for i := 0; i < 5; i++ {
		fakenow = fakenow.Add(time.Second)

		c.Failure(0)
		c.Failure(0)
		c.Success(0)
	}

	// 33.333% error for 5 seconds
	for i := 0; i < 5; i++ {
		fakenow = fakenow.Add(time.Second)

		c.Failure(0)
		c.Success(0)
		c.Success(0)
	}

	if ex, s := 0.34, c.Summary(); s.rate > ex {
		t.Errorf("expected error rate to be under %d%%, got: %f in %+v", int(ex*100), s.rate, s)
	}

}
