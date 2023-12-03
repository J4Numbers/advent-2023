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

const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

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
}

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
}

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

const inverseOnSymbol = (symbolMap, currentPartNumber) => {
  const mapKey = `${currentPartNumber.symbol.symbol}.${currentPartNumber.symbol.lineNumber}.${currentPartNumber.symbol.index}`;
  if (!Object.keys(symbolMap).includes(mapKey)) {
    symbolMap[mapKey] = [];
  }
  symbolMap[mapKey].push(currentPartNumber);
  return symbolMap;
};

let lookBehind;
let prevLineLog = [];
let foundNumberLog = [];
let lineCount = 1;
let runningTotal = 0;

const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    const currentLineSymbols = checkSymbolCount(line, lineCount);
    currentLineSymbols.forEach((symbolObj) => {
      let testString;
      if (lookBehind) {
        testString = lookBehind.substring(symbolObj.index - 1, symbolObj.index + 2);
        if (/[0-9]/.test(testString)) {
          foundNumberLog.push(...retrievePartNumbers(lookBehind, lineCount - 1, symbolObj));
        }
      }
      testString = line.substring(symbolObj.index - 1, symbolObj.index + 2);
      if (/[0-9]/.test(testString)) {
        foundNumberLog.push(...retrievePartNumbers(line, lineCount, symbolObj));
      }
    });
    prevLineLog.forEach((symbolObj) => {
      let testString = line.substring(symbolObj.index - 1, symbolObj.index + 2);
      if (/[0-9]/.test(testString)) {
        foundNumberLog.push(...retrievePartNumbers(line, lineCount, symbolObj));
      }
    });
    foundNumberLog = foundNumberLog.reduce((ongoing, current) => {
      if (ongoing.filter((val) => val.lineNumber === current.lineNumber && val.startIndex === current.startIndex).length === 0) {
        ongoing.push(current);
      }
      return ongoing;
    }, []);
    runningTotal = foundNumberLog.reduce((ongoing, current) => ongoing + current.value, 0);
    debugLine(`Line ${lineCount} - ${currentLineSymbols.length} symbols found - ${foundNumberLog.length} distinct parts - Running total ${runningTotal}`);
    lineCount += 1;
    lookBehind = line;
    prevLineLog = currentLineSymbols;
  });

const invertedSymbolTracking = foundNumberLog.reduce(inverseOnSymbol, {});
const validGears = Object.keys(invertedSymbolTracking)
  .filter((symbolKey) => symbolKey.startsWith('*') && invertedSymbolTracking[symbolKey].length === 2);

const gearTotal = validGears
  .reduce((total, gearPos) => total + invertedSymbolTracking[gearPos]
    .reduce((total, current) => total * current.value, 1), 0);

console.log(`Part total - ${runningTotal}`);
console.log(`Gear total - ${gearTotal}`);
