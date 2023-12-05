#!/usr/bin/env node

import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

import * as fs from 'fs';

const argsEngine = yargs(hideBin(process.argv));
const args = argsEngine.wrap(argsEngine.terminalWidth())
  .env('J4_ADVENT_5')
  .options({
    input: {
      alias: 'i',
      type: 'string',
      description: 'The input file to run this algorithm on',
      demandOption: true,
    },
    mode: {
      type: 'string',
      options: ['seed', 'range'],
      default: 'seed',
      description: 'The mode of operation for the seed values - can be one seed per value, or a range from a start value',
    },
    split: {
      type: 'number',
      default: 1000,
      description: 'The max number to split on a range',
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

let count = 0;
let minPowersTotal = 0;

// Set up some basic regex splits to test individual lines and retrieve information
const seedLine = new RegExp('^seeds: [0-9 ]+$', 'gi');
const mapLookup = new RegExp('^([a-z]+)-to-([a-z]+) map:$', 'gi');
const rangeLookup = new RegExp('^([0-9]+) ([0-9]+) ([0-9]+)$', '');

/**
 * @param line - Print to console if debug is enabled
 */
const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

const retrieveSeeds = (line) => {
  const seedList = [];
  let seedPointer;
  if (args.mode === 'seed') {
    const seedLookup = new RegExp('([0-9]+)', 'gi');
    while ((seedPointer = seedLookup.exec(line)) !== null) {
      seedList.push({
        startVal: Number(seedPointer[1]),
      });
    }
  } else {
    const seedLookup = new RegExp('([0-9]+) ([0-9]+)', 'gi');
    while ((seedPointer = seedLookup.exec(line)) !== null) {
      let startIndex = Number(seedPointer[1]);
      let remainingRange = Number(seedPointer[2]);
      while (remainingRange > args.split) {
        seedList.push({
          startVal: startIndex,
          range: args.split,
        });
        startIndex += args.split;
        remainingRange -= args.split;
      }
      if (remainingRange > 0) {
        seedList.push({
          startVal: startIndex,
          range: remainingRange,
        });
      }
    }
  }
  return seedList;
};

const calculateDestinationValue = (sourceVal, rangeLookups) => {
  let foundVal = sourceVal;
  const validRanges = rangeLookups.filter((range) => range.sourceStartIndex <= sourceVal
    && (range.sourceStartIndex + range.length - 1) >= sourceVal);
  if (validRanges.length > 0) {
    validRanges.forEach((range) => {
      foundVal = range.destStartIndex + (sourceVal - range.sourceStartIndex);
    })
  }
  return foundVal;
};

// MAIN CODE STARTS HERE

const lookupTranslation = {};
const lookupMaps = {};
let currentMapIndex;
let seedLookupList;

let seedListFound = false;

// For each line in the provided file...
const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    if (!seedListFound) {
      if (seedLine.test(line)) {
        seedLookupList = retrieveSeeds(line);
        seedListFound = !seedListFound;
      }
    }
    const mappingTest = mapLookup.exec(line);
    if (mappingTest !== null) {
      lookupTranslation[mappingTest[1]] = mappingTest[2];
      lookupMaps[mappingTest[1]] = [];
      currentMapIndex = mappingTest[1];
    }
    const rangeTest = rangeLookup.exec(line);
    if (rangeTest !== null) {
      lookupMaps[currentMapIndex].push({
        sourceStartIndex: Number(rangeTest[2]),
        destStartIndex: Number(rangeTest[1]),
        length: Number(rangeTest[3]),
      });
    }
  });

let currentMinLocation;

seedLookupList.forEach((seed) => {
  const seedDiary = [{ seed: seed.startVal }];
  if (Object.keys(seed).includes('range')) {
    for (let i=1; i<seed.range; ++i) {
      seedDiary.push({ seed: seed.startVal + i });
    }
  };
  seedDiary.forEach((entry) => {
    let workingType = 'seed';
    while (Object.keys(lookupTranslation).includes(workingType)) {
      entry[lookupTranslation[workingType]] = calculateDestinationValue(entry[workingType], lookupMaps[workingType]);
      debugLine(`Seed ${entry['seed']} - ${workingType} ${entry[workingType]} -> ${lookupTranslation[workingType]} ${entry[lookupTranslation[workingType]]}`);
      workingType = lookupTranslation[workingType];
    }
    if (!currentMinLocation || entry.location < currentMinLocation.location) {
      currentMinLocation = entry;
    }
  });
});

// Print end results
console.log(currentMinLocation);
