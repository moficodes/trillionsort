package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/moficodes/trillionsort/internal/fileops"
	"github.com/moficodes/trillionsort/internal/stat"
	"github.com/sirupsen/logrus"
)

var (
	filename   string
	count      int
	goroutine  int
	buffer     int
	parallel   bool
	linelength int
	ver        bool
	version    string = "v0.0.0"
	commit     string = "unknown"
	date       string = "unknown"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
	flag.StringVar(&filename, "file", "input.txt", "file name to split")
	flag.IntVar(&count, "count", 0, "split the file in these many files")
	flag.BoolVar(&ver, "version", false, "show version")
	flag.IntVar(&buffer, "buffer", 1, "buffer size in MB")
	flag.IntVar(&linelength, "linelength", 17, "length of each line (length of each number + 1 for newline)")
	flag.BoolVar(&parallel, "parallel", false, "split the file in parallel (default false)")
	flag.IntVar(&goroutine, "goroutine", runtime.GOMAXPROCS(-1), "number of concurrent workers")
}

func main() {
	defer stat.Duration("split", time.Now(), log)
	flag.Parse()

	if ver {
		fmt.Printf("splitfile %s (%s) %s\n", version, commit, date)
		os.Exit(0)
	}

	if count == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if goroutine > count {
		goroutine = count
	}

	filenamePrefix := strings.Split(filename, ".")[0]

	log.Infof("splitting file %s into %d files with buffer size %d MB", filename, count, buffer)
	var err error
	if parallel {
		err = fileops.SplitFileParallel(context.Background(), count, goroutine, buffer, linelength, filename, filenamePrefix)
	} else {
		err = fileops.SplitFile(count, buffer, linelength, filename, filenamePrefix)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log.Debug(stat.MemUsage())
}
