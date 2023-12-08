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
NODE_REGEX = "^([0-9a-z]{3}) = \(([0-9a-z]{3}), ([0-9a-z]{3})\)$"
START_NODE = "AAA"
END_NODE = "ZZZ"

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["seek", "haunt"], default="seek",
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


def follow_path(start_nodes, target_nodes, map_of_nodes, directions):
    common_candidate = 0
    voters_agree = []
    last_pos = {}

    while len(voters_agree) < len(start_nodes):
        for s_node in start_nodes:
            if s_node not in voters_agree:
                debug_line(f"{s_node} :: Not found in voter agreement... Attempting to match candidate of {common_candidate}")
                if s_node in last_pos:
                    dir_idx = last_pos[s_node]["idx"]
                    working_node = last_pos[s_node]["node"]
                    hops = last_pos[s_node]["hops"]
                else:
                    dir_idx = 0
                    working_node = s_node
                    hops = 0
                debug_line(f"{s_node} :: Starting from {working_node} - IDX {dir_idx} - {hops}")

                while working_node not in target_nodes or hops < common_candidate:
                    hops += 1
                    prev_node = working_node
                    working_node = map_of_nodes[working_node][0 if directions[dir_idx] == "L" else 1]
                    # debug_line(f"{s_node} :: Going {directions[dir_idx]} from {prev_node} to {working_node} - {hops}")
                    dir_idx = (dir_idx + 1) % len(directions)

                if common_candidate == hops:
                    voters_agree.append(s_node)
                    debug_line(f"{s_node} :: Matched candidate {common_candidate} on {working_node} - {len(start_nodes) - len(voters_agree)} voters left")
                else:
                    voters_agree = [s_node]
                    common_candidate = hops
                    debug_line(f"{s_node} :: New candidate {common_candidate} on {working_node} - {len(start_nodes) - len(voters_agree)} voters left")

                last_pos[s_node] = {
                    "node": working_node,
                    "idx": dir_idx,
                    "hops": hops
                }

    print(last_pos)
    return common_candidate


# MAIN CODE STARTS HERE

# Set up some base variables
dir_line = []
node_map = {}

# For each line in the file, we extract the hand of cards and their points
# that are associated with them and add them onto a working tuple array.
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            if re.search(DIRECTIONS_REGEX, line, flags=re.IGNORECASE):
                dir_line = list(line.strip())
            line_discovery = re.search(NODE_REGEX, line, flags=re.IGNORECASE)
            if line_discovery is not None:
                node_map[line_discovery.group(1)] = (
                        line_discovery.group(2), line_discovery.group(3))

start_nodes = []
end_nodes = []

print(node_map)

if args.mode == "haunt":
    start_nodes = list(filter(lambda n: n.endswith("A"), node_map.keys()))
    end_nodes = list(filter(lambda n: n.endswith("Z"), node_map.keys()))
else:
    start_nodes = [START_NODE]
    end_nodes = [END_NODE]

path = follow_path(start_nodes, end_nodes, node_map, dir_line)
print(path)

