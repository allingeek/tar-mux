// Copyright 2017 Jeff Nickoloff. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) != 2 {
		log.Fatal(`requires one argument - relative name of the file to be included in the tar stream`)
	}
	fn := os.Args[1]

	fi, err := os.Stat(fn)
	if err != nil {
		log.Fatal(err)
	}

	if fi.Mode() & os.ModeNamedPipe == 0 {
		log.Fatal(`Target file is not a named pipe.`)
	}

	type payload struct {
		Header *tar.Header
		Data *bytes.Buffer
	}

	// start stdin reader
	// STDIN is a tar stream, meaning that a payload is not ready for 
	// writing until the whole record has been retrieved.
	stin := make(chan payload)
	go func() {
		defer os.Stdin.Close()
		defer close(stin)
		r := bufio.NewReader(os.Stdin)
		tin := tar.NewReader(r)
		for {
			h, err := tin.Next()
			if err != nil && err == io.EOF {
				return
			} else if err != nil {
				log.Fatalln(err)
			}
			var b bytes.Buffer
			if _, err := io.Copy(&b, tin); err != nil {
				log.Fatal(err)
			}
			stin <- payload{
				Header: h,
				Data: &b,
			}
		}
	}()

	// start pipe reader
	ptin := make(chan payload)
	go func() {
		f, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defer close(ptin)
		r := bufio.NewReader(f)
		tin := tar.NewReader(r)
		for {
			h, err := tin.Next()
			if err != nil && err == io.EOF {
				return
			} else if err != nil {
				log.Fatalln(err)
			}
			var b bytes.Buffer
			if _, err := io.Copy(&b, tin); err != nil {
				log.Fatal(err)
			}
			ptin <- payload{
				Header: h,
				Data: &b,
			}
		}
	}()

	// tar writer loop
	tw := tar.NewWriter(os.Stdout)
	var p payload
	var ok bool
	for {
		if stin == nil && ptin == nil {
			break
		}
		select {
			case p, ok = <- stin:
				if !ok {
					stin = nil
					continue
				}
			case p, ok = <- ptin:
				if !ok {
					ptin = nil
					continue
				}
		}

		if err := tw.WriteHeader(p.Header); err != nil {
			log.Fatalln(err)
		}
		if _, err := io.Copy(tw, p.Data); err != nil {
			log.Fatal(err)
		}
	}

	// close the tar stream on stdout
	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}
}

