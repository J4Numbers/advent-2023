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

let validCards = [];
let copyCount = {};
let maxLines = 0;
let count = 0;
const cardSplitRegex = new RegExp('^card *([0-9]+): +([0-9 ]+) +[|] +([0-9 ]+)$', 'i');

const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

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

const file = fs.readFileSync(args.input).toString('utf-8');
const lines = file.split('\n');
maxLines = lines.length;
lines.forEach((line) => {
    const foundMatches = calculateMatches(line);
    if (foundMatches.cardId > 0) {
      validCards.push(foundMatches.cardId);
      const matchCount = foundMatches.matchedNumbers.length;
      markCopies(foundMatches.cardId, matchCount);
      let cardVal = 0;
      if (matchCount > 0) {
        cardVal = Math.pow(2, matchCount - 1);
      }
      count += cardVal;
      const copyStatement = Object.keys(copyCount).includes(`${foundMatches.cardId}`)
        ? `${copyCount[foundMatches.cardId]} copies`
        : `0 copies`;
      debugLine(`Card ${foundMatches.cardId} - ${matchCount} matches - ${copyStatement} - Value ${cardVal} - Total ${count}`);
    }
  });

const totalCards = Object.keys(copyCount)
  .filter((cardId) => validCards.includes(Number(cardId)))
  .reduce((totalCards, currentId) => totalCards + copyCount[currentId], validCards.length);

console.log(`Valid count - ${count}`);
console.log(`Total scratchcards - ${totalCards}`);
