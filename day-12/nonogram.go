// Main package path for day 12 of the AoC 2023 challenge. This time, we're essentially
// solving a nonogram puzzle on a series of lines.
//
// This can be run with the -i input flag to change the input file accordingly.
package main

import (
  "bufio"
  "flag"
  "fmt"
  "hash/fnv"
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

func hashLogAndPuzzle(logLine string, puzzle []int) uint32 {
  puzzTotal := 0
  for _, puzzLen := range puzzle {
    puzzTotal += puzzLen
  }
  workingLine := logLine + strconv.Itoa(puzzTotal)
  h := fnv.New32a()
  h.Write([]byte(workingLine))
  return h.Sum32()
}

func countArrangements(logLine string, puzzleLayout []int, cache map[uint32]uint64) uint64 {
  var cacheKey = hashLogAndPuzzle(logLine, puzzleLayout)
  cacheVal, cacheHit := cache[cacheKey]
  if cacheHit {
    debugLine(fmt.Sprintf("Cache hit for (%v,%v) := %v", logLine, puzzleLayout, cacheVal))
    return cacheVal
  }

  if len(logLine) == 0 {
    if len(puzzleLayout) > 0 {
      return 0
    } else {
      return 1
    }
  }

  var chrPtr = logLine[0]
  if chrPtr == '.' {
    var ptrIdx int
    for ptrIdx = 0; ptrIdx < len(logLine); ptrIdx++ {
      if logLine[ptrIdx] != '.' {
        break
      }
    }
    cache[cacheKey] = countArrangements(logLine[ptrIdx:], puzzleLayout, cache)

  } else if chrPtr == '?' {
    cache[cacheKey] = countArrangements(logLine[1:], puzzleLayout, cache) + countArrangements("#" + logLine[1:], puzzleLayout, cache)

  } else if chrPtr == '#' {
    if len(puzzleLayout) > 0 {
      lenToCheck := puzzleLayout[0]
      if lenToCheck <= len(logLine) && !strings.ContainsRune(logLine[0:lenToCheck], '.') {
        if lenToCheck == len(logLine) {
          if len(puzzleLayout) > 1 {
            return 0
          } else {
            return 1
          }
        }
        if logLine[lenToCheck] == '#' {
          return 0
        }
        cache[cacheKey] = countArrangements(logLine[lenToCheck+1:], puzzleLayout[1:], cache)
      } else {
        return 0
      }
    } else {
      return 0
    }
  }

  debugLine(fmt.Sprintf("Cache store for (%v,%v) := %v", logLine, puzzleLayout, cache[cacheKey]))
  return cache[cacheKey]
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

  var approaches uint64 = 0

  // Break each line into its line description and the broken spring lengths
  itemLog := breakItemDescriptions(fileContents, folds)
  debugLine(fmt.Sprintf("%v", itemLog))

  for _, item := range itemLog {
    var cache = make(map[uint32]uint64)
    var apprCount = countArrangements(item.logLine, item.nonogram, cache)
    debugLine(fmt.Sprintf("%v found %v approaches", item, apprCount))
    approaches += apprCount
  }

  fmt.Println(approaches)
}
