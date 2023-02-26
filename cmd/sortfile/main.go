package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/moficodes/trillionsort/internal/fileops"
	"github.com/moficodes/trillionsort/internal/stat"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

var (
	input      string
	output     string
	index      string
	linelength int
	batchSize  int
	check      bool
	ver        bool
	version    string = "v0.0.0"
	commit     string = "unknown"
	date       string = "unknown"
)

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
	flag.StringVar(&input, "input", "input.txt", "file to read")
	flag.StringVar(&output, "output", "output.txt", "file to write")
	flag.StringVar(&index, "index", "", "index of file")
	flag.IntVar(&linelength, "linelength", 17, "length of line")
	flag.IntVar(&batchSize, "batchsize", 100000, "batch size")
	flag.BoolVar(&check, "check", false, "check if file is sorted")
	flag.BoolVar(&ver, "version", false, "print version and exit")
}

func main() {
	defer stat.Duration("sortfile", time.Now(), log)
	flag.Parse()

	if ver {
		fmt.Printf("sortfile %s (%s) %s\n", version, commit, date)
		os.Exit(0)
	}

	inputFile, err := fileops.GetFileName(input, index)
	if err != nil {
		log.Fatal(err)
	}

	input, err := os.Open(inputFile)
	if err != nil {
		os.Exit(1)
	}

	defer input.Close()
	fi, err := input.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fi.Size()

	log.Infof("File size: %s", stat.HumanReadableFilesize(fileSize))
	data, err := fileops.ReadDataScan(input)
	if err != nil {
		os.Exit(1)
	}

	if check {
		log.Info("Checking if file is sorted")
		if fileops.IsSorted(data) {
			log.Info("File is sorted")
			os.Exit(0)
		} else {
			log.Info("File is not sorted")
			os.Exit(1)
		}
	}

	outputFileName, err := fileops.GetFileName(output, index)
	if err != nil {
		log.Fatal(err)
	}

	err = fileops.DeleteFileIfExists(outputFileName)
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(data)

	outfile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	written, err := fileops.WriteData(outfile, data, batchSize)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Wrote %s to %s", stat.HumanReadableFilesize(int64(written)), output)
}
