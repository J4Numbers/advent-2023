# Day 3 - Gear ratios

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Each line contains a fixed length string consisting of dots (`.`), arbitrary ascii symbols,
  and numeric values

## Part one

Arbitrary symbols (i.e. `#*/$` etc) act as a marker. Any number which is adjacent (horizontal,
vertical, or diagonal) to one of those symbols is counted as a 'valid' number. The solution to
the puzzle is the sum of all those 'valid' numbers

> A number which is adjacent to multiple symbols is only counted once.

For example:

```txt
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
```

In this example, two numbers are not part numbers because they are not adjacent to a symbol:
`114` (top right) and `58` (middle right). Every other number is adjacent to a symbol and so is
a part number; their sum is `4361`.

## Part two

In addition, we would like to find the sum of all gear ratios within the puzzle block. Any `*`
symbol should be treated as a gear, only if it is adjacent to _exactly_ two numbers. The ratio
of a gear is calculated by multiplying its adjacent numbers together.

For example:

```txt
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
```

In this example, there are only two valid gears. The first is in the top left; it has part
numbers `467` and `35`, so its gear ratio is `16345`. The second gear is in the lower right;
its gear ratio is `451490`. (The `*` adjacent to `617` is not a gear because it is only adjacent
to one part number.) Adding up all of the gear ratios produces `467835`.

## This script

This script uses JavaScript and can be run with the following command:

```bash
npm i
node src/app.js -i input.txt
```

Which will run the above scenario on the given input file after installing any required packages
and will return the answer to both of the above questions.
