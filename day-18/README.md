# Day 18 - Lavaduct Lagoon

This is a simple script which operates on an input file with the following requirements:

* Every line is in the format `[UDLR] [0-9]+ (#[a-f0-9]{6})`

## Part one

A file containing dig instructions has been provided. Starting from a 1x1 hole, the instructions
should be followed accordingly, as shown in the example below:

```text
R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)
```

The digger starts in a 1 meter cube hole in the ground. They then dig the specified number of 
meters up (`U`), down (`D`), left (`L`), or right (`R`), clearing full 1 meter cubes as they go.
The directions are given as seen from above, so if "up" were north, then "right" would be east,
and so on. Each trench is also listed with the color that the edge of the trench should be painted
as an RGB hexadecimal color code.

When viewed from above, the above example dig plan would result in the following loop of trench
(`#`) having been dug out from otherwise ground-level terrain (`.`):

```text
#######
#.....#
###...#
..#...#
..#...#
###.###
#...#..
##..###
.#....#
.######
```

At this point, the trench could contain `38` cubic meters of lava. However, this is just the edge
of the lagoon; the next step is to dig out the interior so that it is one meter deep as well:

```text
#######
#######
#######
..#####
..#####
#######
#####..
#######
.######
.######
```

Now, the lagoon can contain a much more respectable `62` cubic meters of lava. While the interior
is dug out, the edges are also painted according to the color codes in the dig plan.

## Part two

## This script

This script uses Rust and can be run with the following command:

```bash
cargo run -- -i input.txt
```

This will answer part one as described above.
