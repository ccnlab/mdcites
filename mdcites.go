// Copyright (c) 2020, The CCNLab Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mdcites extracts markdown citations in the format [@Ref] from .md files and
// looks those up in a specified .bib file, and writes the refs to target .bib
// file which can then be used by pandoc-citeproc to efficiently process
// references.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goki/ki/dirs"
	"github.com/goki/pi/langs/bibtex"
)

func main() {
	var srcDir string
	var srcBib string
	var outBib string
	flag.StringVar(&srcDir, "dir", "./", "optional directory containing .md files to process")
	flag.StringVar(&srcBib, "bib", "", "required full path to .bib file containing all references that could be cited")
	flag.StringVar(&outBib, "out", "references.bib", "filename for output .bib file containing only the cited references")
	flag.Parse()

	if srcBib == "" {
		flag.Usage()
		return
	}

	exp := regexp.MustCompile(`\[(@([[:alnum:]]+-?)+(;[[:blank:]]+)?)+\]`)

	mds := dirs.ExtFileNames(srcDir, []string{".md"})
	if len(mds) == 0 {
		fmt.Printf("No .md files found in: %v\n", srcDir)
		os.Exit(1)
	}

	bf, err := os.Open(srcBib)
	if err != nil {
		fmt.Println(err)
		return
	}
	parsed, err := bibtex.Parse(bf)
	if err != nil {
		bf.Close()
		fmt.Printf("Bibtex bibliography: %s not loaded due to error(s):\n", srcBib)
		fmt.Println(err)
		os.Exit(1)
	}
	bf.Close()

	refs := map[string]int{}

	for _, md := range mds {
		fn := filepath.Join(srcDir, md)
		fmt.Printf("processing: %v\n", fn)
		f, err := os.Open(fn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		scan := bufio.NewScanner(f)
		for scan.Scan() {
			cs := exp.FindAllString(string(scan.Bytes()), -1)
			for _, c := range cs {
				tc := c[1 : len(c)-1]
				sp := strings.Split(tc, "@")
				for _, ac := range sp {
					a := strings.TrimSpace(ac)
					a = strings.TrimSuffix(a, ";")
					if a == "" {
						continue
					}
					cc, _ := refs[a]
					cc++
					refs[a] = cc
				}
			}
		}
		f.Close()
	}
	fmt.Printf("cites:\n%v\n", refs)

	ob := bibtex.NewBibTex()
	ob.Preambles = parsed.Preambles
	ob.StringVar = parsed.StringVar

	for r, _ := range refs {
		be, has := parsed.Lookup(r)
		if has {
			ob.Entries = append(ob.Entries, be)
		} else {
			fmt.Printf("Error: Reference key: %v not found in %s\n", r, srcBib)
		}
	}
	out := ob.PrettyString()

	of, err := os.Create(outBib)
	if err != nil {
		fmt.Println(err)
		return
	}
	of.WriteString(out)
	of.Close()
}
