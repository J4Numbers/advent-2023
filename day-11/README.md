# Day 11 - Cosmic Expansion

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many fixed-width lines
* Each line can only contain the following symbols `[.#]`

## Part one

A file contains the map of a galaxy, with empty space (`.`) and galaxies (`#`). This map,
however, is somewhat smaller than the true representation. In any line which _there are
no galaxies_, the line should be doubled up. This counts for both horizontal and vertical
empty lines.

i.e.

```text
.#
..
```

becomes

```text
..#
...
...
```

The task, once the map has been expanded correctly, is to find the shortest path between
every galaxy and sum them all up.

For example, given the original map:

```text
...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
```

First we identify the empty rows as follows:

```text
   v  v  v
 ...#......
 .......#..
 #.........
>..........<
 ......#...
 .#........
 .........#
>..........<
 .......#..
 #...#.....
   ^  ^  ^
```

These rows and columns need to be doubled up to produce the following expanded map:

```text
....#........
.........#...
#............
.............
.............
........#....
.#...........
............#
.............
.............
.........#...
#....#.......
```

Equipped with this expanded map, the shortest path between every pair of galaxies can be
found. To simplify things, we give each galaxy a number:

```text
....1........
.........2...
3............
.............
.............
........4....
.5...........
............6
.............
.............
.........7...
8....9.......
```

In these 9 galaxies, there are 36 pairs. Only count each pair once; order within the
pair doesn't matter. For each pair, find any shortest path between the two galaxies
using only steps that move up, down, left, or right exactly one `.` or `#` at a time

> The shortest path between two galaxies is allowed to pass through another galaxy.

For example, here is one of the shortest paths between galaxies 5 and 9:

```text
....1........
.........2...
3............
.............
.............
........4....
.5...........
.##.........6
..##.........
...##........
....##...7...
8....9.......
```

This path has length 9 because it takes a minimum of nine steps to get from galaxy 5 to
galaxy 9 (the eight locations marked `#` plus the step onto galaxy 9 itself). Here are
some other example shortest path lengths:

* Between galaxy 1 and galaxy 7: `15`
* Between galaxy 3 and galaxy 6: `17`
* Between galaxy 8 and galaxy 9: `5`

In this example, after expanding the universe, the sum of the shortest path between all
36 pairs of galaxies is `374`.

## Part two

The expansion rate is now different. Instead of the expansion of 1 per empty line, each
empty row or column should be one million times larger instead. That is, each empty row
should be replaced with `1000000` empty rows, and each empty column should be replaced
with `1000000` empty columns.

(In the example above, if each empty row or column were merely `10` times larger, the sum
of the shortest paths between every pair of galaxies would be `1030`. If each empty row or
column were merely `100` times larger, the sum of the shortest paths between every pair of
galaxies would be `8410`)

## This script

This script uses Go and can be run with the following command:

```bash
go run . -i input.txt
```

This will answer part one by default, but can be overridden to answer part two with the
`-e` flag as follows:

```bash
go run . -i input.txt -e 1000000
```
