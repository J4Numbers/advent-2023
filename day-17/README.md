# Day 17 - Clumsy Crucible

This is a simple script which operates on an input file with the following requirements:

* It several fixed width lines containing any of the characters `[0-9]`

## Part one

A file contains an ambient heat loss map that describes how much heat will be lost on entering
each tile. The task is to go from the starting position of top-left to the exit position of the
bottom right, minimising heat loss across the journey.

Additionally, the journey can only have a maximum of 3 tiles in one direction at once before
a change in direction must be performed.

For example, given the puzzle input:

```text
2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533
```

As heat loss is not triggered on the starting tile unless it is re-entered, the following is
one of the ways to minimise heat loss.

```text
2>>34^>>>1323
32v>>>35v5623
32552456v>>54
3446585845v52
4546657867v>6
14385987984v4
44578769877v6
36378779796v>
465496798688v
456467998645v
12246868655<v
25465488877v5
43226746555v>
```

Leading to a heat loss of only `102`.

## Part two

Alternatively, the minimum distance travelled in a straight line must be 4 (including stopping
at the end), and the maximum distance travelled in a straight mine must be 10.

In the above example, the updated min and max dists could follow this path to minimize heat loss:

```text
2>>>>>>>>1323
32154535v5623
32552456v4254
34465858v5452
45466578v>>>>
143859879845v
445787698776v
363787797965v
465496798688v
456467998645v
122468686556v
254654888773v
432267465553v
```

In the above example, an ultra crucible would incur the minimum possible heat loss of `94`.

## This script

This script uses Rust and can be run with the following command:

```bash
cargo run -- -i input.txt
```

This will answer part one as described above. To perform part two, add the `--ultra` flag.

```bash
cargo run -- -i input.txt --ultra
```
