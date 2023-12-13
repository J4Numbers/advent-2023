# Day 13 - Point of Incidence

This is a simple script which operates on an input file with the following requirements:

* It is made up of several small maps
* Each map is made up of between 1 and many fixed-width lines
* Each line can only contain the following symbols `[.#]`

## Part one

A file contains a map. Somewhere in this map on either the horizontal or diagonal axis,
there are reflective lines where the map is mirrored on the axis. For each of those
reflective lines, run the following algorithm:

* If the reflection axis is vertical, add the number of columns to the left of the axis
  to your answer
* If the reflection axis is horizontal, multiply the number of rows above the axis by 100
  and add them to the answer

i.e.

```text
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
```

Contains two maps with two different lines of reflection. The first one contains a
reflection across a vertical line between two columns; arrows on each of the two
columns point at the line between the columns:

```text
123456789
    ><   
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.
    ><   
123456789
```

In this pattern, the line of reflection is the vertical line between columns 5 and 6.
Because the vertical line is not perfectly in the middle of the pattern, part of the
pattern (column 1) has nowhere to reflect onto and can be ignored; every other column
has a reflected column within the pattern and must match exactly: column 2 matches
column 9, column 3 matches 8, 4 matches 7, and 5 matches 6.

The second pattern reflects across a horizontal line instead:

```text
1 #...##..# 1
2 #....#..# 2
3 ..##..### 3
4v#####.##.v4
5^#####.##.^5
6 ..##..### 6
7 #....#..# 7
```

This pattern reflects across the horizontal line between rows 4 and 5. Row 1 would
reflect with a hypothetical row 8, but since that's not in the pattern, row 1 doesn't
need to match anything. The remaining rows match: row 2 matches row 7, row 3 matches
row 6, and row 4 matches row 5.

To get the result, we look at the first map, which had a reflective line with five columns
to its left, while the second pattern had 4 rows above it, meaning `(4 * 100) + 5 = 405`

## Part two

Alternatively, each map contains 1 difference that changes the reflection axis. This means
that either one `.` turns into a `#`, or one `#` turns into a `.`.

This changes the example above as follows:

The first pattern's smudge is in the top-left corner. If the top-left `#` were instead `.`,
it would have a different, horizontal line of reflection:

```text
1 ..##..##. 1
2 ..#.##.#. 2
3v##......#v3
4^##......#^4
5 ..#.##.#. 5
6 ..##..##. 6
7 #.#.##.#. 7
```

With the smudge in the top-left corner repaired, a new horizontal line of reflection
between rows 3 and 4 now exists. Row 7 has no corresponding reflected row and can be
ignored, but every other row matches exactly: row 1 matches row 6, row 2 matches row 5,
and row 3 matches row 4.

In the second pattern, the smudge can be fixed by changing the fifth symbol on row 2
from `.` to `#`:

```text
1v#...##..#v1
2^#...##..#^2
3 ..##..### 3
4 #####.##. 4
5 #####.##. 5
6 ..##..### 6
7 #....#..# 7
```

Now, the pattern has a different horizontal line of reflection between rows 1 and 2.

With the new reflection lines, the example is now `400`, with the first having 3 rows
above a horizontal line, and the second having 1.

## This script

This script uses Go and can be run with the following command:

```bash
go run . -i input.txt
```

This will answer part one as described above. Part 2 can be answered with the `-d` flag
which sets the number of differences.

```bash
go run . -i input.txt -d 1
```
