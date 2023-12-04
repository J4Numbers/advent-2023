# Day 4 - Scratchcards

This is a simple script which operates on an input file with the following requirements:

* It is made up of between 1 and many lines
* Each line contains a fixed length string consisting of a card ID and two sets of numbers
  separated by a pipe (`|`)

## Part one

Two sets of numbers are given - the first set is a list of 'winning' numbers, and the second set
is a list of 'available' numbers. Any available number which is a 'winner' is counted as a match.
The point value of a scratch card is equal to `2^(match count - 1)` - or 0 if there are no matches.

The solution to the puzzle is the sum of all points across all scratchcards.

For example:

```txt
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
```

In the above example, card 1 has five winning numbers (`41`, `48`, `83`, `86`, and `17`) and eight
available numbers (`83`, `86`, `6`, `31`, `17`, `9`, `48`, and `53`). Of the available numbers, four
of them (`48`, `83`, `17`, and `86`) are winning numbers! That means card 1 is worth `8` points (1 for
the first match, then doubled three times for each of the three matches after the first).

* Card 2 has two winning numbers (`32` and `61`), so it is worth `2` points.
* Card 3 has two winning numbers (`1` and `21`), so it is worth `2` points.
* Card 4 has one winning number (`84`), so it is worth `1` point.
* Card 5 has no winning numbers, so it is worth no points.
* Card 6 has no winning numbers, so it is worth no points.

So, in this example, the Elf's pile of scratchcards is worth `13` points.

## Part two

Alternatively, the scratchcards generate copies of their subsequent cards on matching numbers. For
example, if card 10 had 3 matching numbers, a copy of 11, 12, and 13 would be generated. These copies
also generate copies if they have a match (i.e. if card 11 had one matching number, then there would
be 1 original and 2 copies of card 12).

> We cannot copy card IDs that do not exist past the end of the table.

We would like to find out the total number of scratchcards (copies and originals).

For example:

```txt
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
```

* Card 1 has four matching numbers, so you win one copy each of the next four cards: cards 2, 3,
  4, and 5.
* Your original card 2 has two matching numbers, so you win one copy each of cards 3 and 4.
* Your copy of card 2 also wins one copy each of cards 3 and 4.
* Your four instances of card 3 (one original and three copies) have two matching numbers, so you
  win four copies each of cards 4 and 5.
* Your eight instances of card 4 (one original and seven copies) have one matching number, so you
  win eight copies of card 5.
* Your fourteen instances of card 5 (one original and thirteen copies) have no matching numbers
  and win no more cards.
* Your one instance of card 6 (one original) has no matching numbers and wins no more cards.

Once all of the originals and copies have been processed, you end up with 1 instance of card 1, 2
instances of card 2, 4 instances of card 3, 8 instances of card 4, 14 instances of card 5, and 1
instance of card 6. In total, this example pile of scratchcards causes you to ultimately have `30`
scratchcards!

## This script

This script uses JavaScript and can be run with the following command:

```bash
npm i
node src/app.js -i input.txt
```

Which will run the above scenario on the given input file after installing any required packages
and will return the answer to both of the above questions.
