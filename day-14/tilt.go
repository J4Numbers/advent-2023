// Main package path for day 14 of the AoC 2023 challenge. This time, we're tilting a
// map to find out the weight on a given point.
//
// This can be run with the -i input flag to change the input file. To change the program
// to perform cycles, add in the --cycle flag, along with the --count flag to define the
// number of cycles to perform in a given run.
package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "slices"
  "sort"
)

// ValidLineCheck - the measure of whether a line we read in is valid or not
const ValidLineCheck string = `^[.#O]+$`

// debug - Choose whether to run the program in debug mode
var debug = false

// Print out a given line if debug is enabled during the runtime of this program
func debugLine(lineToDebug string) {
  if debug {
    fmt.Println(lineToDebug)
  }
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

// CacheMap - A data object which is used to show the position of a series of rocks after
// a given number of cycles.
type CacheMap struct {
  rockMapped [][]int
  cyclesMapped int
}

// Given a list of regex index searches (i.e. [[0, 1]]), extract the start index from all
// of them and return a list of all those start indexes.
func extractStartIndex(indexMap [][]int) []int {
  var indexList []int
  for _, searchIndex := range indexMap {
    indexList = append(indexList, searchIndex[0])
  }
  return indexList
}

// Given a base map of '.'s and '#'s, and the 2D array which describes the locations of
// rocks across that file, return a filled map in a string array which displays the
// locations of the movable rocks on top of the base map.
func joinOnTiltRocks(baseMap [][]rune, tiltRocks [][]int) []string {
  var retMap []string
  for lineIdx, line := range baseMap {
    var retLine = ""
    for _, chr := range line {
      retLine += string(chr)
    }
    for _, rockIdx := range tiltRocks[lineIdx] {
      retLine = retLine[:rockIdx] + "O" + retLine[rockIdx + 1:]
    }
    retMap = append(retMap, retLine)
  }
  return retMap
}

// Given a 2D map of rocks, generate a hash key that we're going to use for caching.
// The hash algorithm isn't particularly complex, and is the line index multiplied by
// each individual rock, all added together.
func generateCacheKey(rockMap [][]int) uint64 {
  var hashKey uint64 = 0
  for idx, rockLine := range rockMap {
    for _, rock := range rockLine {
      hashKey += uint64((idx + 1) * rock)
    }
  }
  return hashKey
}

// Perform a single tilt on a base map with a given set of rocks in a given roll direction
// (which can be N, W, S, and E), and return the modified position of those rocks after
// being rolled in a given direction on a base map.
func tilt(baseMap [][]rune, rockMap [][]int, rollDir rune) [][]int {
  var newRocks [][]int

  for i := 0; i < len(rockMap); i++ {
    newRocks = append(newRocks, []int{})
  }

  if rollDir == 'N' {
    for lineIdx, rockLine := range rockMap {
      for _, rockPos := range rockLine {
        var insertIdx = lineIdx
        for rollIdx := lineIdx - 1; rollIdx >= 0; rollIdx-- {
          if baseMap[rollIdx][rockPos] != '.' || slices.Contains(newRocks[rollIdx], rockPos) {
            break
          } else {
            insertIdx -= 1
          }
        }
        debugLine(fmt.Sprintf("Moving rock from [%v, %v] N to [%v, %v]", rockPos, lineIdx, rockPos, insertIdx))
        newRocks[insertIdx] = append(newRocks[insertIdx], rockPos)
      }
    }
  } else if rollDir == 'W' {
    for lineIdx, rockLine := range rockMap {
      for _, rockPos := range rockLine {
        var insertIdx = rockPos
        for rollIdx := rockPos - 1; rollIdx >= 0; rollIdx-- {
          if baseMap[lineIdx][rollIdx] != '.' || slices.Contains(newRocks[lineIdx], rollIdx) {
            break
          } else {
            insertIdx--
          }
        }
        debugLine(fmt.Sprintf("Moving rock from [%v, %v] W to [%v, %v]", rockPos, lineIdx, insertIdx, lineIdx))
        newRocks[lineIdx] = append(newRocks[lineIdx], insertIdx)
      }
    }
  } else if rollDir == 'S' {
    for lineIdx := len(rockMap) - 1; lineIdx >= 0; lineIdx-- {
      for _, rockPos := range rockMap[lineIdx] {
        var insertIdx = lineIdx
        for rollIdx := lineIdx + 1; rollIdx < len(baseMap); rollIdx++ {
          if baseMap[rollIdx][rockPos] != '.' || slices.Contains(newRocks[rollIdx], rockPos) {
            break
          } else {
            insertIdx += 1
          }
        }
        debugLine(fmt.Sprintf("Moving rock from [%v, %v] S to [%v, %v]", rockPos, lineIdx, rockPos, insertIdx))
        newRocks[insertIdx] = append(newRocks[insertIdx], rockPos)
      }
    }
  } else if rollDir == 'E' {
    for lineIdx, rockLine := range rockMap {
      for rockIdx := len(rockLine) - 1; rockIdx >= 0; rockIdx-- {
        var insertIdx = rockLine[rockIdx]
        for rollIdx := rockLine[rockIdx] + 1; rollIdx < len(baseMap[lineIdx]); rollIdx++ {
          if baseMap[lineIdx][rollIdx] != '.' || slices.Contains(newRocks[lineIdx], rollIdx) {
            break
          } else {
            insertIdx++
          }
        }
        debugLine(fmt.Sprintf("Moving rock from [%v, %v] E to [%v, %v]", rockLine[rockIdx], lineIdx, insertIdx, lineIdx))
        newRocks[lineIdx] = append(newRocks[lineIdx], insertIdx)
      }
    }
  }

  // We want to make sure all rocks are in order in the rock map so we can do tilts
  // in linear directions from L -> R or vice versa consistently.
  for rockIdx := range newRocks {
    sort.Ints(newRocks[rockIdx])
  }

  return newRocks
}

// Perform a full cycle of tilts, N, W, S, and E in that order. We do this a given cycleCount
// number of times, which is backed against a cache of rockMap locations against transposed
// locations after a given cycle count.
func cycleTilt(baseMap [][]rune, rockMap [][]int, cache map[uint64]CacheMap, cycleCount int) [][]int {
  var tiltRocks = rockMap

  // Find out if the rock formation is already in the cache
  var cacheKey = generateCacheKey(rockMap)
  cacheVal, cacheHit := cache[cacheKey]
  if cacheHit {
    debugLine(fmt.Sprintf("Cache hit on (%v) => %v", rockMap, cacheVal))
    // If it is and it matches the number of cycles we're trying to find, return that mapping.
    // Otherwise, we should ignore the cache as it does not contain the information we want.
    if cacheVal.cyclesMapped == cycleCount {
      return cacheVal.rockMapped
    }
  }

  // If we are looking for more than one cycle, we want to search back through the cycle stack
  // until we find a cache hit, then set that as our starting point
  if cycleCount > 1 {
    tiltRocks = cycleTilt(baseMap, rockMap, cache, cycleCount - 1)
  }

  // Perform a full cycle of tilting
  tiltRocks = tilt(baseMap, tiltRocks, 'N')
  tiltRocks = tilt(baseMap, tiltRocks, 'W')
  tiltRocks = tilt(baseMap, tiltRocks, 'S')
  tiltRocks = tilt(baseMap, tiltRocks, 'E')

  // Add our new entry to the cache
  var cacheEntry CacheMap
  cacheEntry.rockMapped = tiltRocks
  cacheEntry.cyclesMapped = cycleCount
  cache[cacheKey] = cacheEntry
  debugLine(fmt.Sprintf("Cache store for (%v) := %v", rockMap, cache[cacheKey]))

  // And return the rock array that we're looking for
  return cache[cacheKey].rockMapped
}

// Convert an input map from a file (or other input location) into a 2D array representing
// the base map (of empty space and immovable rocks) and a 2D array of the movable rock
// locations.
func generateTiltMap(tiltMap []string) ([][]rune, [][]int) {
  var oRegex = regexp.MustCompile("O")
  var workingTracker [][]rune
  var rollRockTracking [][]int

  // For each line, we want to extract the location of all movable rocks and get their starting
  // index (we don't care about length because we know they are always 1-char long.
  for _, tiltLine := range tiltMap {
    var lineTracking []rune
    rollRockTracking = append(rollRockTracking, extractStartIndex(oRegex.FindAllStringIndex(tiltLine, -1)))

    // For each character in a file line, we replace any 'O' characters with empty space and
    // otherwise put that character into the base map.
    for _, tiltChr := range tiltLine {
      if tiltChr == '#' || tiltChr == '.' {
        lineTracking = append(lineTracking, tiltChr)
      } else if tiltChr == 'O' {
        lineTracking = append(lineTracking, '.')
      }
    }
    workingTracker = append(workingTracker, lineTracking)
  }

  return workingTracker, rollRockTracking
}

// Calculate the weight against the north of the map when given a 2D array of rock
// positions that would otherwise fit inside a map. The closer we are to y=0, the
// more each rock counts for in the map - which is linear against the length of the
// map.
func calculateNWeight(rockMap [][]int) int {
  weight := 0
  for lineIdx, rockLine := range rockMap {
    weight += (len(rockMap) - lineIdx) * len(rockLine)
  }
  return weight
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  var enableCycle bool
  var cycleCount int
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.BoolVar(&enableCycle, "cycle", false, "Enable spin cycles for tilt maps")
  flag.IntVar(&cycleCount, "count", 1, "The number of spins to do in one run")
  flag.BoolVar(&debug, "debug", false, "Enable debug logging")
  flag.Parse()

  // Read in the given file as a number of maps.
  fileMap, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Split out the map into the base and a set of rocks that exist across several
  // rows
  baseMap, rockMap := generateTiltMap(fileMap)

  if enableCycle {
    // If we're doing cycle, we'll need to cache things because this is going to be
    // a long time running with any serious number of cycles.
    cache := make(map[uint64]CacheMap)

    // While we have cycles remaining to scan
    remCycles := cycleCount
    for remCycles > 0 {
      var readCycles = 1
      discoveredRocks, exists := cache[generateCacheKey(rockMap)]
      if exists {
        // If we have already seen this configuration of rocks, we want to build up
        // the cache further by adding another cycle onto the scan (as long as we
        // have the space remaining to rely on the cache to cycle
        if discoveredRocks.cyclesMapped + 1 < remCycles {
          readCycles = discoveredRocks.cyclesMapped + 1
        }
      }
      // Start the cycle and run a new tilt with a given number of cycles (which only
      // increases when we have the cycles cached to some degree already
      debugLine(fmt.Sprintf("Starting cycle %v - Reading cycles %v...", cycleCount - remCycles, readCycles))
      rockMap = cycleTilt(baseMap, rockMap, cache, readCycles)
      remCycles -= readCycles
      debugLine(fmt.Sprintf("Found weigh %v after %v cycles", calculateNWeight(rockMap), cycleCount - remCycles))
    }
  } else {
    // Run a single tilt to the north as part of default running
    rockMap = tilt(baseMap, rockMap, 'N')
  }

  // Return the weight against the north of our map
  weight := calculateNWeight(rockMap)
  fmt.Println(weight)

  // Print out the map that we have at the moment if we're printing debug
  if debug {
    fillMap := joinOnTiltRocks(baseMap, rockMap)
    for _, line := range fillMap {
      fmt.Println(line)
    }
  }
}
