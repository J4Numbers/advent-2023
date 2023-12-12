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

func breakItemDescriptions(itemDesc []string) []ItemLog {
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
    var tmpLog ItemLog
    tmpLog.logLine = splitLn[logLineIndex]
    tmpLog.nonogram = nonogramIn
    itemLogs = append(itemLogs, tmpLog)
  }
  return itemLogs
}

func addLetter(c chan string, combo string, alphabet string, reqLen int) {
  if reqLen <= 0 {
    c <- combo
    return
  }

  for _, ch := range alphabet {
    addLetter(c, combo + string(ch), alphabet, reqLen - 1)
  }
}

func explodeUnknownPossibilities(reqLen int) <-chan string {
  c := make(chan string)

  go func(c chan string) {
    defer close(c)
    addLetter(c, "", ".#", reqLen)
  }(c)

  return c
}

func findValidNonograms(itemLog ItemLog) int {
  var intervals = ``
  for idx, nonoLen := range itemLog.nonogram {
    if idx > 0 {
      intervals += `\.+`
    }
    intervals += fmt.Sprintf(`#{%v}`, nonoLen)
  }
  intervals = fmt.Sprintf(`^\.*%v\.*$`, intervals)
  validStrRegex := regexp.MustCompile(intervals)
  unknownLocsRegex := regexp.MustCompile(`\?+`)

  unknownLocs := unknownLocsRegex.FindAllStringIndex(itemLog.logLine, -1)
  var unknownLens []int
  var spaceLen = 0
  for _, unknownLoc := range unknownLocs {
    unknownLens = append(unknownLens, unknownLoc[1] - unknownLoc[0])
    spaceLen += unknownLoc[1] - unknownLoc[0]
  }

  var count = 0
  for possibility := range explodeUnknownPossibilities(spaceLen) {
    var line =  itemLog.logLine
    var startIdx = 0
    for _, uLen := range unknownLens {
      line = strings.Replace(line, strings.Repeat("?", uLen), possibility[startIdx:startIdx+uLen], 1)
      startIdx += uLen
    }
    debugLine(fmt.Sprintf("Comparing candidate %v against regex %v", line, intervals))
    if validStrRegex.MatchString(line) {
      debugLine(fmt.Sprintf("Successfully discovered candidate %v", line))
      count++
    }
  }
  //fmt.Println(hashLocs)
  //fmt.Println(unknownLocs)
  return count
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
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
  itemLog := breakItemDescriptions(fileContents)
  debugLine(fmt.Sprintf("%v", itemLog))

  for _, item := range itemLog {
    approaches += findValidNonograms(item)
  }

  fmt.Println(approaches)
}
