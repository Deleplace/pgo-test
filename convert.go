package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
	"sync"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}
	folder := flag.Arg(0)

	if *concurrent {
		processConcurrent(folder)
	} else {
		processSequential(folder)
	}
}

var concurrent = flag.Bool("concurrent", false, "process files concurrently")

func processSequential(d string) {
	_ = os.Mkdir(d+"/png", 0777)

	entries := must(os.ReadDir(d))
	fmt.Println("Processing", len(entries), "sequentially")
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if n := strings.ToLower(entry.Name()); !strings.HasSuffix(n, ".jpg") && !strings.HasSuffix(n, ".jpeg") {
			continue
		}
		src := d + "/" + entry.Name()
		fmt.Println("Decoding", src)
		dest := d + "/png/" + entry.Name() + ".png"
		convertJpgFileToPngFile(src, dest)
		// fmt.Println(dest)
	}
	fmt.Println("All done.")
}

func processConcurrent(d string) {
	_ = os.Mkdir(d+"/png", 0777)

	var wg sync.WaitGroup
	entries := must(os.ReadDir(d))
	fmt.Println("Processing", len(entries), "concurrently")
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if n := strings.ToLower(entry.Name()); !strings.HasSuffix(n, ".jpg") && !strings.HasSuffix(n, ".jpeg") {
			continue
		}
		src := d + "/" + entry.Name()
		dest := d + "/png/" + entry.Name() + ".png"
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Decoding", src)
			convertJpgFileToPngFile(src, dest)
			// fmt.Println("Wrote", dest)
		}()
	}
	wg.Wait()
	fmt.Println("All done.")
}

func convertJpgFileToPngFile(src, dest string) {
	data := must(os.ReadFile(src))
	pngdata := convertJpgToPng(data)
	os.WriteFile(dest, pngdata, 0777)
}

func convertJpgToPng(src []byte) []byte {
	buffer := bytes.NewBuffer(src)
	img, _ := must2(image.Decode(buffer))
	var pngdata bytes.Buffer
	must0(png.Encode(&pngdata, img))
	return pngdata.Bytes()
}

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func must[V any](v V, err error) V {
	must0(err)
	return v
}

func must2[V, W any](v V, w W, err error) (V, W) {
	must0(err)
	return v, w
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <flags> <folder>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}
