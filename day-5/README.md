# Day 5 - If You Give A Seed A Fertilizer

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* There is one line that contains `seeds: ( ?[0-9]+)+` with a list of input seeds
* There is at least one set of lines that is made up of the following:
    * One mapping line in the format of `x-to-y map:`
    * One to many range lines in the format of `([0-9]+) ([0-9]+) ([0-9]+)`

## Part one

Given a set of seed values at the top, figure out which one maps to the lowest location value
after being translated through a series of maps (i.e. `seed-to-soil`, `soil-to-fertiliser`,
etc.).

The maps, along with the mapping values, contain a series of ranges in the format of:
`destStart sourceStart rangeLen`, so a range map of `5 10 3` means that a source value of `11`
would map to `6`.

The answer to the question is the lowest location value - not the original seed

For example:

```txt
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
```

Given the initial seed set of `79`, `14`, `55`, and `13`, we can follow the seeds through the
range map as described above, individually reaching a location of `82`, `43`, `86`, and `35`
respectively, leaving `35` as our lowest location value.

## Part two

Alternatively, the seed set at the top is instead a list of tuples that describe all available
seeds. Each pair of numbers makes up the starting seed, and the number of sequential seeds from
that given seed. I.e. given a pair of seed values `302 20`, that would mean that seeds `302` to
`321` all exist and need to be checked with the same algorithm as above.

For example:

```txt
seeds: 79 14 55 13
```

Given the seed values from the first example, the above instead means that we have the following
seeds available as an input value:

* `79` - `92`
* `55` - `62`

This expands our source set and instead reveals that seed `82` results in a location of `46`
(given that original seed `13` is no longer available).

## This script

This script uses JavaScript and can be run with the following command:

```bash
npm i
node src/app.js -i input.txt
```

Which will run the above scenario on the given input file after installing any required packages
and will return the answer to part one of the above scenario.

To enable part two, run the above command with the `--mode range` flag set.
