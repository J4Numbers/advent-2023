"""
A script that performs the actions required for Day 6 of Advent of Code, 2023.
The script takes in an input file parameter, consisting of one line with a
time series, and one line with a distance series (which are then twinned
together into a tuple).

For these tuples, the total number of permutations is calculated where we can
exceed the distance given a formula.

This script operates in two modes - 'collect' all time series into a set of
tuples, or 'join' all series into one value to be tupled together.
"""
import argparse
import re

PARSER_DESC = """
Calculate the number of permutations that can be found on going the furthest
distance within a fixed time
"""
TIME_REGEX = "^time: [0-9 ]+$"
DIST_REGEX = "^distance: [0-9 ]+$"

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["collect", "join"],
    help="Mode of operation", default="collect")
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


def calculate_permutations(time_allowed, base_distance_covered):
    """
    Calculate the number of permutations that exist given a bounded time
    and the max distance covered within that time as a score to beat. This
    is done by splitting the time allowed into time charging and time
    running, then multiplying them together while we step through the
    available values.

    :param time_allowed: the total time allowed to gain distance.
    :param base_distance_covered: the base value that we have to beat in
    a given permutation.
    """
    over_count = 0
    for time_charging in range(time_allowed):
        time_running = time_allowed - time_charging
        if (time_running * time_charging) > base_distance_covered:
            over_count += 1
    return over_count


# MAIN CODE STARTS HERE

# Set up some base variables
total_permutations = 1
time_list = []
distance_list = []
td_tuples = []

# For each line in the file, we ask if is a time series or a distance series,
# then assign that series accordingly to the working variables above. If the
# mode is set to join, then the series is joined together into a single value
# instead.
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            if re.search(TIME_REGEX, line, flags=re.IGNORECASE) is not None:
                if args.mode == "collect":
                    time_list = re.findall("[0-9]+", line)
                elif args.mode == "join":
                    time_list = ["".join(re.findall("[0-9]+", line))]
            if re.search(DIST_REGEX, line, flags=re.IGNORECASE) is not None:
                if args.mode == "collect":
                    distance_list = re.findall("[0-9]+", line)
                elif args.mode == "join":
                    distance_list = ["".join(re.findall("[0-9]+", line))]

    # Twin all of our discovered tuples together
    for x in range(len(time_list)):
        td_tuples.append((int(time_list[x]), int(distance_list[x])))

    # For each tuple, calculate the number of permutations that exist where we
    # can gain a higher distance in the provided time. Then multiply that
    # number of permutations against the existing permutation factor for our
    # puzzle answer.
    for td in td_tuples:
        permutations = calculate_permutations(td[0], td[1])
        total_permutations *= permutations
        format_str = "Discovered max distance of {} mm covered in {} ms has {} permutations to increase distance"
        debug_line(format_str.format(td[1], td[0], permutations))

# Print the answer to the puzzle
print(total_permutations)
