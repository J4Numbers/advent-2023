// Main package path for day 13 of the AoC 2023 challenge. This time, we're looking
// to find a reflective line in any graph we are presented with.
//
// This can be run with the -i input flag to change the input file, and the -d flag
// to change the number of set differences in each mirror map
package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
)

// ValidLineCheck - the measure of whether a line we read in is valid or not
const ValidLineCheck string = `^[.#]+$`

// debug - Choose whether to run the program in debug mode
var debug = false

// Print out a given line if debug is enabled during the runtime of this program
func debugLine(lineToDebug string) {
  if debug {
    fmt.Println(lineToDebug)
  }
}

// Read in a given file and return each line in a slice
func readFile(filename string) ([][]string, error) {
  var fileContents [][]string

  file, err  := os.Open(filename)
  if err != nil {
    return fileContents, err
  }
  defer file.Close()

  var lineRegex = regexp.MustCompile(ValidLineCheck)

  scanner := bufio.NewScanner(file)
  var newMap []string

  for scanner.Scan() {
    var line = scanner.Text()
    if lineRegex.MatchString(line) {
      newMap = append(newMap, line)
    } else {
      if len(newMap) > 0 {
        fileContents = append(fileContents, newMap)
        newMap = []string{}
      }
    }
  }

  if len(newMap) > 0 {
    fileContents = append(fileContents, newMap)
  }

  return fileContents, nil
}

// DifferenceLog - A small object to track a candidate axis row/column and the number
// of differences we have found for that candidate so far
type DifferenceLog struct {
  candidate int
  differenceCount int
}

// Given two input strings, calculate the number of characters that are different
// between the two. We always expect both input strings to be of equal length.
func calculateDifferenceCount(inputA string, inputB string) int {
  var diffCount = 0
  for strIdx := 0; strIdx < len(inputA); strIdx++ {
    if inputA[strIdx] != inputB[strIdx] {
      diffCount += 1
    }
  }
  return diffCount
}

// Given an array of strings representing a map of mirrors and the number of allowed
// differences in a reflection, return the best candidate (if it exists) for which
// we have a reflection along the vertical (|) axis
func findVerticalReflection(mirrorMap []string, allowedDifferences int) int {
  var candidates []DifferenceLog
  var lineLen = -1

  for _, mirrorLine := range mirrorMap {
    // Set up the candidates as prefill since any column could be a candidate at
    // this point. We give them all a starting difference value of 0 too
    if lineLen < 0 {
      lineLen = len(mirrorLine)
      for fillIdx := 1; fillIdx < len(mirrorLine); fillIdx++ {
        var diffLog DifferenceLog
        diffLog.candidate = fillIdx
        diffLog.differenceCount = 0
        candidates = append(candidates, diffLog)
      }
    }

    debugLine(fmt.Sprintf("Inspecting %v with candidates %v", mirrorLine, candidates))

    // For each column, get a string up to that point and reflect it backwards, then
    // compare it against a string of equal length taken from the other side of the
    // candidate axis.
    for candIdx := len(candidates); candIdx > 0; candIdx-- {
      lSplit := mirrorLine[0:candidates[candIdx-1].candidate]
      var revLSplit = ""
      for lIdx := len(lSplit); lIdx > 0; lIdx-- {
        revLSplit += string(lSplit[lIdx - 1])
      }
      var rSplit string
      if candidates[candIdx-1].candidate * 2 > len(mirrorLine) {
        rSplit = mirrorLine[candidates[candIdx-1].candidate:]
        revLSplit = revLSplit[:len(rSplit)]
      } else {
        rSplit = mirrorLine[candidates[candIdx-1].candidate:candidates[candIdx-1].candidate * 2]
      }

      debugLine(fmt.Sprintf("Comparing reversed lSplit %v against rSplit %v", revLSplit, rSplit))
      if revLSplit != rSplit {
        // If the two strings were different, calculate the difference count and
        // add that to the stored candidate difference. If that ever exceeds our
        // allowed values, we remove that candidate from the array
        var diffCount = calculateDifferenceCount(revLSplit, rSplit)
        candidates[candIdx-1].differenceCount += diffCount
        debugLine(fmt.Sprintf("Found %v differences to a total of %v", diffCount, candidates[candIdx-1].differenceCount))
        if candidates[candIdx-1].differenceCount > allowedDifferences {
          candidates = append(candidates[:candIdx-1], candidates[candIdx:]...)
        }
      }
    }
  }

  // Find a candidate which matches the allowed differences since that is a match
  // rather than an allowance factor.
  if len(candidates) > 0 {
    for _, cand := range candidates {
      if cand.differenceCount == allowedDifferences {
        return cand.candidate
      }
    }
  }
  return 0
}

