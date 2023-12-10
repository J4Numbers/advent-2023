"""
A script that performs the actions required for Day 10 of Advent of Code, 2023.
The script takes in an input file parameter, consisting of many lines made up
of a fixed character maze.

These character strings are made up of straight pieces, bends, ground area, and
one position which is the start for a loop within the maze.
"""
import sys
import argparse
import re

PARSER_DESC = """Calculate the anti-point within a maze which is the furthest
position from a given start, along with the flood value"""
MAP_REGEX = re.compile("([-.|JLF7S]+)", re.IGNORECASE)

parser = argparse.ArgumentParser(description=PARSER_DESC)

parser.add_argument(
    "--input", "-i",
    help="The input file to provide to the script", required=True)
parser.add_argument(
    "--mode", choices=["counter"], default="counter",
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


def get_map_point(mapped_array, pos):
    """
    Return a given coordinate in a two-dimensional array
    :param mapped_array: The 2D array that holds our map
    :param pos: The specific (y,x) position that we're looking up
    :return: The value of pos within mapped_array
    """
    return mapped_array[pos[0]][pos[1]]


def add_coords(a, b):
    """
    Add two co-ordinate tuples together (y,x)
    :param a: First co-ord to add
    :param b: Second co-ord to add (usually a vector)
    :return: The addition of a and b
    """
    return a[0] + b[0], a[1] + b[1]


def remap_direction(pipe):
    """
    Convert the value of a pipe into the available directions (vectors) that
    can be taken from that point.
    :param pipe: The pipe that we're starting from.
    :return: A set of directions that are valid to take from this pipe
    (vectors)
    """
    directions = []
    if pipe == '-':  # LR
        directions.append((0, 1))
        directions.append((0, -1))
    elif pipe == '|':  # UD
        directions.append((1, 0))
        directions.append((-1, 0))
    elif pipe == 'J':  # UL
        directions.append((-1, 0))
        directions.append((0, -1))
    elif pipe == 'L':  # UR
        directions.append((-1, 0))
        directions.append((0, 1))
    elif pipe == 'F':  # DR
        directions.append((1, 0))
        directions.append((0, 1))
    elif pipe == '7':  # DL
        directions.append((1, 0))
        directions.append((0, -1))
    return directions


def map_to_pipe(directions):
    """
    Map a set of directions to the pipe value that fits them - essentially
    reverse engineering the above function.
    :param directions: A set of directions that can be taken from a given
                       pipe.
    :return: The pipe value that those directions correspond to
    """
    pipe = '.'
    if (0, 1) in directions:  # R
        if (0, -1) in directions:  # L
            pipe = '-'
        elif (1, 0) in directions:  # D
            pipe = 'F'
        elif (-1, 0) in directions:  # U
            pipe = 'L'
    elif (0, -1) in directions:  # L
        if (1, 0) in directions:  # D
            pipe = '7'
        elif (-1, 0) in directions:  # U
            pipe = 'J'
    elif (1, 0) in directions and (-1, 0) in directions:  # UD
        pipe = '|'
    return pipe


def reverse(direction):
    """
    Reverse the direction of a given vector - turning L to R and U to D.
    :param direction: The direction to reverse
    :return: The reversed direction
    """
    return direction[0] * -1, direction[1] * -1


def get_valid_directions(maze_copy, start_pos):
    """
    Return a list of valid directions that can be taken from a given point
    within the maze.
    :param maze_copy: A copy of the maze from which we're looking up the loop.
    :param start_pos: The start co-ordinate (y,x) that we're looking for valid
                      directions to travel from.
    :return: A list of vectors describing the directions that can be travelled
    from the provided start position.
    """
    valid_directions = []
    if get_map_point(maze_copy, start_pos) != 'S':
        valid_directions = remap_direction(get_map_point(maze_copy, start_pos))
    else:
        possible_directions = [(0, 1), (0, -1), (1, 0), (-1, 0)]
        for vect in possible_directions:
            if reverse(vect) in remap_direction(
                    get_map_point(maze_copy, add_coords(start_pos, vect))):
                valid_directions.append(vect)
    return valid_directions


def get_next_point(maze_copy, start_pos, current_loop):
    """
    Return the next connecting node for a given loop. This consists of finding
    the valid directions that we can take and asking if they've already been
    covered in our loop. If both or neither of them are already covered in the
    loop, it doesn't particularly matter which is returned.
    :param maze_copy: A copy of the maze that we're searching for the loop
                      inside.
    :param start_pos: The starting coordinate (y,x) within the maze.
    :param current_loop: The current list of discovered nodes within the loop.
    :return: The next point to travel within the maze.
    """
    valid_steps = get_valid_directions(maze_copy, start_pos)
    next_step = valid_steps[0]
    if add_coords(start_pos, next_step) in current_loop:
        next_step = valid_steps[1]
    debug_line(f"Going in vector {next_step} from {start_pos}")
    return add_coords(start_pos, next_step)


def test_escaped(maze_copy, test_pos):
    return (test_pos[0] <= 0
            or test_pos[1] <= 0
            or test_pos[0] >= len(maze_copy) - 1
            or test_pos[1] >= len(maze_copy[0]) - 1)


def test_escaped_on_circuit(maze_copy, test_pos, approach):
    pipe = get_map_point(maze_copy, test_pos)
    if pipe == 'S':
        pipe = map_to_pipe(get_valid_directions(maze_copy, test_pos))

    allowed = []
    if pipe == '|':
        allowed.append((0, -1) if 'R' in list(approach["qual"]) else (0, 1))
    if pipe == '-':
        allowed.append((1, 0) if 'U' in list(approach["qual"]) else (-1, 0))
    if pipe == '7' and approach["qual"] != 'UR':
        allowed.append((-1, 0))
        allowed.append((0, 1))
    if pipe == 'J' and approach["qual"] != "DR":
        allowed.append((1, 0))
        allowed.append((0, 1))
    if pipe == 'L' and approach["qual"] != "DL":
        allowed.append((1, 0))
        allowed.append((0, -1))
    if pipe == 'F' and approach["qual"] != "UL":
        allowed.append((-1, 0))
        allowed.append((0, -1))

    escapeable = False
    for vect in allowed:
        escape_pos = add_coords(test_pos, vect)
        if escape_pos[0] < 0 or escape_pos[1] < 0\
            or escape_pos[0] >= len(maze_copy) or escape_pos[1] >= len(maze_copy[0]):
            escapeable = True

    return escapeable


def consider_valid_flood_options(maze_copy, seen_nodes, seen_pipes, potential_directions):
    valid_directions = []
    for test_dir in potential_directions:
        test_pos = add_coords(test_dir["pos"], test_dir["dir"])
        if test_pos in seen_nodes:
            continue
        if { "pos": test_pos, "approach": test_dir["qual"] } in seen_pipes:
            continue
        if (test_pos[0] < 0
                or test_pos[1] < 0
                or test_pos[0] > len(maze_copy) - 1
                or test_pos[1] > len(maze_copy[0]) - 1):
            continue
        valid_directions.append(test_dir)
    return valid_directions


def translate_circuit_flood_options(circuit_pipe, pos, approach):
    """
    I admit, this _hurts_... But it works... So, let's split a pipe into four quadrants and
    say that a direction is approaching in one of those four quadrants in some combination of
    up/down and left/right.

    :param circuit_pipe: The type of pipe we're converting into directions
    :param pos: The position of the pipe in the maze
    :param approach: The direction we're taking to reach the pipe
    :return: A list of possible directions that can be taken from this pipe
    with the given approach vector
    """
    directions = []
    if circuit_pipe == "|":
        if "L" in list(approach["qual"]):
            directions.append({"pos": pos, "dir": (0, 1), "qual": "R"})
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DL"})
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UL"})
        else:
            directions.append({"pos": pos, "dir": (0, -1), "qual": "L"})
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DR"})
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UR"})
    if circuit_pipe == "-":
        if "U" in list(approach["qual"]):
            directions.append({"pos": pos, "dir": (1, 0), "qual": "D"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "UR"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "UL"})
        else:
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "U"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "DR"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "DL"})
    if circuit_pipe == "7":
        if approach["qual"] == "UR":
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DR"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "UL"})
        else:
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "U"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "R"})
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DL"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "DL"})
    if circuit_pipe == "J":
        if approach["qual"] == "DR":
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UR"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "DL"})
        else:
            directions.append({"pos": pos, "dir": (1, 0), "qual": "D"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "R"})
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UL"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "UL"})
    if circuit_pipe == "L":
        if approach["qual"] == "DL":
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UL"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "DR"})
        else:
            directions.append({"pos": pos, "dir": (1, 0), "qual": "D"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "L"})
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "UR"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "UR"})
    if circuit_pipe == "F":
        if approach["qual"] == "UL":
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DL"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "UR"})
        else:
            directions.append({"pos": pos, "dir": (-1, 0), "qual": "U"})
            directions.append({"pos": pos, "dir": (0, -1), "qual": "L"})
            directions.append({"pos": pos, "dir": (1, 0), "qual": "DR"})
            directions.append({"pos": pos, "dir": (0, 1), "qual": "DR"})
    return directions


