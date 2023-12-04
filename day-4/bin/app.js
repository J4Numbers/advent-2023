#!/usr/bin/env node

import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

import * as fs from 'fs';

const argsEngine = yargs(hideBin(process.argv));
const args = argsEngine.wrap(argsEngine.terminalWidth())
  .env('J4_ADVENT_4')
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

// Set up all working variables
let validCards = [];
let copyCount = {};
let maxLines = 0;
let count = 0;

// Set a split regex to partition each line into the card ID, the winning numbers,
// and the available numbers
const cardSplitRegex = new RegExp('^card *([0-9]+): +([0-9 ]+) +[|] +([0-9 ]+)$', 'i');

/**
 * @param line - Print this line if debug is enabled
 */
const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

/**
 * For a given line, calculate how many matches there are on this card between the
 * set of winning numbers and the set of available numbers. This function returns
 * that data as an object with multiple distinct items of information.
 *
 * @param line - The line that we're finding matches on
 * @returns An object detailing the winning, discovered, and matched numbers within
 * a single card
 */
const calculateMatches = (line) => {
  let cardId = 0;
  let winningNumbers = [];
  let discoveredNumbers = [];
  let matchedNumbers = [];
  const cardDetails = cardSplitRegex.exec(line);
  if (cardDetails) {
    cardId = Number(cardDetails[1]);
    winningNumbers = cardDetails[2].split(/ +/).map((val) => Number(val)).sort();
    discoveredNumbers = cardDetails[3].split(/ +/).map((val) => Number(val)).sort();
    matchedNumbers = winningNumbers.filter((winner) => discoveredNumbers.includes(winner));
  }
  return {
    cardId,
    winningNumbers,
    discoveredNumbers,
    matchedNumbers,
  };
};

/**
 * Given an original card ID from which we are triggering copies of subsequent
 * cards and the number of matching numbers on that card, update our copy map with
 * the new numbers of required copies of that card.
 *
 * Note: If the original card ID has a number of copies already linked against it,
 * then that number of copies is added to every copy that this card requests
 *
 * @param originalCardId - The ID of the card from which we are springing copies
 *                         from
 * @param matchingCount - The number of subsequent cards that need an additional
 *                        copy, or set of copies
 */
const markCopies = (originalCardId, matchingCount) => {
  let additionValue = 1;
  if (Object.keys(copyCount).includes(`${originalCardId}`)) {
    additionValue += copyCount[originalCardId];
  }
  for (let i=1; i<=matchingCount; ++i) {
    if (maxLines >= originalCardId + i) {
      if (!Object.keys(copyCount).includes(`${originalCardId + i}`)) {
        copyCount[originalCardId + i] = 0;
      }
      copyCount[originalCardId + i] += additionValue;
    }
  }
};

// MAIN CODE STARTS HERE

// Split a file on new lines and iterate through all of them - noting down the maximum
// number of available lines that there are
const file = fs.readFileSync(args.input).toString('utf-8');
const lines = file.split('\n');
maxLines = lines.length;
lines.forEach((line) => {
  // Calculate the matches found on each line
  const foundMatches = calculateMatches(line);
  if (foundMatches.cardId > 0) {
    // If the match was valid, count this as a valid card, and mark the subsequent cards
    // for copies if there were any matches
    validCards.push(foundMatches.cardId);
    const matchCount = foundMatches.matchedNumbers.length;
    markCopies(foundMatches.cardId, matchCount);

    // Calculate the value of this card
    let cardVal = 0;
    if (matchCount > 0) {
      cardVal = Math.pow(2, matchCount - 1);
    }
    count += cardVal;

    // Debug the value of this card and the number of copies
    const copyStatement = Object.keys(copyCount).includes(`${foundMatches.cardId}`)
      ? `${copyCount[foundMatches.cardId]} copies`
      : `0 copies`;
    debugLine(`Card ${foundMatches.cardId} - ${matchCount} matches - ${copyStatement} - Value ${cardVal} - Total ${count}`);
  }
});

// Calculate the total number of scratchcards, including originals and copies
const totalCards = Object.keys(copyCount)
  .filter((cardId) => validCards.includes(Number(cardId)))
  .reduce((totalCards, currentId) => totalCards + copyCount[currentId], validCards.length);

// Log outputs
console.log(`Valid count - ${count}`);
console.log(`Total scratchcards - ${totalCards}`);
