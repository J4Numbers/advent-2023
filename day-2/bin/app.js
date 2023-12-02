#!/usr/bin/env node

import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

import * as fs from 'fs';

const argsEngine = yargs(hideBin(process.argv));
const args = argsEngine.wrap(argsEngine.terminalWidth())
  .env('J4_ADVENT_2')
  .options({
    input: {
      alias: 'i',
      type: 'string',
      description: 'The input file to run this algorithm on',
      demandOption: true,
    },
    maxRed: {
      alias: ['red', 'r'],
      type: 'number',
      description: 'The maximum number of red cubes that can exist in a game',
      default: 12,
    },
    maxBlue: {
      alias: ['blue', 'b'],
      type: 'number',
      description: 'The maximum number of blue cubes that can exist in a game',
      default: 14,
    },
    maxGreen: {
      alias: ['green', 'g'],
      type: 'number',
      description: 'The maximum number of green cubes that can exist in a game',
      default: 13,
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
const allowedCubeMax = {
  red: args.maxRed,
  blue: args.maxBlue,
  green: args.maxGreen,
}

let count = 0;
const gameRegex = new RegExp('^game ([0-9]+): ((([^;]+)(;|\n)?)+)$', 'i');
const gameSplitterRegex = new RegExp('(([0-9]+ green|[0-9]+ red|[0-9]+ blue)+,? ?)+(;|$)', 'ig');
const pullSplitterRegex = new RegExp('([0-9]+) (green|red|blue)', 'ig');

const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

const calculateGameOutcome = (line) => {
  let gameId = 0;
  let maxRed = 0;
  let maxBlue = 0;
  let maxGreen = 0;
  const gameDetails = gameRegex.exec(line);
  if (gameDetails) {
    let onePull;
    gameId = Number(gameDetails[1]);
    while ((onePull = gameSplitterRegex.exec(gameDetails[2])) !== null) {
      let oneColour;
      while ((oneColour = pullSplitterRegex.exec(onePull[0])) !== null) {
        const cubeCount = Number(oneColour[1]);
        switch (oneColour[2]) {
          case 'red':
            if (cubeCount > maxRed) {
              maxRed = cubeCount;
              break;
            }
          case 'blue':
            if (cubeCount > maxBlue) {
              maxBlue = cubeCount;
              break;
            }
          case 'green':
            if (cubeCount > maxGreen) {
              maxGreen = cubeCount;
            }
        }
      }
      debugLine(`Game ${gameId} - ${onePull[0]} - Max red ${maxRed} - Max blue ${maxBlue} - Max green ${maxGreen}`);
    }
  }
  return {
    gameId,
    red: maxRed,
    green: maxGreen,
    blue: maxBlue,
  };
};

const testGameValid = (maxDiscoveredCubes, maxAllowedCubes) => {
  return !((maxDiscoveredCubes.red > maxAllowedCubes.red)
    || (maxDiscoveredCubes.blue > maxAllowedCubes.blue)
    || (maxDiscoveredCubes.green > maxAllowedCubes.green));
}

const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    const discoveredCubeMax = calculateGameOutcome(line);
    if (discoveredCubeMax.gameId > 0) {
      const baseGameLog = `Game ${discoveredCubeMax.gameId} - Red ${discoveredCubeMax.red}/${allowedCubeMax.red} - Blue ${discoveredCubeMax.blue}/${allowedCubeMax.blue} - Green ${discoveredCubeMax.green}/${allowedCubeMax.green}`;
      if (testGameValid(discoveredCubeMax, allowedCubeMax)) {
        count += discoveredCubeMax.gameId;
        debugLine(`${baseGameLog} - VALID - Count ${count}`);
      } else {
        debugLine(`${baseGameLog} - INVALID - Count ${count}`);
      }
    }
  });

console.log(count);
