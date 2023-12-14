# Day 14 - Parabolic Reflector Dish

This is a simple script which operates on an input file with the following requirements:

* It is made up of several fixed-width lines
* Each line can only contain the following symbols `[.#O]`

## Part one

A file contains a map. This map can contain three types of symbol:

* `.` - Empty space
* `-` - A flat-faced rock that is fixed in place
* `O` - A rounded rock that will move when the map is tiled

We would like to know what the total weight on the north support is when tilting all
rocks towards the north.

i.e.

```text
O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
```

Start by tilting the platform so all of the rocks will slide north as far as they will
go:

```text
OOOO.#.O..
OO..#....#
OO..O##..O
O..#.OO...
........#.
..#....#.#
..O..#.O.O
..O.......
#....###..
#....#....
```

To calculate the load on the north side of the platform, we use the following algorithm.

The amount of load caused by a single rounded rock (`O`) is equal to the number of rows
from the rock to the south edge of the platform, including the row the rock is on.
(Cube-shaped rocks (`#`) don't contribute to load.) So, the amount of load caused by each
rock in each row is as follows:

```text
OOOO.#.O.. 10
OO..#....#  9
OO..O##..O  8
O..#.OO...  7
........#.  6
..#....#.#  5
..O..#.O.O  4
..O.......  3
#....###..  2
#....#....  1
```

The total load is the sum of the load caused by all of the rounded rocks. In this example,
the total load is `136`

## Part two

Now, instead of a single tilt, we want to start a tilt cycle of N, W, S, E - which counts
for one cycle. After each tilt, the rounded rocks roll as far as they can before the
platform tilts in the next direction. After one cycle, the platform will have finished
rolling the rounded rocks in those four directions in that order.

Here's what happens in the example above after each of the first few cycles:

```text
After 1 cycle:
.....#....
....#...O#
...OO##...
.OO#......
.....OOO#.
.O#...O#.#
....O#....
......OOOO
#...O###..
#..OO#....

After 2 cycles:
.....#....
....#...O#
.....##...
..O#......
.....OOO#.
.O#...O#.#
....O#...O
.......OOO
#..OO###..
#.OOO#...O

After 3 cycles:
.....#....
....#...O#
.....##...
..O#......
.....OOO#.
.O#...O#.#
....O#...O
.......OOO
#...O###.O
#.OOO#...O
```

We want to calculate the weight on the north support beam after 1000000000. In the above
example, the weight would be `64`.

## This script

This script uses Go and can be run with the following command:

```bash
go run . -i input.txt
```

This will answer part one as described above. To add in cycles, we add in the two flags
`--cycle` and `--count 1000000000` - the former of which turns on the spin cycle while
the latter sets the number of cycles.

```bash
go run . -i input.txt --cycle --count 1000000000
```