def attempt_escape(maze_copy, start_pos, circuit, trapped_nodes, free_nodes, trapped_pipes, free_pipes):
    escaped = False
    pathway = [start_pos]
    considered_pipes = []
    options = consider_valid_flood_options(maze_copy, pathway, considered_pipes, [
        {"pos": start_pos, "dir": (0, 1), "qual": "R"},
        {"pos": start_pos, "dir": (0, -1), "qual": "L"},
        {"pos": start_pos, "dir": (1, 0), "qual": "D"},
        {"pos": start_pos, "dir": (-1, 0), "qual": "U"}
    ])
    test_pos = start_pos
    test_pipe = None

    debug_line(f"Running new escape attempt on {start_pos} "
               f"- {len(trapped_nodes)} confirmed trapped nodes "
               f"- {len(free_pipes)} confirmed free nodes "
               f"- {len(trapped_pipes)} confirmed trapped pipe approaches "
               f"- {len(free_pipes)} confirmed free pipe approaches")

    while len(options) > 0:
        debug_line(f"{len(options)} potential options remaining - current pathway {pathway}")

        test_option = options.pop()
        test_pos = add_coords(test_option["pos"], test_option["dir"])
        debug_line(f"Considering option {test_option} at {test_pos}...")

        if test_pos in trapped_nodes or test_pos in free_nodes:
            break
        if test_escaped(maze_copy, test_pos):
            if test_pos in circuit:
                if test_escaped_on_circuit(maze_copy, test_pos, test_option):
                    break
            else:
                pathway.append(test_pos)
                break

        test_directions = []
        if test_pos in circuit:
            cir_pipe = get_map_point(maze_copy, test_pos)
            if cir_pipe == 'S':
                cir_pipe = map_to_pipe(get_valid_directions(maze_copy, test_pos))
            debug_line(f"Considering intersecting circuit pipe {cir_pipe} at {test_pos} on vector quality {test_option['qual']}")

            test_pipe = { "pos": test_pos, "approach": test_option["qual"] }
            if test_pipe in trapped_pipes or test_pipe in free_pipes:
                break
            considered_pipes.append(test_pipe)
            test_directions = translate_circuit_flood_options(cir_pipe, test_pos, test_option)
        else:
            pathway.append(test_pos)
            test_directions.append({ "pos": test_pos, "dir": (0, 1), "qual": "R" })
            test_directions.append({ "pos": test_pos, "dir": (0, -1), "qual": "L" })
            test_directions.append({ "pos": test_pos, "dir": (1, 0), "qual": "D" })
            test_directions.append({ "pos": test_pos, "dir": (-1, 0), "qual": "U" })

        for valid_option in consider_valid_flood_options(
                maze_copy, pathway, considered_pipes, test_directions):
            options.append(valid_option)

    if (test_escaped(maze_copy, test_pos)
            or test_pos in free_nodes
            or test_pipe in free_pipes):
        escaped = True
    return escaped, pathway, considered_pipes

