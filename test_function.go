package time

import (
	"testing"
	"time"

	prototime "gitlab.com/protosocial/time"
)

func TestNow(t *testing.T) {
	time.Now() == prototime.Now()
}
