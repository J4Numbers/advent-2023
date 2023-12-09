# Day 7 - Mirage Maintenance

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Every line is in the format `( ?[-0-9]+)+`

## Part one

A file contains several sequences of integers. The task is to predict the next value in the
sequence using a set algorithm, then add all the next numbers in each sequence together.

The algorithm goes as follows with the following example:

```text
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
```

Start by making a new sequence from the difference at each step in the sequence. If that
sequence is not all zeroes, repeat this process using the sequence you just generated as the
input sequence. Once all of the values in your latest sequence are zeroes, you can extrapolate
what the next value of the original history should be.

In the above dataset, the first history is `0 3 6 9 12 15`. Because the values increase by `3`
each step, the first sequence of differences that you generate will be `3 3 3 3 3`. Since these
values aren't all zero, repeat the process: the values differ by `0` at each step, so the next
sequence is `0 0 0 0`. This means you have enough information to extrapolate the history!

> Note that each subsequent sequence has one fewer value than the input sequence because at
> each step it considers two numbers from the input.

Visually, these sequences can be arranged like this:

```text
0   3   6   9  12  15
  3   3   3   3   3
    0   0   0   0
```

To extrapolate, start by adding a new zero to the end of your list of zeroes; because the zeroes
represent differences between the two values above them, this also means there is now a
placeholder in every sequence above it:

```text
0   3   6   9  12  15   B
  3   3   3   3   3   A
    0   0   0   0   0
```

You can then start filling in placeholders from the bottom up. `A` needs to be the result of
increasing `3` (the value to its left) by `0` (the value below it); this means `A` must be `3`:

```text
0   3   6   9  12  15   B
  3   3   3   3   3   3
    0   0   0   0   0
```

Finally, you can fill in `B`, which needs to be the result of increasing `15` (the value to its
left) by `3` (the value below it), or `18`:

```text
0   3   6   9  12  15  18
  3   3   3   3   3   3
    0   0   0   0   0
```

So, the next value of the first history is `18`.

Finding all-zero differences for the second history requires an additional sequence:

```text
1   3   6  10  15  21
  2   3   4   5   6
    1   1   1   1
      0   0   0
```

Then, following the same process as before, work out the next value in each sequence from the
bottom up:

```text
1   3   6  10  15  21  28
  2   3   4   5   6   7
    1   1   1   1   1
      0   0   0   0
```

So, the next value of the second history is `28`.

The third history requires even more sequences, but its next value can be found the same way:

```text
10  13  16  21  30  45  68
   3   3   5   9  15  23
     0   2   4   6   8
       2   2   2   2
         0   0   0
```

So, the next value of the third history is `68`.

If you find the next value for each history in this example and add them together, you get `114`.

## Part two

Part two is to do exactly the same, but in reverse instead - so instead of a prediction into the
future, it is a divination of the past. Using the example above:

```text
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
```

We can do the following for the third example:

```text
5  10  13  16  21  30  45
  5   3   3   5   9  15
   -2   0   2   4   6
      2   2   2   2
        0   0   0
```

To reveal that the left-most value is `5`. Doing the same for the other two lines reveals a value
of `0` for the second history, and `-3` for the first.

Summing these up together returns a final value of `2`.

## This script

This script uses Python and can be run with the following command:

```bash
python step_predict.py -i input.txt
```

This will answer part one of the question as described above. To return the value for part two,
append the `--mode previous` flag to your run to select that. i.e.

```bash
python step_predict.py -i input.txt --mode previous
```
