package stat

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// HumanReadableFilesize returns a human readable string representation of a file size.
func HumanReadableFilesize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}

func Duration(msg string, start time.Time, log *logrus.Logger) {
	log.Infof("%s took %s", msg, time.Since(start))
}
