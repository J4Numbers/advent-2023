// Main package path for day 12 of the AoC 2023 challenge. This time, we're essentially
// solving a nonogram puzzle on a series of lines.
//
// This can be run with the -i input flag to change the input file accordingly.
package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "strconv"
  "strings"
)

// ValidLineCheck - the measure of whether a line we read in is a valid puzzle input
const ValidLineCheck string = `^(?P<LogLine>[.#?]+)\s+(?P<LogDesc>[0-9,]+)$`

// debug - Choose whether to run the program in debug mode
var debug = false

// Print out a given line if debug is enabled during the runtime of this program
func debugLine(lineToDebug string) {
  if debug {
    fmt.Println(lineToDebug)
  }
}

type ItemLog struct {
  logLine string
  nonogram []int
}

// Read in a given file and return each line in a slice
func readFile(filename string) ([]string, error) {
  var fileContents []string

  file, err  := os.Open(filename)
  if err != nil {
    return fileContents, err
  }
  defer file.Close()

  var lineRegex = regexp.MustCompile(ValidLineCheck)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    var line = scanner.Text()
    if lineRegex.MatchString(line) {
      fileContents = append(fileContents, line)
    }
  }
  return fileContents, nil
}

func breakItemDescriptions(itemDesc []string, folds int) []ItemLog {
  var itemLogs []ItemLog
  lineSplitter := regexp.MustCompile(ValidLineCheck)
  logLineIndex := lineSplitter.SubexpIndex("LogLine")
  breakSplit := regexp.MustCompile(`[0-9]+`)
  for i := 0; i < len(itemDesc); i++ {
    var nonogramIn []int
    splitLn := lineSplitter.FindStringSubmatch(itemDesc[i])
    puzzleVals := breakSplit.FindAllString(itemDesc[i], -1)
    for puzzleIdx := 0; puzzleIdx < len(puzzleVals); puzzleIdx++ {
      convVal, convErr := strconv.Atoi(puzzleVals[puzzleIdx])
      if convErr != nil {
        fmt.Println(convErr)
        break
      }
      nonogramIn = append(nonogramIn, convVal)
    }

    var logLine string
    var nonogramOut []int
    for j := 0; j < folds; j++ {
      if j > 0 {
        logLine += "?"
      }
      logLine += splitLn[logLineIndex]
      nonogramOut = append(nonogramOut, nonogramIn...)
    }

    var tmpLog ItemLog
    tmpLog.logLine = logLine
    tmpLog.nonogram = nonogramOut
    itemLogs = append(itemLogs, tmpLog)
  }
  return itemLogs
}

func generateTestRegex(input string) []*regexp.Regexp {
  var regexCache []*regexp.Regexp
  for i := 0; i <= len(input); i++ {
    var logLinePoss = ""
    for j := 0; j < i; j++ {
      if input[j] == '?' {
        logLinePoss += "[.#]"
      } else {
        logLinePoss += fmt.Sprintf("\\%v", string(input[j]))
      }
    }
    logLinePoss = fmt.Sprintf("^%v$", logLinePoss)
    regexCache = append(regexCache, regexp.MustCompile(logLinePoss))
  }
  return regexCache
}

func addLetter(c chan string, combo string, nonogram []int, cache []*regexp.Regexp, minLen int, reqLen int) {
  if reqLen < minLen || !cache[len(combo)].MatchString(combo) {
    return
  }
  if reqLen <= 0 {
    c <- combo
    return
  }

  if len(nonogram) > 0 {
    var nextCombo = strings.Repeat("#", nonogram[0])
    if reqLen > nonogram[0] {
      nextCombo += "."
    }

    addLetter(
      c, combo + nextCombo, nonogram[1:],
      cache, minLen - len(nextCombo), reqLen - len(nextCombo))
  }
  addLetter(c, combo + string('.'), nonogram, cache, minLen, reqLen - 1)
}

func explodeUnknownPossibilities(nonogram []int, inputLn string, reqLen int) <-chan string {
  c := make(chan string)
  minLen := 0
  for i := 0; i < len(nonogram); i++ {
    minLen += nonogram[i]
  }

  regexCache := generateTestRegex(inputLn)

  go func(c chan string) {
    defer close(c)
    addLetter(c, "", nonogram, regexCache, minLen + len(nonogram) - 1, reqLen)
  }(c)

  return c
}

func findValidNonograms(itemLog ItemLog) int {
  var intervals = ``
  var brokenWidth = 0
  for idx, nonoLen := range itemLog.nonogram {
    if idx > 0 {
      intervals += `\.+`
    }
    intervals += fmt.Sprintf(`#{%v}`, nonoLen)
    brokenWidth += nonoLen
  }
  intervals = fmt.Sprintf(`^\.*%v\.*$`, intervals)
  validStrRegex := regexp.MustCompile(intervals)

  var count = 0
  for possibility := range explodeUnknownPossibilities(itemLog.nonogram, itemLog.logLine, len(itemLog.logLine)) {
    debugLine(fmt.Sprintf("Comparing candidate %v against regex %v", possibility, intervals))
    if validStrRegex.MatchString(possibility) {
      debugLine(fmt.Sprintf("Successfully discovered candidate %v", possibility))
      count++
    }
  }
  return count
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  var folds int
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.IntVar(&folds, "f", 1, "Number of times to repeat a given item line")
  flag.BoolVar(&debug, "debug", false, "Enable debug logging")
  flag.Parse()

  // Read in the given file as a number of lines.
  fileContents, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  var approaches = 0

  // Break each line into its line description and the broken spring lengths
  itemLog := breakItemDescriptions(fileContents, folds)
  debugLine(fmt.Sprintf("%v", itemLog))

  for _, item := range itemLog {
    approaches += findValidNonograms(item)
  }

  fmt.Println(approaches)
}
