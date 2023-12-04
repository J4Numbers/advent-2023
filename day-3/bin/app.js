#!/usr/bin/env node

import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

import * as fs from 'fs';

const argsEngine = yargs(hideBin(process.argv));
const args = argsEngine.wrap(argsEngine.terminalWidth())
  .env('J4_ADVENT_3')
  .options({
    input: {
      alias: 'i',
      type: 'string',
      description: 'The input file to run this algorithm on',
      demandOption: true,
    },
    debug: {
      type: 'boolean',
      default: false,
      description: 'Enable debug logging',
    },
  })
  .help()
  .parse();

const debug = args.debug;
let lineSearch = new RegExp(/[^0-9.]/, 'ig');

/**
 * @param line - Print a debug line if the program is set to debug
 */
const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

/**
 * Find all symbols on a given line and generate a symbol object reference for them
 *
 * @param line - The line to scan over
 * @param lineNumber - The line number in the total file
 * @returns A symbol object containing the line number, index of the symbol, and the
 * symbol itself
 */
const checkSymbolCount = (line, lineNumber) => {
  const symbolPositions = [];
  let foundSymbol;
  while ((foundSymbol = lineSearch.exec(line)) !== null) {
    symbolPositions.push({
      lineNumber,
      index: foundSymbol.index,
      symbol: foundSymbol[0],
    });
  }
  return symbolPositions;
};

/**
 * Given a line and the index where a numeral is in that line, find the first
 * numeral in that number and return the whole number from that line.
 *
 * @param line - The full line that is being looked up
 * @param index - The starting position where a number has been confirmed to exist
 * @returns A found number with the starting index of the first numeral and the full
 * value of that number when worked through to the end
 */
const lookBack = (line, index) => {
  let workingIndex = index;
  let pointerChar = line.charAt(workingIndex);
  while (/[0-9]/.test(pointerChar)) {
    workingIndex -= 1;
    pointerChar = line.charAt(workingIndex);
  }
  const fullNumber = /^[0-9]+/.exec(line.substring(workingIndex + 1));
  return {
    startIndex: workingIndex + 1,
    value: Number(fullNumber[0]),
  };
};

/**
 * Given a line and a symbol which has been found (can be above, below, or on this
 * line), find any numbers which are around that symbol on this line. This is done
 * by generating three separate search strings - one after the symbol, one on the
 * symbol's location, and one before the symbol - and performing a lookback check
 * if any of those points are numeric.
 *
 * These three may be resolved down into a single number later in this function,
 * but for the time being, this exhaustive method allows us to find all available
 * possibilities.
 *
 * @param line - The line to check
 * @param lineNumber - The current line number to record against any findings
 * @param symbolObj - The symbol that we're checking adjacent to for any numbers
 * @returns A list of found numbers around the given symbol - between 0 and 2 in
 * length - containing the value that was found, its start index and line number,
 * and the symbol that it maps back against.
 */
const retrievePartNumbers = (line, lineNumber, symbolObj) => {
  const foundNumbers = [];
  if (/[0-9]/.test(line.charAt(symbolObj.index + 1))) {
    foundNumbers.push(lookBack(line, symbolObj.index + 1));
  }
  if (/[0-9]/.test(line.charAt(symbolObj.index))) {
    foundNumbers.push(lookBack(line, symbolObj.index));
  }
  if (/[0-9]/.test(line.charAt(symbolObj.index - 1))) {
    foundNumbers.push(lookBack(line, symbolObj.index - 1));
  }
  return foundNumbers.reduce((ongoing, current) => {
    if (ongoing.filter((val => val.startIndex === current.startIndex)).length === 0) {
      ongoing.push({
        lineNumber,
        symbol: symbolObj,
        ...current,
      });
    }
    return ongoing;
  }, []);
};

