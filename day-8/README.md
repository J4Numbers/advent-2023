# Day 8 - Haunted wasteland

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* One line is made up of `[LR]+` in a set of directions
* All other lines are made up of `[A-Z]{3} = \([A-Z]{3}, [A-Z]{3}\)`

## Part one

A file containing a set of nodes and edges has been provided. The starting node is always
`AAA` and the ending node is always `ZZZ`. We must follow the list of directions ad
nauseum until we reach the ending node.

For example:

```txt
RL

AAA = (BBB, CCC)
BBB = (DDD, EEE)
CCC = (ZZZ, GGG)
DDD = (DDD, DDD)
EEE = (EEE, EEE)
GGG = (GGG, GGG)
ZZZ = (ZZZ, ZZZ)
```

Starting with `AAA`, you need to look up the next element based on the next left/right instruction
in your input. In this example, start with `AAA` and go right (`R`) by choosing the right element
of `AAA`, `CCC`. Then, `L` means to choose the left element of `CCC`, `ZZZ`. By following the
left/right instructions, you reach `ZZZ` in 2 steps.

Of course, you might not find `ZZZ` right away. If you run out of left/right instructions, repeat
the whole sequence of instructions as necessary: `RL` really means `RLRLRLRLRLRLRLRL...` and so on.
For example, here is a situation that takes 6 steps to reach `ZZZ`:

```txt
LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)
```

## Part two

Alternatively, instead of just starting from `AAA`, we instead start simultaneously from all nodes
ending in `A`. The challenge now is to find the number of steps required before all nodes are
ending in `Z`.

> If some nodes are ending in `Z`, but not all of them, then we continue on all nodes until we
> find a better match

For example:

```txt
LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)
```

Here, there are two starting nodes, `11A` and `22A` (because they both end with A). As you follow
each left/right instruction, use that instruction to simultaneously navigate away from both nodes
you're currently on. Repeat this process until all of the nodes you're currently on end with `Z`.
In this example, you would proceed as follows:

* Step 0: You are at `11A` and `22A`.
* Step 1: You choose all of the left paths, leading you to `11B` and `22B`.
* Step 2: You choose all of the right paths, leading you to `11Z` and `22C`.
* Step 3: You choose all of the left paths, leading you to `11B` and `22Z`.
* Step 4: You choose all of the right paths, leading you to `11Z` and `22B`.
* Step 5: You choose all of the left paths, leading you to `11B` and `22C`.
* Step 6: You choose all of the right paths, leading you to `11Z` and `22Z`.

So, in this example, you end up entirely on nodes that end in `Z` after 6 steps.

## This script

This script uses Python and can be run with the following command:

```bash
python pathfinder.py -i input.txt
```

This will answer part one of the question as described above. To switch to answering part two, add
the `--mode haunt` parameter to the program. i.e.

```bash
python pathfinder.py -i input.txt --mode haunt
```
