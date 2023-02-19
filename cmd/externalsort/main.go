package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/moficodes/trillionsort/internal/fileops"
	"github.com/moficodes/trillionsort/internal/stat"
	"github.com/sirupsen/logrus"
)

var (
	inputDir   string
	pattern    string
	output     string
	bufferSize int
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
	flag.StringVar(&inputDir, "dir", "", "directory to search for files")
	flag.StringVar(&pattern, "pattern", "", "file pattern to join, e.g. input or input_ for input_00001.txt")
	flag.StringVar(&output, "output", "", "output file name")
	flag.IntVar(&bufferSize, "buffer", 1, "buffer size in Mb")
	flag.IntVar(&linelength, "linelength", 17, "line length in bytes")
	flag.BoolVar(&ver, "version", false, "print version and exit")
}

func main() {
	defer stat.Duration("externalsort", time.Now(), log)
	flag.Parse()

	if ver {
		fmt.Printf("externalsort %s (%s) %s\n", version, commit, date)
		os.Exit(0)
	}

	if pattern == "" || output == "" {
		flag.Usage()
		os.Exit(1)
	}

	if inputDir == "" {
		inputDir = "."
	}

	files, err := fileops.FindFilesInDir(inputDir, pattern)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("merging %d files", len(files))
	fileops.MergeSortedFiles(files, output, bufferSize, linelength)
	log.Info(stat.MemUsage())
}