/**
 * Flip a provided map into a map of symbols to the value numbers that surround them
 *
 * @param symbolMap - The ongoing map of symbols that will map to the values surrounding
 *                    them. Keys are `${symbol}.${line}.${index}` - i.e. '@.10.29'
 * @param currentPartNumber - The part value that is currently being iterated over. Must
 *                            contain a symbol object to be reduced in this function
 * @returns An updated symbol map with the currentPartNumber included inside
 */
const inverseOnSymbol = (symbolMap, currentPartNumber) => {
  const mapKey = `${currentPartNumber.symbol.symbol}.${currentPartNumber.symbol.lineNumber}.${currentPartNumber.symbol.index}`;
  if (!Object.keys(symbolMap).includes(mapKey)) {
    symbolMap[mapKey] = [];
  }
  symbolMap[mapKey].push(currentPartNumber);
  return symbolMap;
};

// MAIN CODE STARTS HERE

// A variable containing the previous line that has been read
let lookBehind;
// An array containing all symbols that were in the previous line
let prevLineLog = [];
// An array of all unique numbers that were found adjacent to symbols
let foundNumberLog = [];
// Currently read line (used for tracking of unique symbols and numbers)
let lineCount = 1;
// Ongoing total count of part numbers
let runningTotal = 0;

// Read the file in and split on newline
const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    // Get all symbols on the current line, and for each symbol...
    const currentLineSymbols = checkSymbolCount(line, lineCount);
    currentLineSymbols.forEach((symbolObj) => {
      let testString;
      if (lookBehind) {
        // If there is a previous line available, see whether there is a number in the three
        // characters above this found symbol
        testString = lookBehind.substring(symbolObj.index - 1, symbolObj.index + 2);
        if (/[0-9]/.test(testString)) {
          foundNumberLog.push(...retrievePartNumbers(lookBehind, lineCount - 1, symbolObj));
        }
      }
      // See if there is a number in the two characters around this symbol
      testString = line.substring(symbolObj.index - 1, symbolObj.index + 2);
      if (/[0-9]/.test(testString)) {
        foundNumberLog.push(...retrievePartNumbers(line, lineCount, symbolObj));
      }
    });
    // For the previous set of symbols on the last line, go over them and ask if there is
    // a number in the three characters below that found symbol
    prevLineLog.forEach((symbolObj) => {
      let testString = line.substring(symbolObj.index - 1, symbolObj.index + 2);
      if (/[0-9]/.test(testString)) {
        foundNumberLog.push(...retrievePartNumbers(line, lineCount, symbolObj));
      }
    });
    // Reduce down the found numbers on line and start index to ensure we have no duplicates
    foundNumberLog = foundNumberLog.reduce((ongoing, current) => {
      if (ongoing.filter((val) => val.lineNumber === current.lineNumber && val.startIndex === current.startIndex).length === 0) {
        ongoing.push(current);
      }
      return ongoing;
    }, []);

    // Update the running total and debug!
    runningTotal = foundNumberLog.reduce((ongoing, current) => ongoing + current.value, 0);
    debugLine(`Line ${lineCount} - ${currentLineSymbols.length} symbols found - ${foundNumberLog.length} distinct parts - Running total ${runningTotal}`);

    // Update our working variables and set the last lines as required
    lineCount += 1;
    lookBehind = line;
    prevLineLog = currentLineSymbols;
  });

// GEAR WORK STARTS HERE

// Only consider numbers surrounding a '*' which have exactly two associated numbers
const invertedSymbolTracking = foundNumberLog.reduce(inverseOnSymbol, {});
const validGears = Object.keys(invertedSymbolTracking)
  .filter((symbolKey) => symbolKey.startsWith('*') && invertedSymbolTracking[symbolKey].length === 2);

// Calculate those gear ratios!
const gearTotal = validGears
  .reduce((total, gearPos) => total + invertedSymbolTracking[gearPos]
    .reduce((total, current) => total * current.value, 1), 0);

console.log(`Part total - ${runningTotal}`);
console.log(`Gear total - ${gearTotal}`);
