// Main package path for day 11 of the AoC 2023 challenge. This time, we're looking
// to calculate the steps between several points in a graph which is potentially folded
// down given rows and columns.
//
// This can be run with the -i input flag and the -e expansion rate flag to change the
// input file and rate of expansion accordingly.
package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "slices"
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

// Given a set of lines which are in our space format, return two slices containing
// the X columns and Y rows that do not contain any galaxies
func findEmptySpace(fileContents []string) ([]int, []int) {
  var xLines []int
  var yLines []int
  var lineLen = -1

  var hashRegex = regexp.MustCompile(`#`)

  for lineIdx := range fileContents {
    // If we don't know the line length yet, we can find this out on the first
    // iteration and pre-fill the yLines array with every possible position as
    // we do not yet know which lines are empty
    if lineLen == -1 {
      lineLen = len(fileContents[lineIdx])
      for i := 0; i < len(fileContents[lineIdx]); i++ {
        yLines = append(yLines, i)
      }
    }

    // If the line contains any hashes at all, we can remove their positions from
    // the Y lines as we have found proof that they exist.
    if hashRegex.MatchString(fileContents[lineIdx]) {
      var locations = hashRegex.FindAllStringIndex(fileContents[lineIdx], -1)
      for locIdx := range locations {
        for idx, loc := range yLines {
          if loc == locations[locIdx][0] {
            yLines = slices.Delete(yLines, idx, idx + 1)
          }
        }
      }
    } else {
      // Otherwise, the line contains no hashes and is an empty X row.
      xLines = append(xLines, lineIdx)
    }
  }

  return xLines, yLines
}

// Given a map of space which is several lines containing .s and #s, along with the
// desired expansion rate, the slice of empty X rows and Y columns, we return a map
// of all galaxies within the map, along with their expanded positions.
func extractGalaxies(spaceMap []string, expansionRate int, xLines []int, yLines []int) []map[string]int {
  var galaxyList []map[string]int
  var hashRegex = regexp.MustCompile(`#`)

  for i := 0; i < len(spaceMap); i++ {
    // Find the positions of all hashes within a given line
    var locations = hashRegex.FindAllStringIndex(spaceMap[i], -1)
    for locIdx := range locations {
      // And calculate the modified X location, which is the position within the regex
      // search, + x(expansionRate), where x is the number of Y columns that the hash
      // appears after that contain empty space.
      xLoc := locations[locIdx][0]
      for _, yLine := range yLines {
        if yLine < locations[locIdx][0] {
          xLoc += expansionRate
        }
      }

      // Calculate the modified Y location, which is the current row + x(expansionRate),
      // where x is the number of X rows that the # appears after that contain empty space.
      yLoc := i
      for _, xLine := range xLines {
        if xLine < i {
          yLoc += expansionRate
        }
      }

      // Create and fill in the map representation of the galaxy and add it to the list.
      tmpMap := make(map[string]int)
      tmpMap["x"] = xLoc
      tmpMap["y"] = yLoc
      galaxyList = append(galaxyList, tmpMap)
    }
  }
  return galaxyList
}

// Given a list of galaxies (which contain and x and y coordinate), calculate the number of
// steps between every galaxy. The number of galaxy paths that need checking correlates
// exponentially with the number of galaxies - where 2 galaxies return 1 path, 3 generates
// 3 paths, 4 generates 6, and so on.
func calculateStepsBetweenGalaxies(galaxyList []map[string]int) int {
  stepCount := 0

  // We do an inverse for loop to avoid replaying any step counts that we have already
  // calculated.
  for i := 0; i < len(galaxyList); i++ {
    for j := i + 1; j < len(galaxyList); j++ {
      // To calculate the required steps between any point is the absolute value of the
      // difference between two points on both axis added together.
      xSteps := galaxyList[i]["x"] - galaxyList[j]["x"]
      if xSteps < 0 {
        xSteps = -xSteps
      }
      ySteps := galaxyList[i]["y"] - galaxyList[j]["y"]
      if ySteps < 0 {
        ySteps = -ySteps
      }
      stepCount += xSteps + ySteps
    }
  }
  return stepCount
}

// Main function to kick the work
func main() {
  // Do some initial CLI parsing to figure out what the requested operation is.
  var filename string
  var expansionRate int
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.IntVar(&expansionRate, "e", 1, "Specify the rate of expansion")
  flag.BoolVar(&debug, "debug", false, "Enable debug logging")
  flag.Parse()

  // Read in the given file as a number of lines.
  fileContents, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Calculate the empty rows and columns within the file and debug that out to
  // the console.
  xLines, yLines := findEmptySpace(fileContents)
  debugLine(fmt.Sprintf("Found empty X lines at %v", xLines))
  debugLine(fmt.Sprintf("Found empty Y lines at %v", yLines))

  // Find the galaxies that are contained within the space map along with and track
  // them with their expanded coordinates, then print them as debug.
  galaxies := extractGalaxies(fileContents, expansionRate, xLines, yLines)
  debugLine(fmt.Sprintf("Found galaxies %v", galaxies))

  // Calculate the steps between all galaxies, then return that as output.
  steps := calculateStepsBetweenGalaxies(galaxies)
  fmt.Println(steps)
}