# MAIN CODE STARTS HERE

# Set up some base variables
max_steps = 0
starting_pos = None
loop = []
maze_map = []

inside_nodes = []
inside_pipes = []
outside_nodes = []
outside_pipes = []

# For each line in the file, we extract the maze and throw it into a working
# array to track the maze.
if args.input:
    with open(args.input, encoding="utf-8") as fp:
        for line in fp:
            map_line = MAP_REGEX.match(line)
            if map_line is not None:
                maze_map.append(list(map_line.group(0)))

# Attempt to find the starting position within the maze. If we can't find it,
# then we can't start the program
for idx, row in enumerate(maze_map):
    for col in [col for col, val in enumerate(row) if val == 'S']:
        starting_pos = (idx, col)
        debug_line(f"Starting position found in row {idx + 1}, "
                   f"column {col + 1}")

if starting_pos is None:
    print("Unable to find starting position... exiting")
    sys.exit(1)

# Construct the loop from all the information we have available and store it
# in a working array until we start backtracking on ourselves.
working_pos = starting_pos
while working_pos not in loop:
    loop.append(working_pos)
    working_pos = get_next_point(maze_map, working_pos, loop)

debug_line(f"Found {len(loop)} nodes in the loop")

# FLOOD!
for y, row in enumerate(maze_map):
    for x, col in enumerate(row):
        working_pos = (y, x)
        if working_pos not in loop\
                and working_pos not in inside_nodes\
                and working_pos not in outside_nodes:
            # Consider whether (y,x) is inside or outside the loop
            (node_escaped, nodes_checked, pipes_checked) = attempt_escape(
                    maze_map, working_pos, loop,
                    inside_nodes, outside_nodes,
                    inside_pipes, outside_pipes)

            # Efficiency addition (for what it's worth)
            if node_escaped:
                outside_nodes += nodes_checked
                outside_pipes += pipes_checked
            else:
                inside_nodes += nodes_checked
                inside_pipes += pipes_checked

# Print out the final result
print(f"{len(loop) / 2} steps to reach the anti-point of the loop")
print(f"{len(set(inside_nodes))} trapped nodes and {len(set(outside_nodes))} discovered")
