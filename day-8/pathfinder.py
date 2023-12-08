"""
A script that performs the actions required for Day 8 of Advent of Code, 2023.
The script takes in an input file parameter, consisting of many lines made up
of a 3 character node, and two 3 character options.

These nodes form a set of nodes and edges that can be used to run from AAA to
ZZZ in the fewest number of directions (where the leftmost direction is L and
the rightmost direction is R).

This script operates in two modes - 'seek', which is the original mode of
simply following directions, or 'jokers', where J is jokers, which can fill in for any
other card at the cost of losing all high card battles.
"""
import argparse
import re

from functools import reduce

PARSER_DESC = "Calculate the directions required to travel from AAA to ZZZ"
DIRECTIONS_REGEX = "^[LR]+$"
NODE_REGEX = "^([a-z]{3}) = \(([a-z]{3}), ([a-z]{3})\)$"
START_NODE = "AAA"
END_NODE = "ZZZ"

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["seek"], default="seek",
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


def flat_directions(so_far, current_node):
    return so_far + current_node["direction"]


def history_contains(history, test_node):
    """
    Test whether a given history array already contains a test node as
    visited.
    :param history: the complete journey from start to finish we have
                    taken to reach this point.
    :param test_node: the node to test the presence of in our history.
    :returns: True if the node exists in our history
    """
    test_presence = list(filter(lambda point: point["node"] == test_node, history))
    return len(test_presence) > 0


def split_on_node(split_point, target_node, map_of_nodes, history):
    """
    Split a given node in the node map to follow both paths and construct an
    eventual path to the target node. The rules are: no overlapping history
    as it means backtracking, and the final node is it.

    :param split_point: the node that we are splitting down the middle.
    :param target_node: the final node we are aiming to reach.
    :param map_of_nodes: the map of all nodes and edges that connect them.
    :param history: the journey we have taken to reach this split node.
    :return: A list of separate histories that are either final or terminated
    """
    histories = []
    options = map_of_nodes[split_point]
    dir_so_far = reduce(flat_directions, history, "");

    if split_point == target_node:
        history_point = {
            "node": split_point,
            "final": True
        }
        histories.append([*history, history_point])
        debug_line(f"Target node {target_node} found after {len(history) + 1} hops - {dir_so_far}")

    else:
        avail_paths = 0

        if not history_contains(history, options[0]) and options[0] not in redlisted_nodes:
            hist_left = {
                "node": split_point,
                "direction": "L",
                "next_node": options[0]
            }
            split_hist_l = split_on_node(options[0], target_node, map_of_nodes, [*history, hist_left])
            histories += split_hist_l
            if len(split_hist_l) == 0:
                debug_line(f"Redlisting {options[0]} as no paths to victory were found")
                redlisted_nodes.append(options[0])
        else:
            debug_line(f"Repeat L node {options[0]} found after {len(history)} hops - {dir_so_far}")

        if not history_contains(history, options[1]) and options[1] not in redlisted_nodes:
            hist_right = {
                "node": split_point,
                "direction": "R",
                "next_node": options[1]
            }
            split_hist_r = split_on_node(options[1], target_node, map_of_nodes, [*history, hist_right])
            histories += split_hist_r
            if len(split_hist_r) == 0:
                debug_line(f"Redlisting {options[1]} as no paths to victory were found")
                redlisted_nodes.append(options[1])
        else:
            debug_line(f"Repeat R node {options[1]} found after {len(history)} hops - {dir_so_far}")

    return histories


# MAIN CODE STARTS HERE

# Set up some base variables
node_map = {}
travel_history = []
redlisted_nodes = []

# For each line in the file, we extract the hand of cards and their points
# that are associated with them and add them onto a working tuple array.
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            line_discovery = re.search(NODE_REGEX, line, flags=re.IGNORECASE)
            if line_discovery is not None:
                node_map[line_discovery.group(1)] = (
                        line_discovery.group(2), line_discovery.group(3))

# Explode all possible paths from the start to the target node
paths = split_on_node(START_NODE, END_NODE, node_map, [])

# Reduce all paths down to the winning path
winning_path = reduce(
    lambda winner, challenger: challenger if not winner else (winner if len(winner) < len(challenger) else challenger),
    paths, None)

print(f"{len(paths)} different paths from {START_NODE} to {END_NODE} were found")
# Print out the final result
print(winning_path)
print(len(winning_path))
