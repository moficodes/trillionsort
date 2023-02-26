package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/moficodes/trillionsort/internal/fileops"
	"github.com/moficodes/trillionsort/internal/stat"
	"github.com/sirupsen/logrus"
)

var (
	count         int
	goroutine     int
	filename      string
	bufferSize    int
	ver           bool
	linelength    int
	filenameIndex string
	version       string = "v0.0.0"
	commit        string = "unknown"
	date          string = "unknown"
)

var log *logrus.Logger
var dataPerGoroutine int

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
	flag.IntVar(&count, "count", 0, "number of record to generate")
	flag.IntVar(&goroutine, "goroutine", 0, "number of goroutine to run")
	flag.StringVar(&filename, "file", "input.txt", "name of the file")
	flag.IntVar(&bufferSize, "buffer", 1, "buffer size in Mb")
	flag.IntVar(&linelength, "linelength", 17, "length of the line (length of each number + 1 for newline)")
	flag.StringVar(&filenameIndex, "fileindex", "", "name of the file index")
	flag.BoolVar(&ver, "version", false, "print version and exit")
}

func main() {
	defer stat.Duration("generate", time.Now(), log)
	flag.Parse()

	if ver {
		fmt.Printf("generate %s (%s) %s\n", version, commit, date)
		os.Exit(0)
	}

	if goroutine > count {
		goroutine = count
	}

	if goroutine == 0 {
		goroutine = runtime.GOMAXPROCS(-1)
	}

	if count == 0 {
		log.Error("no data to produce")
		os.Exit(1)
	}

	outputFile, err := fileops.GetFileName(filename, filenameIndex)
	if err != nil {
		log.Fatal(err)
	}

	dataPerGoroutine = count / goroutine
	count = count - (count % goroutine)

	log.Infof("total count: %d, goroutine: %d, gen per goroutine: %d", count, goroutine, dataPerGoroutine)

	err = fileops.DeleteFileIfExists(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	err = fileops.WriteToFile(context.Background(), outputFile, goroutine, dataPerGoroutine, bufferSize, linelength)
	if err != nil {
		log.Error(err)
	}
	log.Infof("total gen: %s, filename: %s", stat.HumanReadableFilesize(int64(count*linelength)), filename)
	log.Debug(stat.MemUsage())
}
