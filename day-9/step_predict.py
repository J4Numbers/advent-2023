"""
A script that performs the actions required for Day 7 of Advent of Code, 2023.
The script takes in an input file parameter, consisting of many lines made up
of a 5 character string, and a point value.

These character strings are sorted on a given rule, and a final result is
then printed according to the point value and their ranks.

This script operates in two modes - 'camel', which is the original mode of
poker using jacks, or 'jokers', where J is jokers, which can fill in for any
other card at the cost of losing all high card battles.
"""
import argparse
import re

PARSER_DESC = """Calculate the summation of the next value in a series of
sequences"""
VALUE_REGEX = re.compile("([-0-9]+)", re.IGNORECASE)

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["next", "previous"], default="next",
    help="The mode of operation for the program")
parser.add_argument("--debug", action="store_true",
                    help="Enable debug logging")

args = parser.parse_args()


def debug_line(line_to_print):
    """
    Debug a given line if we have debut output enabled for this program.
    Otherwise, do nothing.
    :param line_to_print: the line to debug out to console
    """
    if args.debug:
        print(line_to_print)


def predict_next(sequence_breakdown):
    """
    Predict the next value in a sequence of numbers which has already been
    broken down into its stepped values.
    :param sequence_breakdown: A broken down sequence of sequences, each with
                               the stepped value required to reach the next
                               value in the above sequence.
    :return: The next value in a sequence according to the stepped parent
    sequences
    """
    seq_iteration = len(sequence_breakdown) - 1
    addition_marker = 0
    while seq_iteration >= 0:
        last_index = len(sequence_breakdown[seq_iteration]) - 1
        addition_marker = (sequence_breakdown[seq_iteration][last_index]
                           + addition_marker)
        seq_iteration -= 1

    debug_line(f"Sequence starting with {sequence_breakdown[0][0]} predicts "
               f"next value of {addition_marker}")
    return addition_marker


def predict_previous(sequence_breakdown):
    """
    Predict the previous value in a sequence of numbers which has already been
    broken down into its stepped values.
    :param sequence_breakdown: A broken down sequence of sequences, each with
                               the stepped value required to reach the next
                               value in the above sequence.
    :return: The previous value in a sequence according to the stepped parent
    sequences
    """
    seq_iteration = len(sequence_breakdown) - 1
    addition_marker = 0
    while seq_iteration >= 0:
        addition_marker = (sequence_breakdown[seq_iteration][0]
                           - addition_marker)
        seq_iteration -= 1

    debug_line(f"Sequence starting with {sequence_breakdown[0][0]} predicts "
               f"previous value of {addition_marker}")
    return addition_marker


def break_sequence(sequence):
    """
    Take a sequence and break it down into individual layers so that the
    difference in each step is tracked.
    :param sequence: the first hand of five cards to compare
    :return: An array of sequence levels from 0 as the original to x where
    all entries in that sequence are 0
    """
    sequence_breakdown = [sequence]
    working_sequence = sequence
    while len(list(filter(lambda s: s != 0, working_sequence))) > 0:
        prev_val = None
        next_seq = []
        for step in working_sequence:
            if prev_val is not None:
                next_seq.append(step - prev_val)
            prev_val = step
        sequence_breakdown.append(next_seq)
        working_sequence = next_seq

    debug_line(f"Sequence starting with {sequence[0]} broken down into "
               f"{len(sequence_breakdown)} sequences")
    return sequence_breakdown


# MAIN CODE STARTS HERE

# Set up some base variables
summation = 0
predictions = []
sequence_collection = []

# For each line in the file, we extract the list of sequence numbers and store
# them inside an array
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            sequence_numbers = VALUE_REGEX.findall(line)
            if len(sequence_numbers) > 0:
                sequence_collection.append(
                        [int(number) for number in sequence_numbers])

# Break down all the sequences into their stepped parts
broken_sequences = [break_sequence(seq) for seq in sequence_collection]

# Depending on the mode of operation, either predict the next value in the
# sequence or the previous value
if args.mode == "next":
    predictions = [predict_next(seq) for seq in broken_sequences]
else:
    predictions = [predict_previous(seq) for seq in broken_sequences]

# Sum up all of our predictions...
for prediction in predictions:
    summation += prediction

# Print out the final result
print(summation)
