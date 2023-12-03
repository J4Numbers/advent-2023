# Day 2 - Cube Conundrum

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Each line contains a string representation of a series of observations in the form of
  `Game [0-0]+:(<round>;? (([0-9]+ red|green|blue),? ?))+`

## Part one

Each game should be marked as valid or invalid depending on whether the observed number of
cubes fits within an established maximum. The ID of these valid games should be summed up
to create a final result.

For example:

```txt
Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
```

If the allowed maximum cubes are `12` red cubes, `13` green cubes, and `14` blue cubes,
then only games `1`, `2`, and `5` would have been possible. The total value of these IDs
being `8`.

## Part two

In addition, we should find out the minimum number of cubes required to play each game, and
create a power set of those minimum values. With this power set (the minimum number of red,
green, and blue cubes in each game multiplied together), it should be summed with all other
games to give a final result.

For example:

```txt
Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
```

The various minimum power sets from lines 1-5 are `48`, `12`, `1560`, `630`, and `36`
respectively. The sum of all of these is `2286`.

## This script

This script uses JavaScript and can be run with the following command:

```bash
npm i
node src/app.js -i input.txt
```

Which will run the above scenario on the given input file after installing any required packages
and will return the answer to both of the above questions.
