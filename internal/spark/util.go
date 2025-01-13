package spark

import (
	"fmt"
	"time"
)

const (
	RFC3339Milli = "2006-01-02T15:04:05.000GMT"
)

func calcDateDifferenceMillis(startDate string, endDate string) (
	int64, error,
) {
	start, err := time.Parse(RFC3339Milli, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date: %s", startDate)
	}

	end, err := time.Parse(RFC3339Milli, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date: %s", endDate)
	}

	return end.Sub(start).Milliseconds(), nil
}
