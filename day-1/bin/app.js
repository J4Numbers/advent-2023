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
const twoDigitRegexTest = new RegExp('^[^0-9]*([0-9]).*([0-9])[^0-9]*$', 'i');
const oneDigitRegexTest = new RegExp('^[^0-9]*([0-9])[^0-9]*$', 'i');


const literalNumbers = 'one|two|three|four|five|six|seven|eight|nine';
const reversedLiteralNumbers = literalNumbers.split('').reverse().join('');
const literalOptions = `([0-9]|${literalNumbers})`;
const reversedLiteralOptions = `([0-9]|${reversedLiteralNumbers})`;

const literalRegex = new RegExp(literalOptions, 'i');
const reversedLiteralRegex = new RegExp(reversedLiteralOptions, 'i');


const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

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

const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    let outValue = 0;
    if (args.mode === 'digits') {
      outValue = testDigits(line);
    }
    if (args.mode === 'literals') {
      outValue = testLiterals(line);
    }
    count += outValue;
    debugLine(`Line ${lineCount} - ${line} - Found ${outValue} - Running ${count}`);
    lineCount += 1;
  });

console.log(count);
