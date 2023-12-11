package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
  "slices"
  "strings"
)

const ValidLineCheck string = `^[.#]+$`

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

func explodeSpace(fileContents []string, expansionRate int) [][]rune {
  var galaxyMap [][]rune
  var verticalSpaces []int
  var horizontalStretchMap []string
  var lineLen = -1

  var hashRegex = regexp.MustCompile(`#`)

  for lineIdx := range fileContents {
    if lineLen == -1 {
      lineLen = len(fileContents[lineIdx])
      for i := 1; i <= len(fileContents[lineIdx]); i++ {
        verticalSpaces = append(verticalSpaces, i)
      }
    }
    if hashRegex.MatchString(fileContents[lineIdx]) {
      var locations = hashRegex.FindAllStringIndex(fileContents[lineIdx], -1)
      for locIdx := range locations {
        for idx, loc := range verticalSpaces {
          if loc == locations[locIdx][1] {
            verticalSpaces = slices.Delete(verticalSpaces, idx, idx + 1)
          }
        }
      }
    } else {
      for addLines := 0; addLines < expansionRate; addLines++ {
        horizontalStretchMap = append(horizontalStretchMap, strings.Repeat(".", lineLen))
      }
    }
    horizontalStretchMap = append(horizontalStretchMap, fileContents[lineIdx])
  }

  for i := 0; i < len(horizontalStretchMap); i++ {
    var tmpSlice []rune
    for j := 0; j < lineLen; j++ {
      for _, loc := range verticalSpaces {
        if loc == j {
          for addLines := 0; addLines < expansionRate; addLines++ {
            tmpSlice = append(tmpSlice, '.')
          }
        }
      }
      tmpSlice = append(tmpSlice, []rune(horizontalStretchMap[i])[j])
    }
    galaxyMap = append(galaxyMap, tmpSlice)
  }

  return galaxyMap
}

func extractGalaxies(galaxyMap [][]rune) []map[string]int {
  var galaxyList []map[string]int
  for i := 0; i < len(galaxyMap); i++ {
    for j := 0; j < len(galaxyMap[i]); j++ {
      if galaxyMap[i][j] == '#' {
        tmpMap := make(map[string]int)
        tmpMap["x"] = j
        tmpMap["y"] = i
        galaxyList = append(galaxyList, tmpMap)
      }
    }
  }
  return galaxyList
}

func calculateStepsBetweenGalaxies(galaxyList []map[string]int) int {
  stepCount := 0
  for i := 0; i < len(galaxyList); i++ {
    for j := i + 1; j < len(galaxyList); j++ {
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

func main() {
  var filename string
  var expansionRate int
  flag.StringVar(&filename, "i", "input.txt", "Specify input file for the program")
  flag.IntVar(&expansionRate, "e", 1, "Specify the rate of expansion")
  flag.Parse()

  fileContents, err := readFile(filename)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  spaceMap := explodeSpace(fileContents, expansionRate)
  galaxies := extractGalaxies(spaceMap)

  steps := calculateStepsBetweenGalaxies(galaxies)

  fmt.Println(steps)
}
