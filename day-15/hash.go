// Main package path for day 15 of the AoC 2023 challenge. This time, we're going to
// hash a series of strings within a file.
//
// This can be run with the -i input flag to change the input file. Adding the --focus
// flag sets the program to build a focus map of many boxes
package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strconv"
  "strings"
)

// debug - Choose whether to run the program in debug mode
var debug = false

// Print out a given line if debug is enabled during the runtime of this program
func debugLine(lineToDebug string) {
  if debug {
    fmt.Println(lineToDebug)
  }
}

// Read in a given file and return the list of strings to hash and focus on
func readFile(filename string) ([]string, error) {
  var hashStrings []string

  file, err  := os.Open(filename)
  if err != nil {
    return hashStrings, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    var line = scanner.Text()
    hashPotentials := strings.Split(line, ",")
    for _, hashPot := range hashPotentials {
      if strings.TrimSpace(hashPot) != "" {
        hashStrings = append(hashStrings, hashPot)
      }
    }
  }

  return hashStrings, nil
}

// LabelFocus - A small object to match a label against its given focus value
type LabelFocus struct {
  label string
  focus int
}

// Update a label in-place within the focus map if it already exists, otherwise we add
// the label to the list on the end and return that instead
func updateLabel(focus LabelFocus, focusMap []LabelFocus) []LabelFocus {
  var idx = -1
  for foundIdx, labelUnderTest := range focusMap {
    if labelUnderTest.label == focus.label {
      idx = foundIdx
      break
    }
  }

  if idx >= 0 {
    debugLine(fmt.Sprintf("Label %v already exists in box, overwriting with focus %v", focus.label, focus.focus))
    return append(append(focusMap[:idx], focus), focusMap[idx+1:]...)
  }
  debugLine(fmt.Sprintf("Label %v not yet existing in box, setting with focus %v", focus.label, focus.focus))
  return append(focusMap, focus)
}

// Attempt to remove a given label from a map of labels to focuses - if the label does not
// exist, then we do nothing.
func removeLabel(focusLabel string, focusMap []LabelFocus) []LabelFocus {
  var idx = -1
  for foundIdx, labelUnderTest := range focusMap {
    if focusLabel == labelUnderTest.label {
      idx = foundIdx
      break
    }
  }
  if idx >= 0 {
    debugLine(fmt.Sprintf("Label %v exists in box, removing...", focusLabel))
    return append(focusMap[:idx], focusMap[idx + 1:]...)
  }
  debugLine(fmt.Sprintf("Label %v does not exist in box... Doing nothing", focusLabel))
  return focusMap
}

// Generate the value of a set of boxes according to our algorithm as described in the
// readme (label hash + 1) * (focus order + 1) * focus, all added together into one
// total value.
func generateBoxCount(boxMap map[int][]LabelFocus) int {
  count := 0
  for boxNum, focusList := range boxMap {
    for orderIdx, focus := range focusList {
      count += (boxNum + 1) * (orderIdx + 1) * focus.focus
    }
  }
  return count
}

// Given a string of characters, hash that string according to a small algorithm,
// bounded to 256 options.
func hashWord(wordToHash string) int {
  hash := 0

  for _, chr := range wordToHash {
    hash = ((hash + int(chr)) * 17) % 256
  }

  return hash
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  var focus bool
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.BoolVar(&focus, "focus", false, "Enable focus mode")
  flag.BoolVar(&debug, "debug", false, "Enable debug logging")
  flag.Parse()

  // Read in the given file as a number of maps.
  hashStrings, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Start up some working variables
  var count = 0
  var focusMap = make(map[int][]LabelFocus)

  if focus {
    // If we're building a focus map, for each option, split the computing into whether
    // we're working on an assignment or a removal, then perform those actions
    // accordingly.
    for _, strToMap := range hashStrings {
      if strings.ContainsRune(strToMap, '=') {
        putInLabel := strings.Split(strToMap, "=")
        hashKey := hashWord(putInLabel[0])
        debugLine(fmt.Sprintf("Found hash value of %v for string %v", hashKey, putInLabel[0]))

        var label LabelFocus
        label.label = putInLabel[0]
        label.focus, err = strconv.Atoi(putInLabel[1])
        focusMap[hashKey] = updateLabel(label, focusMap[hashKey])
      } else if strings.ContainsRune(strToMap, '-') {
        remLabel := strings.Split(strToMap, "-")
        hashKey := hashWord(remLabel[0])
        debugLine(fmt.Sprintf("Found hash value of %v for string %v", hashKey, remLabel[0]))
        focusMap[hashKey] = removeLabel(remLabel[0], focusMap[hashKey])
      }
    }
    count = generateBoxCount(focusMap)
  } else {
    // If we're just hashing strings, hash each string together and return the result.
    for _, strToHash := range hashStrings {
      var addCount = hashWord(strToHash)
      debugLine(fmt.Sprintf("Found hash value of %v for string %v", addCount, strToHash))
      count += addCount
    }
  }

  // Return the hash/focus score of our input
  fmt.Println(count)
}
