# Day 1 - Trebuchet!?

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Each line contains a string of unbroken characters in the set of `[a-z][0-9]`

Each line should be converted to a single two-digit number made up of the first digit and the
last digit of each line (or the same digit if there is only one) to form a single two-digit
number.

All of these numbers should be added together to create a final summed result.

For example:

```txt
1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
```

In this example, the values of these four lines are `12`, `38`, `15`, and `77`. Adding these
together produces 142.

## This script

This script uses JavaScript and can be run with the following command:

```bash
npm i
node src/app.js -i input.txt
```

Which will run the above scenario on the given input file after installing any required packages.

## Part two

Part two says that each line may also include literal numbers alongside the digits.

For example:

```txt
two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen
```

In this example, the calibration values are `29`, `83`, `13`, `24`, `42`, `14`, and `76`. Adding
these together produces `281`.

### This script

To run the script in the literals mode of operation, the following command can be used after
installation:

```bash
node src/app.js -i input.txt --mode literals
```
