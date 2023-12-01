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
    debug: {
      type: 'boolean',
      default: false,
      description: 'Enable debug logging',
    }
  })
  .help()
  .parse();

const debug = args.debug;
let lineCount = 1;
let count = 0;
const twoDigitRegexTest = new RegExp('^[^0-9]*([0-9]).*([0-9])[^0-9]*$', 'i');
const oneDigitRegexTest = new RegExp('^[^0-9]*([0-9])[^0-9]*$', 'i');

const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
}

const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    const output = twoDigitRegexTest.exec(line);
    if (output) {
      const value = `${output[1]}${output[2]}`;
      count += Number(value);
      debugLine(`Line ${lineCount} - Found ${Number(value)} - Running ${count}`);
    } else {
      const oneDigitOutput = oneDigitRegexTest.exec(line);
      if (oneDigitOutput) {
        const value = `${oneDigitOutput[1]}${oneDigitOutput[1]}`;
        count += Number(value);
        debugLine(`Line ${lineCount} - Found ${value} - Running ${count}`);
      } else {
        debugLine(`Line ${lineCount} - Not valid`);
      }
    }
    lineCount += 1;
  });

console.log(count);
