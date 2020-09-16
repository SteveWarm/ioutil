package ioutil

import (
	"io"
	"os"
	"regexp"
	"bufio"
)

func GrepFromReader(r io.Reader, regstr string) (lines []string, err error) {
	regex, e := regexp.Compile(regstr)
	if e != nil {
		err = e
		return
	}

	reader := bufio.NewReader(r)
	var line string
	for {
		d, isprefix, e := reader.ReadLine()
		if e != nil && e != io.EOF {
			err = e
			return
		}

		line += string(d)

		if isprefix {
			if e == io.EOF {
				break
			}
			continue
		}

		matchs := regex.FindStringSubmatch(line)
		for _, s := range matchs {
			lines = append(lines, s)
		}

		line = ""

		if e == io.EOF {
			break
		}
	}

	return
}

func GrepFromFile(filestr, regstr string) (lines []string, err error) {
	f, e := os.OpenFile(filestr, os.O_RDONLY, os.FileMode(0644))
	if e != nil {
		err = e
		return
	}
	defer f.Close()

	return GrepFromReader(f, regstr)
}
