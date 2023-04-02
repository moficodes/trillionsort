package fileops

import (
	"bytes"
	"context"
	"io"
	"os"
	"sort"
	"testing"
)

func BenchmarkReadDataScan(b *testing.B) {
	f, err := os.Open("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		_, err := ReadDataScan(f)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadData(b *testing.B) {
	f, err := os.Open("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		_, err := ReadData(f, 1000000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSortStrings(b *testing.B) {
	f, err := os.Open("testdata/input_small.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 10000)
		if err != nil {
			b.Fatal(err)
		}
		sort.Strings(data)
	}
}

func BenchmarkSortSlice(b *testing.B) {
	f, err := os.Open("testdata/input_small.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 10000)
		if err != nil {
			b.Fatal(err)
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i] < data[j]
		})
	}
}

func BenchmarkWriteDataLineByLine(b *testing.B) {
	f, err := os.Open("testdata/input_small.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 10000)
		if err != nil {
			b.Fatal(err)
		}
		w := io.Discard
		_, err = WriteDataLineByLine(w, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteDataLineByLine_large(b *testing.B) {
	f, err := os.Open("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 1_000_000)
		if err != nil {
			b.Fatal(err)
		}
		w := io.Discard
		_, err = WriteDataLineByLine(w, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteDataBatch(b *testing.B) {
	f, err := os.Open("testdata/input_small.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 10000)
		if err != nil {
			b.Fatal(err)
		}
		w := io.Discard
		_, err = WriteData(w, data, 1000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteDataBatch_large(b *testing.B) {
	f, err := os.Open("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		data, err := ReadData(f, 1_000_000)
		if err != nil {
			b.Fatal(err)
		}
		w := io.Discard
		_, err = WriteData(w, data, 10000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWrite10_000_seq(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		write(ctx, &buffer, 1, 10_000, 1, 17)
	}
}

func BenchmarkWriteFile10_000_seq(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		tmpfile, err := os.CreateTemp("", "")
		if err != nil {
			b.Fatal(err)
		}

		cfg := &Config{
			Goroutine:  1,
			Count:      10000,
			BufferSize: 1,
			LineLength: 17,
		}

		WriteToFile(ctx, tmpfile, cfg)
	}
}

func BenchmarkDiscard10_000_seq(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		write(ctx, io.Discard, 1, 10_000, 1, 17)
	}
}

func BenchmarkDiscard10_000_4workers(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		write(ctx, io.Discard, 4, 2_500, 1, 17)
	}
}

func BenchmarkDiscard10_000_50workers(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		write(ctx, io.Discard, 50, 200, 1, 17)
	}
}

func BenchmarkWriteFile10_000_4workers(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		tmpfile, err := os.CreateTemp("", "")
		if err != nil {
			b.Fatal(err)
		}

		cfg := &Config{
			Goroutine:  4,
			Count:      2_500,
			BufferSize: 1,
			LineLength: 17,
		}

		WriteToFile(ctx, tmpfile, cfg)
	}
}

func BenchmarkWriteFile10_000_50workers(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		tmpfile, err := os.CreateTemp("", "")
		if err != nil {
			b.Fatal(err)
		}

		cfg := &Config{
			Goroutine:  50,
			Count:      200,
			BufferSize: 1,
			LineLength: 17,
		}

		WriteToFile(ctx, tmpfile, cfg)
	}
}

func BenchmarkMergeSortedFiles(b *testing.B) {
	files, err := FindFilesInDir("testdata", "output")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		err := MergeSortedFiles(files, "testdata/output.txt", 1, 17)
		if err != nil {
			b.Fatal(err)
		}
	}
}
