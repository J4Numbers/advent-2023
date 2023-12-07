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

from functools import reduce, cmp_to_key

PARSER_DESC = """Calculate the rank power of a given set of card hands and
their point values
"""
CARD_REGEX = "^([akqjt0-9]{5}) +([0-9]+)$"

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["camel", "jokers"],
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


def collect_hand(ongoing, current_card):
    """
    A reduce function to count the number of a given card within a hand.
    :param ongoing: the ongoing reduction of card numbers
    :param current_card: the current card that we are inspecting
    :return: the updated reduction of card numbers in a hand
    """
    if current_card not in ongoing:
        ongoing[current_card] = 0
    ongoing[current_card] += 1
    return ongoing


def rank_hand(hand):
    """
    Given a hand of five cards, rank that hand into its type starting
    from 1 as a high card to 7 as a five of a kind.
    :param hand: a hand of five cards to categorise by type
    :return: A number value for the rank of the hand, from 1 to 7
    """
    split_str = reduce(collect_hand, list(hand), {})
    joker_count = 0
    if (args.mode == "jokers" and "J" in split_str):
        joker_count = split_str["J"]
        split_str["J"] = 0
    count = sorted(split_str.values(), reverse=True)
    rank = 1
    count[0] += joker_count
    if count[0] > 3:
        rank = count[0] + 2
    if count[0] == 3:
        if count[1] == 2:
            rank = 5
        else:
            rank = 4
    if count[0] == 2:
        if count[1] == 2:
            rank = 3
        else:
            rank = 2
    return rank


def convert_card(card):
    """
    Convert a given card ([0-9TJQKA]) to a score value relative to other
    cards, running from 1 for jokers (if enabled) to 14 for aces.
    :param card: the card to convert
    :return: the score value of this card
    """
    ret_val = 0
    if re.match("[2-9]", card):
        ret_val = int(card)
    if card == "T":
        ret_val = 10
    if card == "J":
        ret_val = 1 if args.mode == "jokers" else 11
    if card == "Q":
        ret_val = 12
    if card == "K":
        ret_val = 13
    if card == "A":
        ret_val = 14
    return ret_val


def hand_walk(hand_a, hand_b):
    """
    Walk through two card hands to consider which one has the higher value.
    The actual conversion of a card to a value is offloaded to another
    function.
    :param hand_a: the first hand of five cards to compare
    :param hand_b: the second hand of five cards to compare
    :return -1, 0, or 1 depending on the relative order of the hands
    """
    cards_a = list(hand_a)
    cards_b = list(hand_b)
    ret_val = 0
    for idx, val in enumerate(cards_a):
        if val == cards_b[idx]:
            continue
        ret_val = convert_card(val) - convert_card(cards_b[idx])
        break
    return ret_val


def compare_hands(hand_a, hand_b):
    """
    Compare a given pair of hands together using the following sorting
    algorithm:
    * Five of a kind trumps four of a kind
    * Four of a kind trumps a full house
    * A full house trumps three of a kind
    * Three of a kind trumps two pairs
    * Two pairs trump one pair
    * One pair trumps a high card
    * High card runs AKQJT98765432

    If two hands have the same type, a second ordering rule takes
    effect. We compare the first card in each hand, and if these cards
    are different, the hand with the stronger first card is considered
    stronger according to our high card rules. If the first card in each
    hand have the same label, however, then we move on to consider the
    second card in each hand, then the third, fourth, and fifth, until we
    find a difference.

    :param hand_a: the first hand to compare.
    :param hand_b: the second hand to compare.
    :return: -1, 0, or 1 depending on the comparison result
    """
    e_hand_a = rank_hand(hand_a[0])
    e_hand_b = rank_hand(hand_b[0])

    ret_val = e_hand_a - e_hand_b

    if ret_val == 0:
        ret_val = hand_walk(hand_a[0], hand_b[0])

    return ret_val


# MAIN CODE STARTS HERE

# Set up some base variables
rank_power = 0
hand_tuples = []

# For each line in the file, we extract the hand of cards and their points
# that are associated with them and add them onto a working tuple array.
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            line_discovery = re.search(CARD_REGEX, line, flags=re.IGNORECASE)
            if line_discovery is not None:
                hand_tuples.append(
                        (line_discovery.group(1), int(line_discovery.group(2))))

# Sort the tuples based on hand value
hand_tuples.sort(key=cmp_to_key(compare_hands))

# And derive the rank power based on relative rank in the array and point
# value.
for rank0, tpl in enumerate(hand_tuples):
    rank_power += (rank0 + 1) * tpl[1]
    debug_line(f"Rank {rank0 + 1} - Hand {tpl[0]} - Value {tpl[1]} - Running total {rank_power}")

# Print out the final result
print(rank_power)
