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

// Set the allowed maximums as required
const allowedCubeMax = {
  red: args.maxRed,
  blue: args.maxBlue,
  green: args.maxGreen,
}

let count = 0;
let minPowersTotal = 0;

// Set up some basic regex splits to split out individual games and pulls and colours from
// each line
const gameRegex = new RegExp('^game ([0-9]+): ((([^;]+)(;|\n)?)+)$', 'i');
const gameSplitterRegex = new RegExp('(([0-9]+ green|[0-9]+ red|[0-9]+ blue)+,? ?)+(;|$)', 'ig');
const pullSplitterRegex = new RegExp('([0-9]+) (green|red|blue)', 'ig');

/**
 * @param line - Print to console if debug is enabled
 */
const debugLine = (line) => {
  if (debug) {
    console.log(line);
  }
};

/**
 * Given a line which describes a game made up of several rounds, where each round contains
 * a set of tuples, containing the number and colour of cubes that were found, we want to
 * calculate the minimum required cubes of each colour to be in this bag for this game to
 * have been possible.
 *
 * @param line - A line containing inforamtion about a single game made up of several rounds
 * @returns an object of the minimum required red, green, and blue cubes to have played this
 * game
 */
const calculateGameOutcome = (line) => {
  let gameId = 0;
  let maxRed = 0;
  let maxBlue = 0;
  let maxGreen = 0;
  const gameDetails = gameRegex.exec(line);
  // If this line contains a valid game...
  if (gameDetails) {
    let onePull;
    gameId = Number(gameDetails[1]);
    // Split the game into rounds
    while ((onePull = gameSplitterRegex.exec(gameDetails[2])) !== null) {
      let oneColour;
      // Split the round into individual colours
      while ((oneColour = pullSplitterRegex.exec(onePull[0])) !== null) {
        const cubeCount = Number(oneColour[1]);
        // If there is more of a given colour than we have already seen, then update the
        // required minimum of that colour accordingly
        switch (oneColour[2]) {
          case 'red':
            if (cubeCount > maxRed) {
              maxRed = cubeCount;
            }
            break;
          case 'blue':
            if (cubeCount > maxBlue) {
              maxBlue = cubeCount;
            }
            break;
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

/**
 * Given a reported game, return whether that game is even possible if we have a set
 * number of available cubes.
 *
 * @param maxDiscoveredCubes - The minimum cubes of a given colour to have played a
 *                             given game
 * @param maxAllowedCubes - The known maximum number of cubes that we have available
 *                          in each colour
 * @returns true if there were more (or the same amount of) allowed cubes of every
 * colour than discovered cubes
 */
const testGameValid = (maxDiscoveredCubes, maxAllowedCubes) => {
  return !((maxDiscoveredCubes.red > maxAllowedCubes.red)
    || (maxDiscoveredCubes.blue > maxAllowedCubes.blue)
    || (maxDiscoveredCubes.green > maxAllowedCubes.green));
};

// MAIN CODE STARTS HERE

// For each line in the provided file...
const file = fs.readFileSync(args.input).toString('utf-8');
file.split('\n')
  .forEach((line) => {
    // Calculate the minimum required cubes to play this game
    const discoveredCubeMax = calculateGameOutcome(line);
    if (discoveredCubeMax.gameId > 0) {
      // Generate the power value of the cubes in this game
      let minCubesPower = (discoveredCubeMax.red * discoveredCubeMax.green * discoveredCubeMax.blue);
      minPowersTotal += minCubesPower;

      // Set up a base debug log template
      const baseGameLog = `Game ${discoveredCubeMax.gameId} - Red ${discoveredCubeMax.red}/${allowedCubeMax.red} - Blue ${discoveredCubeMax.blue}/${allowedCubeMax.blue} - Green ${discoveredCubeMax.green}/${allowedCubeMax.green} - Power ${minCubesPower} - Total power ${minPowersTotal}`;

      // Debug differently depending on whether this game is valid according to our set amount
      // of allowed cubes
      if (testGameValid(discoveredCubeMax, allowedCubeMax)) {
        count += discoveredCubeMax.gameId;
        debugLine(`${baseGameLog} - VALID - Count ${count}`);
      } else {
        debugLine(`${baseGameLog} - INVALID - Count ${count}`);
      }
    }
  });

// Print end results
console.log(`Powers - ${minPowersTotal}`);
console.log(`Valid count - ${count}`);

