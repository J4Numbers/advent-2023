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

## This script

This script uses Python and can be run with the following command:

```bash
python pathfinder.py -i input.txt
```

This will answer part one of the question as described above.
