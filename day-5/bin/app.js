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

/**
 * Convert the line containing all of our seeds to the list of seeds that
 * we will iterate over later in this program. Depending on the mode of
 * operation, this will either return a simple array of seed values, or will
 * pair each seed with a range that it can contain. This range is bounded
 * by the split option on the program that allows us to not just run out
 * of memory immediately.
 *
 * @param line - The line containing a list of seeds that we will extract
 * @returns A collection of extracted seeds (and potential ranges from that
 * seed
 */
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


/**
 * Given the source value and the set of ranges to look up its corresponding
 * destination within, return the destination value accordingly. This is going
 * to be the same as the source value if no ranges match the source value, or
 * the corresponding destination value in one of our ranges.
 *
 * @param sourceVal - Our input field that will fit inside at most one of our
 *                    range lookups
 * @param rangeLookups - A list of range lookups containing source-to-destination
 *                       mapping details that we can use to translate our source
 *                       into its destination
 * @returns The corresponding destination value to our source value
 */
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

// Set up lookup maps for the word translations from seed to location, and all the
// range values from a given source
const lookupTranslation = {};
const lookupMaps = {};
let currentMapIndex;
let seedLookupList;

let seedListFound = false;

// For each line in the provided file...
// Build up the seed to location map and all its steps
const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    // Extract out the seed line if we don't already have one
    if (!seedListFound) {
      if (seedLine.test(line)) {
        seedLookupList = retrieveSeeds(line);
        seedListFound = !seedListFound;
      }
    }

    // Start a new mapping relationship when we discover the details
    // of one, including the translation, setting up a new range array,
    // and updating the map index for any future range lookup values
    const mappingTest = mapLookup.exec(line);
    if (mappingTest !== null) {
      lookupTranslation[mappingTest[1]] = mappingTest[2];
      lookupMaps[mappingTest[1]] = [];
      currentMapIndex = mappingTest[1];
    }

    // Add the details of any range addresses when we find them to the current
    // map index for the future calculation
    const rangeTest = rangeLookup.exec(line);
    if (rangeTest !== null) {
      lookupMaps[currentMapIndex].push({
        sourceStartIndex: Number(rangeTest[2]),
        destStartIndex: Number(rangeTest[1]),
        length: Number(rangeTest[3]),
      });
    }
  });

// Set up a working variable for our end location
let currentMinLocation;

// For each of the seeds that we have discovered...
seedLookupList.forEach((seed) => {
  const seedDiary = [{ seed: seed.startVal }];
  // If we are operating over a range of seeds, expand this seed into all
  // of its possible ranges...
  if (Object.keys(seed).includes('range')) {
    for (let i=1; i<seed.range; ++i) {
      seedDiary.push({ seed: seed.startVal + i });
    }
  };

  // For each seed, perform a translation lookup until we reach the final translation
  // as described in our translation dictionary, building up the translation until there
  // are no more translations to look up.
  seedDiary.forEach((entry) => {
    let workingType = 'seed';
    while (Object.keys(lookupTranslation).includes(workingType)) {
      entry[lookupTranslation[workingType]] = calculateDestinationValue(entry[workingType], lookupMaps[workingType]);
      debugLine(`Seed ${entry['seed']} - ${workingType} ${entry[workingType]} -> ${lookupTranslation[workingType]} ${entry[lookupTranslation[workingType]]}`);
      workingType = lookupTranslation[workingType];
    }

    // If our current seed has a lower location value than the one we currently have
    // stored, then update it and move on
    if (!currentMinLocation || entry.location < currentMinLocation.location) {
      currentMinLocation = entry;
    }
  });
});

// Print end results
console.log(currentMinLocation);
