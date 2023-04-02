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

func run() error {
	defer stat.Duration("sortfile", time.Now(), log)

	if ver {
		fmt.Printf("sortfile %s (%s) %s\n", version, commit, date)
		return nil
	}

	inputFile, err := fileops.GetFileName(input, index)
	if err != nil {
		return err
	}

	log.Infof("Reading file: %s", inputFile)

	input, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	defer input.Close()
	fi, err := input.Stat()
	if err != nil {
		return err
	}
	fileSize := fi.Size()

	log.Infof("File size: %s", stat.HumanReadableFilesize(fileSize))
	if check {
		log.Info("Checking if file is sorted")
		if fileops.IsSorted(input) {
			log.Info("File is sorted")
			return nil
		} else {
			log.Info("File is not sorted")
			return nil
		}
	}
	data, err := fileops.ReadDataScan(input)
	if err != nil {
		return err
	}

	outputFileName, err := fileops.GetFileName(output, index)
	if err != nil {
		return err
	}

	err = fileops.DeleteFileIfExists(outputFileName)
	if err != nil {
		return err
	}

	sort.Strings(data)

	log.Infof("Writing to file: %s", outputFileName)

	outfile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer outfile.Close()

	written, err := fileops.WriteData(outfile, data, batchSize)
	if err != nil {
		return err
	}

	log.Infof("Wrote %s to %s", stat.HumanReadableFilesize(int64(written)), outputFileName)
	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
