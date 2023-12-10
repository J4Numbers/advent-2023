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


# MAIN CODE STARTS HERE

# Set up some base variables
max_steps = 0
starting_pos = None
loop = []
maze_map = []

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

# Print out the final result
print(len(loop) / 2)
