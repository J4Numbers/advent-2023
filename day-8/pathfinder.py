"""
A script that performs the actions required for Day 8 of Advent of Code, 2023.
The script takes in an input file parameter, consisting of many lines made up
of a 3 character node, and two 3 character options.

These nodes form a set of nodes and edges that can be used to run from AAA to
ZZZ in the fewest number of directions (where the leftmost direction is L and
the rightmost direction is R).

This script operates in two modes - 'seek', which is the original mode of
simply following directions, or 'haunt', which treats all nodes ending in A as
starting nodes, and all nodes ending in Z as target nodes. Searching for the
first common point.
"""
import argparse
import re

from math import lcm

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
    """
    Given a list of start nodes and a list of potential target nodes, follow
    a given set of directions across a given graph to calculate the number of
    steps required to get from an individual starting node to any target node
    within our inputs.

    If we run out of steps in the directions to follow, we start again from
    scratch.
    :param start_nodes: A list of starting nodes that we will begin our search
                        from.
    :param target_nodes: A list of target nodes that any node can count as
                         their destination.
    :param map_of_nodes: A complete dictionary of nodes which map onto two
                         edges which describes the possible paths within this
                         system.
    :param directions: A list of L/R directions that describe which of the
                       edges in a node graph we will take at any one time.
    :returns: A map of starting nodes against how long it took to reach any
    target node within our system, along with what that target node
    specifically was.
    """
    first_multiples = {}

    for s_node in start_nodes:
        dir_idx = 0
        working_node = s_node
        hops = 0
        debug_line(f"{s_node} :: Starting from {working_node} - IDX {dir_idx} - {hops}")

        while working_node not in target_nodes:
            hops += 1
            prev_node = working_node
            working_node = map_of_nodes[working_node][0 if directions[dir_idx] == "L" else 1]
            debug_line(f"{s_node} :: Going {directions[dir_idx]} from {prev_node} to {working_node} - {hops}")
            dir_idx = (dir_idx + 1) % len(directions)

        first_multiples[s_node] = {
            "node": working_node,
            "idx": dir_idx,
            "hops": hops
        }

    return first_multiples


# MAIN CODE STARTS HERE

# Set up some base variables
dir_line = []
node_map = {}

# For each line in the file, we extract the node and its two edges, along
# with the one line which describes our allowed directions when traversing
# the graph
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            if re.search(DIRECTIONS_REGEX, line, flags=re.IGNORECASE):
                dir_line = list(line.strip())
            line_discovery = re.search(NODE_REGEX, line, flags=re.IGNORECASE)
            if line_discovery is not None:
                node_map[line_discovery.group(1)] = (
                        line_discovery.group(2), line_discovery.group(3))

# Some working variables for the start and end locations
starting_nodes = []
ending_nodes = []

# If we're haunting, we're going over all nodes that end in A with all nodes
# that end in Z as a target. Otherwise it's simply AAA -> ZZZ
if args.mode == "haunt":
    starting_nodes = list(filter(lambda n: n.endswith("A"), node_map.keys()))
    ending_nodes = list(filter(lambda n: n.endswith("Z"), node_map.keys()))
else:
    starting_nodes = [START_NODE]
    ending_nodes = [END_NODE]

# Follow the paths
paths = follow_path(starting_nodes, ending_nodes, node_map, dir_line)

# Debug out what we found in the paths
for start_node, path in paths.items():
    debug_line(f"{start_node} :: Reached {path['node']} after {path['hops']} hops")

# Print the final result
print(lcm(*list(map(lambda p: p["hops"], paths.values()))))
