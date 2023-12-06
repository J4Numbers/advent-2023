# Day 6 - Wait for it

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 2 and many lines
* At least one line contains a time series in the format of `Time:( +[0-9]+)+`
* At least one line contains a distance series in the format of `Distance:( +[0-9]+)+`

## Part one

The file contains two tuple sets - one for the time provided in a given round, and one for
the distance that acts as a score to beat for a given round. In the time available, time can
be spent charging or running - one point of charging increases our speed by 1 point (which
is then constant for the rest of the round).

For example:

```txt
Time:      7  15   30
Distance:  9  40  200
```

This document describes three races:

* The first race lasts 7 milliseconds. The record distance in this race is 9 millimeters.
* The second race lasts 15 milliseconds. The record distance in this race is 40 millimeters.
* The third race lasts 30 milliseconds. The record distance in this race is 200 millimeters.

In the first race, this leaves us with 8 possibilities, where we charge for `0-7` seconds, and
run for `7-0` seconds in an inverse relationship. They are as follows:

* `0` charge, `7` travel. Because no speed has been built up, the player will never move, resulting in
  a final score of `0`, not surpassing the high score of `9`.
* `1` charge, `6` travel. This will result in a constant speed of `1`, resulting in a final score of
  `6`, not surpassing the high score of `9`.
* `2` charge, `5` travel. This will result in a constant speed of `2`, resulting in a final score of
  `10`, surpassing the high score.
* `3` charge, `4` travel. This will result in a constant speed of `3`, resulting in a final score of
  `12`, surpassing the high score.
* `4` charge, `3` travel. This will result in a constant speed of `4`, resulting in a final score of
  `12`, surpassing the high score.
* `5` charge, `2` travel. This will result in a constant speed of `5`, resulting in a final score of
  `10`, surpassing the high score.
* `6` charge, `1` travel. This will result in a constant speed of `6`, resulting in a final score of
  `6`, not surpassing the high score of `9`.
* `7` charge, `0` travel. While the end speed will be `7`, there will be no movement time, resulting
  in a final score of `0`.

This leaves us `4` different options for how we could win in this race.

In the second race, you could charge for at least 4 milliseconds and at most 11 milliseconds and beat
the high score, a total of `8` different ways to win.

In the third race, you could charge for at least 11 milliseconds and no more than 19 milliseconds and
still beat the high score, a total of `9` ways you could win.

To see how much margin of error you have, determine the number of ways you can beat the record in each
race; in this example, if you multiply these values together, you get `288` (`4 * 8 * 9`).

## Part two

Alternatively, the race sheet above only describes one race (just with all the numbers spread out).

So, the example from before:

```txt
Time:      7  15   30
Distance:  9  40  200
```

...now instead means this:

```txt
Time:      71530
Distance:  940200
```

Now, you have to figure out how many ways there are to win this single race. In this example, the race
lasts for `71530` milliseconds and the record distance you need to beat is `940200` millimeters. You
could hold the button anywhere from `14` to `71516` milliseconds and beat the record, a total of
`71503` ways!

## This script

This script uses Python and can be run with the following command:

```bash
python speedster.py -i input.txt
```

This will answer part one of the question as described above. To perform part two, the `--mode` flag
must be set to `join`. i.e.

```bash
python speedster.py -i input.txt --mode join
```
