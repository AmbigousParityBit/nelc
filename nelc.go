package main

import (
	"AmbigousParityBit/tools/errors"
	"AmbigousParityBit/tools/linux/lnxsmp"
	"AmbigousParityBit/tools/logext"
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	docopt "github.com/docopt/docopt-go"
)

var (
	version = "1.3"
	commit  string
	branch  string
)

const usageInfo = `:: nelc :: (c) AmbigousParityBit ::
display information about count of non empty lines in text files in given directory, searched recursively

Usage:
	nelc 
	nelc <directory>
	nelc -h | --help | --version

Options:
  -h, --help  	Show this screen.
  --version  	Show version.`

type countLinesType struct {
	empty, nonEmpty int
}

type countLinesTypeMap map[string]countLinesType

func init() {
	logext.Default("nelc")
}

func printInfo(extensions *countLinesTypeMap) {
	fmt.Printf(":: nelc %s :: (c) AmbigousParityBit ::\n:: information about count of non empty lines in text files in given directory, searched recursively\n", version)
	maxCountSize, maxExtSize, items, percent := 0, 0, 1, 0.
	total, totalNonEmpty := 0, 0
	for k, v := range *extensions {
		if len(k) > maxExtSize {
			maxExtSize = len(k)
		}
		if maxCountSize < v.empty {
			maxCountSize = v.empty
		}
		if maxCountSize < v.nonEmpty {
			maxCountSize = v.nonEmpty
		}
	}
	maxExtSize++
	maxCountSize = int(math.Floor(math.Log10(float64(maxCountSize)) + 1))
	colsDiv := int(math.Floor(80./float64(maxExtSize+maxCountSize+8) + 1))
	formatString := fmt.Sprintf("%s%v%s%v%s", "%", maxExtSize, "s %", maxCountSize, "d %3.0f%% ")

	for k, v := range *extensions {
		percent = 100. * float64(v.nonEmpty) / float64(v.nonEmpty+v.empty)
		totalNonEmpty += v.nonEmpty
		total += v.nonEmpty + v.empty
		if items%colsDiv == 0 {
			fmt.Printf(formatString+"\n", k, v.nonEmpty, percent)
		} else {
			fmt.Printf(formatString+"| ", k, v.nonEmpty, percent)
		}
		items++
	}
	fmt.Printf("\n\n:: total %d %3.0f%%\n", totalNonEmpty, 100*float64(totalNonEmpty)/float64(total))
}

func getFileContentType(filePath string) string {
	f, err := os.Open(filePath)
	errors.FatalErr(err, "opening file "+filePath+" ")
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err == io.EOF {
		return ""
	}
	errors.FatalErr(err, "reading file "+filePath+" ")
	contentType := http.DetectContentType(buffer)

	return contentType
}

func getLinesCount(filePath string) (count countLinesType) {
	f, err := os.Open(filePath)
	errors.FatalErr(err, "opening file "+filePath+" ")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 8*1024), 1024*1024*1024)
	line := ""
	count = countLinesType{0, 0}

	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			count.empty++
		} else {
			count.nonEmpty++
		}
	}
	errors.FatalErr(err, "reading file "+filePath+" ")

	return count
}

func scanDir(dir string, extensions *countLinesTypeMap) {
	contentType := ""
	ext := ""
	count := countLinesType{0, 0}
	if lnxsmp.IsFileDirExisting(dir) {
		err := filepath.Walk(dir, func(pathString string, info os.FileInfo, err error) error {
			errors.PrintErr(err, "recursively walking directory "+pathString+" ")
			if info.Name() != "." && !lnxsmp.IsDir(pathString) {
				contentType = getFileContentType(pathString)
				errors.FatalErr(err, "getting file content type "+pathString+" ")
				if strings.HasPrefix(contentType, "text/plain") {
					ext = path.Ext(info.Name())
					count = getLinesCount(pathString)
					if _, exist := (*extensions)[ext]; exist {
						count.empty += (*extensions)[ext].empty
						count.nonEmpty += (*extensions)[ext].nonEmpty
					} else {
						(*extensions)[ext] = count
					}
				}
			}
			return nil
		})
		errors.FatalErr(err, "recursively walking directory "+dir+" ")
	}
}

func parseArgs() string {
	arguments, err := docopt.Parse(usageInfo, os.Args[1:], true /*show help?*/, version /*version*/, true)
	errors.FatalErr(err, "parsing command line arguments")
	dir, _ := arguments["<directory>"].(string)
	errors.FatalErr(err, "parsing command line arguments (directory)")

	if dir == "" {
		dir = "."
	}

	return dir
}

func main() {
	dir := parseArgs()
	extensions := countLinesTypeMap{}
	scanDir(dir, &extensions)
	printInfo(&extensions)
}
