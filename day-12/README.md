# Day 12 - Hot Springs

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Each line is in the format `[.?#]+ (,?[0-9]+)+`

## Part one

A file contains a series of storage notices, broken down into working (`.`), damaged (`#`),
and unknown (`?`) items. The numbers of the end of these symbols describes the number of
damaged items in each series.

For example, the following would be a completely known record:

```text
#.#.### 1,1,3
.#...#....###. 1,1,3
.#.###.#.###### 1,3,1,6
####.#...#... 4,1,1
#....######..#####. 1,6,5
.###.##....# 3,2,1
```

In reality, the number of unknown items means this isn't as clear as before, and we would.
like to know the number of valid arrangements in each item series that can fit the damaged
item log.

For example:

```text
???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
```

In the first line (`???.### 1,1,3`), there is exactly one way separate groups of one, one,
and three broken springs (in that order) can appear in that row: the first three unknown
springs must be broken, then operational, then broken (`#.#`), making the whole row 
`#.#.###`.

The second line is more interesting: `.??..??...?##. 1,1,3` could be a total of four
different arrangements. The last `?` must always be broken (to satisfy the final contiguous
group of three broken springs), and each `??` must hide exactly one of the two broken
springs. (Neither `??` could be both broken springs or they would form a single contiguous
group of two; if that were true, the numbers afterward would have been `2,3` instead.)
Since each `??` can either be `#.` or `.#`, there are four possible arrangements of
springs.

The last line is actually consistent with ten different arrangements! Because the first
number is `3`, the first and second `?` must both be `.` (if either were `#`, the first
number would have to be `4` or higher). However, the remaining run of unknown spring
conditions have many different ways they could hold groups of two and one broken springs:

```text
?###???????? 3,2,1
.###.##.#...
.###.##..#..
.###.##...#.
.###.##....#
.###..##.#..
.###..##..#.
.###..##...#
.###...##.#.
.###...##..#
.###....##.#
```

In this example, the number of possible arrangements for each row is:

* `???.### 1,1,3` - 1 arrangement
* `.??..??...?##. 1,1,3` - 4 arrangements
* `?#?#?#?#?#?#?#? 1,3,1,6` - 1 arrangement
* `????.#...#... 4,1,1` - 1 arrangement
* `????.######..#####. 1,6,5` - 4 arrangements
* `?###???????? 3,2,1` - 10 arrangements

Adding all of the possible arrangement counts together produces a total of `21`
arrangements.

## Part two

## This script

This script uses Go and can be run with the following command:

```bash
go run . -i input.txt
```

This will answer part one as described above.
