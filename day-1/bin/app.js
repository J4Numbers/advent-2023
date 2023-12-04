#!/usr/bin/env node

import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

import * as fs from 'fs';

const argsEngine = yargs(hideBin(process.argv));
const args = argsEngine.wrap(argsEngine.terminalWidth())
  .env('J4_ADVENT_1')
  .options({
    input: {
      alias: 'i',
      type: 'string',
      description: 'The input file to run this algorithm on',
      demandOption: true,
    },
    mode: {
      options: ['digits', 'literals'],
      default: 'digits',
      description: 'Choose between calibrating on digits only, or to include number words too',
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
let lineCount = 1;
let count = 0;

// Used for numeric-only searching
const twoDigitRegexTest = new RegExp('^[^0-9]*([0-9]).*([0-9])[^0-9]*$', 'i');
const oneDigitRegexTest = new RegExp('^[^0-9]*([0-9])[^0-9]*$', 'i');

// Used for literal searching - generating a set of standard literal searches...
const literalNumbers = 'one|two|three|four|five|six|seven|eight|nine';
const literalOptions = `([0-9]|${literalNumbers})`;
const literalRegex = new RegExp(literalOptions, 'i');

// And reversed literal searches...
const reversedLiteralNumbers = literalNumbers.split('').reverse().join('');
const reversedLiteralOptions = `([0-9]|${reversedLiteralNumbers})`;
const reversedLiteralRegex = new RegExp(reversedLiteralOptions, 'i');

/**
 * @param line - Print this line if console debug is enabled
 */
const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

/**
 * Translate a literal text number into its numeric value. The literal text can
 * be normal or reversed.
 *
 * @param literal - the literal text of a number to be translated into a numeral
 * @returns A numeral value of the provided literal, or 0 if no translation could
 * be found
 */
const transformLiteral = (literal) => {
  let value = '0';
  switch (literal) {
    case 'one':
    case 'eno':
      value = '1';
      break;
    case 'two':
    case 'owt':
      value = '2';
      break;
    case 'three':
    case 'eerht':
      value = '3';
      break;
    case 'four':
    case 'ruof':
      value = '4';
      break;
    case 'five':
    case 'evif':
      value = '5';
      break;
    case 'six':
    case 'xis':
      value = '6';
      break;
    case 'seven':
    case 'neves':
      value = '7';
      break;
    case 'eight':
    case 'thgie':
      value = '8';
      break;
    case 'nine':
    case 'enin':
      value = '9';
  }
  return value;
};

/**
 * Run a test on a provided line to extract the first and last numerals from
 * that line and return their joined value. If there is only one numeric value
 * on a line, then return that one value twice over (i.e. 11, 22, etc.).
 *
 * @param line - the line to search for numerals in
 * @returns a two-digit number made of the first numeral that appears in the
 * string and the last numeral that appears in the string
 */
const testDigits = (line) => {
  const output = twoDigitRegexTest.exec(line);
  let value = 0;
  if (output) {
    value = Number(`${output[1]}${output[2]}`);
  } else {
    const oneDigitOutput = oneDigitRegexTest.exec(line);
    if (oneDigitOutput) {
      value = Number(`${oneDigitOutput[1]}${oneDigitOutput[1]}`);
    }
  }
  return value;
};

/**
 * With literal numbers (i.e. one, two, etc.), alongside numeric numbers (1, 2, etc),
 * the previous method of one regex string to rule them all (tm) is no-longer usable
 * (or, at least, isn't tenable). Instead, we follow the same rule as before, just more
 * explicitly...
 *
 * For the first digit, we find the first occurence in the string of 0-9 or one-nine, and
 * use that. For the second digit, we then reverse the string and find the first occurence
 * of 0-9 or eno-enin - then use that as the last occuring number in each line.
 *
 * @param line - The line that we are inspecting for all of our numbers
 * @returns A two digit number made up of the first and last numbers in the given line
 */
const testLiterals = (line) => {
  const firstDigit = literalRegex.exec(line);
  let digitOne = '0';
  let digitTwo = '0';
  let value = 0;
  if (firstDigit) {
    if (new RegExp(literalNumbers, 'i').test(firstDigit[1])) {
      digitOne = transformLiteral(firstDigit[1]);
    } else {
      digitOne = firstDigit[1];
    }
  }
  let reversedLine = line.split('').reverse().join('');
  const secondDigit = reversedLiteralRegex.exec(reversedLine);
  if (secondDigit) {
    if (new RegExp(reversedLiteralNumbers, 'i').test(secondDigit[1])) {
      digitTwo = transformLiteral(secondDigit[1]);
    } else {
      digitTwo = secondDigit[1];
    }
  }
  value = Number(`${digitOne}${digitTwo}`);
  return value;
};

// MAIN CODE STARTS HERE

// Read individual lines of the file one by one
const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    let outValue = 0;
    // If we're operating on digits-only, run that function
    if (args.mode === 'digits') {
      outValue = testDigits(line);
    }
    // And if we're on literals, choose that function instead
    if (args.mode === 'literals') {
      outValue = testLiterals(line);
    }
    count += outValue;

    // Debug line and increment working variables
    debugLine(`Line ${lineCount} - ${line} - Found ${outValue} - Running ${count}`);
    lineCount += 1;
  });

console.log(count);
