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

func generateTestRegex(input []int) *regexp.Regexp {
  var logLinePoss = ``
  for i := 0; i < len(input); i++ {
    logLinePoss += fmt.Sprintf(`#{%v}`, input[i])
    if i < len(input) - 1 {
      logLinePoss += `\.+`
    }
  }
  return regexp.MustCompile(fmt.Sprintf(`^\.*%v\.*$`, logLinePoss))
}

func addLetter(c chan string, combo string, testRegex *regexp.Regexp, nonogram []int, inputStr string, idx int) {
  if idx >= len(inputStr) {
    if testRegex.MatchString(combo) {
      c <- combo
    }
    return
  }

  if inputStr[idx] == '?' {
    addLetter(c, combo + string('.'), testRegex, nonogram, inputStr, idx + 1)
    addLetter(c, combo + string('#'), testRegex, nonogram, inputStr, idx + 1)
  } else {
    addLetter(c, combo + string(inputStr[idx]), testRegex, nonogram, inputStr, idx + 1)
  }
}

func explodeUnknownPossibilities(nonogram []int, inputLn string) <-chan string {
  c := make(chan string)

  go func(c chan string) {
    defer close(c)
    addLetter(c, "", generateTestRegex(nonogram), nonogram, inputLn, 0)
  }(c)

  return c
}

func calculateMinimumNonogramLength(nonogram []int) int {
  var minLen = 0
  for _, nonLen := range(nonogram) {
    minLen += nonLen
  }
  return minLen + len(nonogram) - 1
}

func containCalculation(c chan []ItemLog, ongoingSplit []ItemLog, splitStrings []string, nonograms []int, strIdx int, nonoFromIdx int) {
  if strIdx >= len(splitStrings) && nonoFromIdx < len(nonograms) {
    return
  }
  if nonoFromIdx >= len(nonograms) && strIdx >= len(splitStrings) {
    c <- ongoingSplit
    return
  }

  var blankLog ItemLog
  blankLog.logLine = splitStrings[strIdx]

  containCalculation(c, append(ongoingSplit, blankLog), splitStrings, nonograms, strIdx + 1, nonoFromIdx)
  for count := nonoFromIdx + 1; count <= len(nonograms); count++ {
    var splitNonos []int
    if count == len(nonograms) {
      splitNonos = nonograms[nonoFromIdx:]
    } else {
      splitNonos = nonograms[nonoFromIdx:count]
    }
    workingLen := calculateMinimumNonogramLength(splitNonos)

    if workingLen <= len(splitStrings[strIdx]) {
      var splitLog ItemLog
      splitLog.logLine = splitStrings[strIdx]
      splitLog.nonogram = splitNonos
      containCalculation(c, append(ongoingSplit, splitLog), splitStrings, nonograms, strIdx + 1, count)
    }
  }
}

func explodeCanContain(splitStrings []string, nonograms []int) <-chan []ItemLog {
  c := make(chan []ItemLog)

  var canContain []ItemLog

  go func(c chan []ItemLog) {
    defer close(c)
    containCalculation(c, canContain, splitStrings, nonograms, 0, 0)
  }(c)

  return c
}

func findValidNonograms(itemLog ItemLog) int {
  possibilityAreas := regexp.MustCompile(`[#?]+`).FindAllString(itemLog.logLine, -1)
  fmt.Println(itemLog.nonogram)
  fmt.Println(possibilityAreas)

  var count = 0

  for logItem := range explodeCanContain(possibilityAreas, itemLog.nonogram) {
    var logAppr = 1
    for _, splitLog := range logItem {
      var possCount = 0
      for poss := range explodeUnknownPossibilities(splitLog.nonogram, splitLog.logLine) {
        debugLine(fmt.Sprintf("%v -> %v -> %v", logItem, splitLog, poss))
        possCount += 1
      }
      logAppr *= possCount
    }
    count += logAppr
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
    var apprCount = findValidNonograms(item)
    debugLine(fmt.Sprintf("%v found %v approaches", item, apprCount))
    approaches += apprCount
  }

  fmt.Println(approaches)
}
