package fileops

import (
	"bufio"
	"container/heap"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func copyChunk(in io.Reader, out io.ReaderFrom, n int64) (int64, error) {
	// ReaderFrom is a Writer that has the "Read from..." capability
	part := io.LimitReader(in, n)
	return out.ReadFrom(part)
}

func split(count, bufferMB, linelength int, f io.ReadSeeker, fileSize int64, filenamePrefix string) error {
	// each line is 17 bytes
	// so we can calculate the number of lines per chunk
	linesPerChunk := int((fileSize / int64(linelength)) / int64(count))
	// each chunk is 17 bytes per line (16 digits + 1 newline)
	// we need to do it this way to avoid any rounding errors
	// for example if we have say 100 lines. that is 1700 bytes
	// if we want to split it into 3 files. then each file should have 566 bytes
	// 566 bytes does not divide in 17 bytes per line.
	// instead if we calculate lines per chunk it comes to be 33
	// then reach chunk size is 660 bytes exactly
	chunkSize := linesPerChunk * linelength

	for i := 0; i < count; i++ {
		_, err := f.Seek(int64(chunkSize*i), io.SeekStart)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(fmt.Sprintf("%s_%04d.txt", filenamePrefix, i+1), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		// TODO use bufio to take "bufferMB" into account

		if i == count-1 {
			// Write everything we have left
			// The last file may be larger than the previous chunks!
			_, err := file.ReadFrom(f)
			return err
		}

		_, err = copyChunk(f, file, int64(chunkSize))
		if err != nil {
			return err
		}

		if err = file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func SplitFile(count, buffer, linelength int, filename, filenamePrefix string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	fileSize := fi.Size()
	return split(count, buffer, linelength, f, fileSize, filenamePrefix)
}

func SplitFileParallel(ctx context.Context, count, goroutine, bufferMB, linelength int, filename, filenamePrefix string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	fileSize := fi.Size()

	// see logic in split function
	linesPerChunk := int((fileSize / int64(linelength)) / int64(count))
	chunkSize := linesPerChunk * linelength
	errs, _ := errgroup.WithContext(ctx)
	errs.SetLimit(goroutine)

	for i := 0; i < count; i++ {
		i := i
		errs.Go(func() error {
			source, err := os.Open(filename)
			if err != nil {
				return err
			}
			destination, err := os.OpenFile(fmt.Sprintf("%s_%04d.txt", filenamePrefix, i+1), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil
			}
			// TODO use bufio to take "bufferMB" into account
			_, err = source.Seek(int64(chunkSize*i), io.SeekStart)
			if err != nil {
				return err
			}

			if i == count-1 {
				// Write everything we have left
				// The last file may be larger than the previous chunks!
				_, err := destination.ReadFrom(source)
				return err
			}

			_, err = copyChunk(source, destination, int64(chunkSize))
			if err != nil {
				return err
			}

			source.Close()
			return destination.Close()
		})
	}
	return errs.Wait()
}

func FindFilesInDir(dirPath string, pattern string) ([]string, error) {
	// Create the file pattern by joining the directory path and the pattern.
	filePattern := filepath.Join(dirPath, pattern) + "*"

	// Use the Glob function to find all files that match the pattern.
	matchingFiles, err := filepath.Glob(filePattern)
	if err != nil {
		return nil, err
	}

	return matchingFiles, nil
}

func IsSorted(data []string) bool {
	for i := 0; i < len(data)-1; i++ {
		if data[i] > data[i+1] {
			return false
		}
	}
	return true
}

func CopyFiles(files []string, output string, buffersize int) (int64, error) {
	// Open the output file for writing.
	dest, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return -1, err
	}
	defer dest.Close()

	var total int64 = 0

	// Loop through the files.
	for _, file := range files {
		// Open the source file for reading.
		src, err := os.Open(file)
		if err != nil {
			return -1, err
		}
		defer src.Close()

		// Create a buffered reader for the source file.
		sourceReader := bufio.NewReaderSize(src, buffersize)

		// Create a buffered writer for the destination file.
		destWriter := bufio.NewWriterSize(dest, buffersize)
		defer destWriter.Flush()

		// Copy the file contents to the destination file.
		read, err := io.Copy(destWriter, sourceReader)
		total += read
		if err != nil {
			return -1, err
		}
	}

	return int64(total), nil
}

func ReadDataScan(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res, scanner.Err()
}

func ReadData(r io.Reader, linecount int64) ([]string, error) {
	scanner := bufio.NewScanner(r)
	res := make([]string, linecount)
	idx := 0
	for scanner.Scan() {
		res[idx] = scanner.Text()
		idx++
	}
	return res, scanner.Err()
}

func readInt(r io.Reader) ([]uint64, error) {
	scanner := bufio.NewScanner(r)
	var res []uint64
	for scanner.Scan() {
		i, err := strconv.ParseUint(scanner.Text(), 16, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, scanner.Err()
}

func WriteDataLineByLine(w io.Writer, data []string) (int, error) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()
	written := 0
	for _, d := range data {
		count, err := writer.WriteString(d + "\n")
		written += count
		if err != nil {
			return -1, err
		}
	}
	return written, nil
}

func WriteData(w io.Writer, data []string, batchSize int) (int, error) {
	idx := 0
	written := 0
	for idx+batchSize+1 < len(data) {
		write, err := w.Write([]byte(strings.Join(data[idx:idx+batchSize], "\n") + "\n"))
		written += write
		if err != nil {
			return -1, err
		}
		idx += batchSize
	}
	write, err := w.Write([]byte(strings.Join(data[idx:], "\n") + "\n"))
	written += write
	return written, err
}

func WriteToFile(ctx context.Context, filename string, goroutines, dataPerGoroutine, bufferSize, lineLength int) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	bufferByteSize := bufferSize * 1024 * 1024
	bf := bufio.NewWriterSize(f, bufferByteSize)
	err = write(ctx, bf, goroutines, dataPerGoroutine, bufferSize, lineLength)
	if err != nil {
		return err
	}
	return bf.Flush()
}

func write(ctx context.Context, w io.Writer, goroutines, dataPerGoroutine, bufferSize, linelength int) error {
	errs, _ := errgroup.WithContext(ctx)
	var filelock sync.Mutex
	n := bufferSize * 1024 * 4 // number of lines in 1 buffered batch

	for i := 0; i < goroutines; i++ {
		errs.Go(func() error {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomBytes := make([]byte, 8*n)
			r.Read(randomBytes)
			randomHexDigits := make([]byte, (linelength-1)*n)
			outputBuffer := make([]byte, 0, linelength*n+1) // 1 '\n' after every 16 digits

			flushRefreshReuse := func() error {
				// Flush to w
				filelock.Lock()
				_, err := w.Write(outputBuffer)
				filelock.Unlock()
				// Refresh work buffers with new random bytes
				r.Read(randomBytes)
				hex.Encode(randomHexDigits, randomBytes)
				// Reuse output buffer
				outputBuffer = outputBuffer[:0]
				return err
			}

			for j := 0; j < dataPerGoroutine; j++ {
				k := j % n
				if k == 0 {
					if err := flushRefreshReuse(); err != nil {
						return err
					}
				}
				outputBuffer = append(outputBuffer, randomHexDigits[(linelength-1)*k:(linelength-1)*k+(linelength-1)]...)
				outputBuffer = append(outputBuffer, '\n')
			}

			return flushRefreshReuse()
		})
	}

	return errs.Wait()
}

func MergeSortedFiles(fileNames []string, outputFileName string, bufferSize, linelength int) error {
	files := make([]*os.File, len(fileNames))
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	// Open all the input files
	for i, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		files[i] = file
	}

	// Create the output file
	outFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Initialize the min heap with the first value from each file
	h := &minHeap{}

	limit := bufferSize * 1024 * 1024
	outputBuffer := make([]byte, 0, limit)

	for _, file := range files {
		buf := make([]byte, 17)
		_, err := file.Read(buf)
		if err != nil {
			continue
		}

		// Remove the newline character from the line
		heap.Push(h, &fileItem{
			file: file,
			val:  buf,
		})
	}

	flushAndReuse := func() error {
		_, err := outFile.Write(outputBuffer)
		outputBuffer = outputBuffer[:0]
		return err
	}

	// Pop the smallest value from the heap and write it to the output file
	for h.Len() > 0 {
		item := heap.Pop(h).(*fileItem)

		outputBuffer = append(outputBuffer, item.val...)
		if len(outputBuffer) >= limit {
			if err := flushAndReuse(); err != nil {
				return err
			}

		}
		// Read the next value from the file and add it to the heap
		buf := make([]byte, 17)
		_, err := item.file.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			return err
		}

		heap.Push(h, &fileItem{
			file: item.file,
			val:  buf,
		})
	}

	return flushAndReuse()
}

func ReadLine(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
}

func GetFileName(filename, fileindex string) (string, error) {
	if fileindex != "" {
		index, err := strconv.Atoi(fileindex)
		if err != nil {
			return "", err
		}
		paths := strings.Split(filename, "/")
		file := paths[len(paths)-1]
		name := strings.Split(file, ".")[0]
		ext := strings.Split(file, ".")[1]
		if len(paths) == 1 {
			filename = fmt.Sprintf(("%s_%04d.%s"), name, index, ext)
		} else {
			filename = fmt.Sprintf(("%s/%s_%04d.%s"), strings.Join(paths[:len(paths)-1], "/"), name, index, ext)
		}
	}
	return filename, nil
}
