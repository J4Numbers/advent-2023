# Day 10 - Pipe Maze

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many fixed-width lines
* Each line can only contain the following symbols `[-|JFL7.]`
* The letter `S` can only appear at most once in the puzzle input

## Part one

A file contains a large maze of pipes. The pipes are either straight sections (`-|`) or bends in
direction (`JLF7`). Within this file, there is also a starting position (`S`) that is within a
looping segment of the maze. Any position in the maze that is not important is covered with a dot
(`.`).

The task is to find the number of anti-point of that starting position within the loop and the
number of steps that would be required to reach that anti-point from that starting position.

For example:

```text
-L|F7
7S-7|
L|7||
-L-J|
L|-JF
```

The loop is only the segment in the middle, but there are many extraneous parts that might try
to connect to that loop unsuccessfully. To discover the anti-point of this loop, we can count
the number of steps from the starting position as follows:

```text
.....
.012.
.1.3.
.234.
.....
```

Leaving us with an answer of `4` as the maximum number of steps one can take within the loop.

## Part two

In addition to the above, we also want to know the number of squares that have been completely
enclosed by the loop. This is not quite as easy as it sounds, however, as the flood algorithm
also needs to consider that `||` is a valid escape.

For example:

```text
..........
.S------7.
.|F----7|.
.||OOOO||.
.||OOOO||.
.|L-7F-J|.
.|II||II|.
.L--JL--J.
..........
```

While the bottom 4 squares marked as `I` are inside the loop, the eight squares marked as `O`
are able to escape through the double-pipe sets at the bottom as none of them directly cross
over the gap.

## This script

This script uses Python and can be run with the following command:

```bash
python maze_runner.py -i input.txt
```

This will answer part one and part 2 of the question as described above.
