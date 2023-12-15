# Day 13 - Lens Library

This is a simple script which operates on an input file with the following requirements:

* It is a single line of character sequence broken up by commas

## Part one

A file contains a series of instructions. Those instructions should be converted into a 
hash value between 0 and 255 for each instruction (as seperated by commas). The algorithm
for calculating the hash is as follows:

* Starting from a current value of zero (0)
* Determine the ASCII code for the current character of the string.
* Increase the current value by the ASCII code you just determined.
* Set the current value to itself multiplied by 17.
* Set the current value to the remainder of dividing itself by 256.

So, to find the result of running the HASH algorithm on the string `HASH`:

* The current value starts at `0`.
* The first character is `H`; its ASCII code is `72`.
* The current value increases to `72`.
* The current value is multiplied by `17` to become `1224`.
* The current value becomes `200` (the remainder of `1224` divided by `256`).
* The next character is `A;` its ASCII code is `65`.
* The current value increases to `265`.
* The current value is multiplied by `17` to become `4505`.
* The current value becomes `153` (the remainder of `4505` divided by `256`).
* The next character is `S`; its ASCII code is `83`.
* The current value increases to `236`.
* The current value is multiplied by `17` to become `4012`.
* The current value becomes `172` (the remainder of `4012` divided by `256`).
* The next character is `H`; its ASCII code is `72`.
* The current value increases to `244`.
* The current value is multiplied by `17` to become `4148`.
* The current value becomes `52 (`the remainder of `4148` divided by `256`).

So, the result of running the HASH algorithm on the string HASH is `52`.

## Part two

Alternatively, each item in the list falls into two categories:

* `[a-z]+=[1-9]` - A mapping of a label to a value between 1-9
* `[a-z]+-` - Shorthand to remove a given label from the system

Following this list of instructions will put a series of label->value pairings into
several boxes that are mapped to by the hash function from part one. The label part is
what directly maps the two.

Instructions are carried out sequentially, meaning that a removal before a label has been
added will do nothing, and that setting a value when the label is already mapped will
overwrite the original value without changing its position in the box.

Here is the contents of every box after each step in the example initialization sequence
above:

```text
After "rn=1":
Box 0: [rn 1]

After "cm-":
Box 0: [rn 1]

After "qp=3":
Box 0: [rn 1]
Box 1: [qp 3]

After "cm=2":
Box 0: [rn 1] [cm 2]
Box 1: [qp 3]

After "qp-":
Box 0: [rn 1] [cm 2]

After "pc=4":
Box 0: [rn 1] [cm 2]
Box 3: [pc 4]

After "ot=9":
Box 0: [rn 1] [cm 2]
Box 3: [pc 4] [ot 9]

After "ab=5":
Box 0: [rn 1] [cm 2]
Box 3: [pc 4] [ot 9] [ab 5]

After "pc-":
Box 0: [rn 1] [cm 2]
Box 3: [ot 9] [ab 5]

After "pc=6":
Box 0: [rn 1] [cm 2]
Box 3: [ot 9] [ab 5] [pc 6]

After "ot=7":
Box 0: [rn 1] [cm 2]
Box 3: [ot 7] [ab 5] [pc 6]
```

All `256` boxes are always present; only the boxes that contain any lenses are shown here.
Within each box, lenses are listed from front to back; each lens is shown as its label and
focal length in square brackets.

To confirm that all of the lenses are installed correctly, add up the focusing power of
all of the lenses. The focusing power of a single lens is the result of multiplying
together:

* One plus the box number of the lens in question.
* The slot number of the lens within the box: 1 for the first lens, 2 for the second lens, 
  and so on.
* The focal length of the lens.

At the end of the above example, the focusing power of each lens is as follows:

* `rn`: `1` (box 0) `* 1` (first slot) `* 1 `(focal length) `= 1`
* `cm`: `1` (box 0) `* 2` (second slot) `* 2` (focal length) `= 4`
* `ot`: `4` (box 3) `* 1` (first slot) `* 7` (focal length) `= 28`
* `ab`: `4` (box 3) `* 2` (second slot) `* 5` (focal length) `= 40`
* `pc`: `4` (box 3) `* 3` (third slot) `* 6` (focal length) `= 72`

So, the above example ends up with a total focusing power of `145`.

## This script

This script uses Go and can be run with the following command:

```bash
go run . -i input.txt
```

This will answer part one as described above. To enable the focusing mode, add the `--focus`
flag to the command.

```bash
go run . -i input.txt --focus
```
