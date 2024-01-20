package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type fdata struct {
	lines int
	words int
	bytes int
	fname string
}

var (
	want_lines = false
	want_words = false
	want_bytes = false
	total      = fdata{fname: "total"}
	files      = []fdata{}
)

func print_stats(x fdata) {
	if want_lines {
		fmt.Printf("%8d", x.lines)
	}
	if want_words {
		fmt.Printf("%8d", x.words)
	}
	if want_bytes {
		fmt.Printf("%8d", x.bytes)
	}
	fmt.Printf(" %-s\n", x.fname)
}

// *os.File implements io.Reader
func process_file(f io.Reader, fname string) (fdata, error) {
	data := fdata{fname: fname}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		data.lines++
		data.words += len(strings.Fields(t))
		data.bytes += len(t) + 1
	}
	if scanner.Err() != nil {
		return data, scanner.Err()
	} else {
		return data, nil
	}
}

func main() {
	flag.BoolVar(&want_lines, "l", false, "Line count")
	flag.BoolVar(&want_words, "w", false, "Word count")
	flag.BoolVar(&want_bytes, "c", false, "Byte count")
	flag.Parse()

	if !want_lines && !want_words && !want_bytes {
		want_lines, want_words, want_bytes = true, true, true
	}

	if len(flag.Args()) == 0 {
		data, _ := process_file(os.Stdin, "-")
		print_stats(data)
		return
	}

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't open %s: %s\n", fname, err)
			continue
		}
		if data, err := process_file(f, fname); err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading %s: %s\n", fname, err)
		} else {
			files = append(files, data)
			total.lines += data.lines
			total.bytes += data.bytes
			total.words += data.words
		}
		f.Close()
	}
	for _, v := range files {
		print_stats(v)
	}
	if len(files) > 1 {
		print_stats(total)
	}
}
