package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"cloud.google.com/go/storage"
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
	bucket        string
	objectStorage bool
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
	flag.IntVar(&goroutine, "goroutine", 1, "number of goroutine to run")
	flag.StringVar(&filename, "file", "input.txt", "name of the file")
	flag.IntVar(&bufferSize, "buffer", 1, "buffer size in Mb")
	flag.IntVar(&linelength, "linelength", 17, "length of the line (length of each number + 1 for newline)")
	flag.StringVar(&filenameIndex, "fileindex", "", "name of the file index")
	flag.StringVar(&bucket, "bucket", "", "name of the bucket")
	flag.BoolVar(&objectStorage, "objectstorage", false, "use object storage")
	flag.BoolVar(&ver, "version", false, "print version and exit")
}

func write(ctx context.Context, filename string, cfg *fileops.Config) error {
	outputFile, err := fileops.GetFileName(filename, filenameIndex)
	if err != nil {
		return err
	}

	log.Infof("total count: %d, goroutine: %d, gen per goroutine: %d", cfg.Count, cfg.Goroutine, cfg.DataPerGoroutine)

	var writer io.WriteCloser
	if objectStorage {
		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		writer = client.Bucket(cfg.Bucket).Object(outputFile).NewWriter(ctx)
	} else {
		err = fileops.DeleteFileIfExists(outputFile)
		if err != nil {
			return err
		}

		writer, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
	}

	defer writer.Close()

	return fileops.WriteSync(ctx, writer, cfg)
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

	dataPerGoroutine = count / goroutine
	count = count - (count % goroutine)

	cfg := &fileops.Config{
		Count:            count,
		Goroutine:        goroutine,
		DataPerGoroutine: dataPerGoroutine,
		BufferSize:       bufferSize,
		LineLength:       linelength,
		Bucket:           bucket,
	}

	if err := write(context.Background(), filename, cfg); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Infof("total gen: %s, filename: %s", stat.HumanReadableFilesize(int64(count*linelength)), filename)

	log.Debug(stat.MemUsage())
}
