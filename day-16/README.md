# Day 16 - The floor will be lava

This is a simple script which operates on an input file with the following requirements:

* It several fixed width lines containing any of the characters `.-|\/`

## Part one

A file containing a map of mirrors and splitters is provided as a puzzle input. The mirrors
are `\/` while the splitters are `-|`. There's also empty space `.` which will fill the rest
of the map.

A beam of light enters at the top left of the map, heading right. Depending on the symbol it
encounters, it behaves differently:

* If it passes through empty space (`.`), it continues unchanged in the same direction.
* If it hits a mirror (`\` or `/`) - it is reflected at a 90 degree angle in the slant of the
  mirror (i.e. a rightward moving beam hitting `/` will be redirected upwards).
* If a beam encounters the thin end of a splitter (`|` or `-`), the beam treats the tile as
  empty space (i.e. a rightward moving beam hitting `-` would continue in the same direction).
* If a beam encounters the flat end of a splitter (`|` or `-`), the beam splits in two going
  in split directions from the origin of the splitter.

We would like to know the number of lit tiles in the map after the light has finished
travelling.

For example, given the puzzle input:

```text
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....
```

After the light enters at the top left tile travelling rightwards and finishes moving around
the map (even if it's found an infinite loop), then the map will look like this instead:

```text
>|<<<\....
|v-.\^....
.v...|->>>
.v...v^.|.
.v...v^...
.v...v^..\
.v../2\\..
<->-/vv|..
.|<<<2-|.\
.v//.|.v..
```

If we say that a tile which has been lit is 'energised' (`#`), then the map will look like
this:

```text
######....
.#...#....
.#...#####
.#...##...
.#...##...
.#...##...
.#..####..
########..
.#######..
.#...#.#..
```

With `46` energised tiles.

## Part two

Alternatively, the beam can start at any edge tile in any direction. So, the beam can travel
downwards from any top tile, upwards from any bottom tile, leftwards from any right tile, and
rightwards from any left tile.

We should look to see which entrance energises the most tiles, and how many tiles are energised
in that specific configuration.

In the above example, this can be achieved by starting the beam in the fourth tile from the
left in the top row:

```text
.|<2<\....
|v-v\^....
.v.v.|->>>
.v.v.v^.|.
.v.v.v^...
.v.v.v^..\
.v.v/2\\..
<-2-/vv|..
.|<<<2-|.\
.v//.|.v..
```

Using this configuration, `51` tiles are energized:

```text
.#####....
.#.#.#....
.#.#.#####
.#.#.##...
.#.#.##...
.#.#.##...
.#.#####..
########..
.#######..
.#...#.#..
```

## This script

This script uses Rust and can be run with the following command:

```bash
cargo run -- -i input.txt
```

This will answer part one as described above. To modify the program to find the best location
to start the beam from, add the `--explore` flag.

```bash
cargo run -- -i input.txt --explore
```