func findHorizontalReflection(mirrorMap []string, allowedDifferences int) int {
  // Set up the candidates as prefill since any row could be a candidate at
  // this point. We give them all a starting difference value of 0 too
  var candidates []DifferenceLog
  for candIdx := 1; candIdx < len(mirrorMap); candIdx++ {
    var diffLog DifferenceLog
    diffLog.candidate = candIdx
    diffLog.differenceCount = 0
    candidates = append(candidates, diffLog)
  }

  // For each candidate, we work from the top down and compare two adjacent rows,
  // then if those were the same, the two surround rows, then check for sameness,
  // then check the surrounding rows again until we reach the edge of the map.
  for candIdx := len(candidates); candIdx > 0; candIdx-- {
    debugLine(fmt.Sprintf("Inspecting horizontal candidate %v", candidates[candIdx-1]))
    // Reverse is going back up the array
    for revVertIdx := candidates[candIdx-1].candidate-1; revVertIdx >= 0; revVertIdx-- {
      // Forward is peeking forwards down the rest of the array in parallel to the
      // reverse
      var forwardIdx = candidates[candIdx-1].candidate + (candidates[candIdx-1].candidate - revVertIdx) - 1
      debugLine(fmt.Sprintf("Comparing forward %v against reverse %v", forwardIdx, revVertIdx))
      if forwardIdx < len(mirrorMap) {
        if mirrorMap[forwardIdx] != mirrorMap[revVertIdx] {
          // If the two strings were different, calculate the difference count and
          // add that to the stored candidate difference. If that ever exceeds our
          // allowed values, we remove that candidate from the array
          var diffCount = calculateDifferenceCount(mirrorMap[forwardIdx], mirrorMap[revVertIdx])
          candidates[candIdx-1].differenceCount += diffCount
          debugLine(fmt.Sprintf("Found %v differences to a total of %v", diffCount, candidates[candIdx-1].differenceCount))
          if allowedDifferences < candidates[candIdx-1].differenceCount {
            candidates = append(candidates[:candIdx-1], candidates[candIdx:]...)
            break
          }
        }
      }
    }
  }

  // Find a candidate which matches the allowed differences since that is a match
  // rather than an allowance factor.
  if len(candidates) > 0 {
    for _, cand := range candidates {
      if cand.differenceCount == allowedDifferences {
        return cand.candidate
      }
    }
  }
  return 0
}

// Given an array of strings representing a map and the number of allowed differences
// in each dimension, return a number which is either the number of rows to the left
// of a vertical reflection, or the number of rows above a horizontal reflection
// multiplied by 100
func findReflectionScore(mirrorMap []string, allowedDifferences int) int {
  score := 0

  score += findVerticalReflection(mirrorMap, allowedDifferences)
  score += findHorizontalReflection(mirrorMap, allowedDifferences) * 100

  return score
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  var allowedDifferences int
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.IntVar(&allowedDifferences, "d", 0, "Specify number of allowed differences between mirrors")
  flag.BoolVar(&debug, "debug", false, "Enable debug logging")
  flag.Parse()

  // Read in the given file as a number of maps.
  fileMaps, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  var count = 0

  // For each map, figure out the reflection score and add it to the count
  for _, fileMap := range fileMaps {
    var addCount = findReflectionScore(fileMap, allowedDifferences)
    debugLine(fmt.Sprintf("%v -> Found reflection score of %v", fileMap, addCount))
    count += addCount
  }

  // Return the reflection score of all of our maps
  fmt.Println(count)
}
