use std::collections::{HashMap, HashSet};

/// Describe a node in the system, along with its connections to other nodes.
#[derive(Debug, PartialEq, Clone)]
pub(crate) struct Node {
    pub(crate) coord: (usize, usize),
    pub(crate) connections: Vec<Connection>,
    pub(crate) history: Vec<History>,
}

/// Describe the connection one node may have to another node, along with the weight and
/// the direction of travel (vector NSEW)
#[derive(Debug, PartialEq, Clone, Copy)]
pub(crate) struct Connection {
    pub(crate) to_node: (usize, usize),
    pub(crate) weight: i8,
    pub(crate) direction: (i8, i8),
}

/// Describe the vector of travel that we have previously gone down in our history. This
/// tracks the furthest we have travelled in a straight line.
#[derive(Debug, PartialEq, Clone, Copy, Eq, Hash)]
pub(crate) struct Vector {
    pub(crate) direction: (i8, i8),
    pub(crate) distance: u8,
}

/// Describe the historic journey that has gone into reaching a given node, alongside its
/// current distance from the start and the last vector of straight travel.
#[derive(Debug, PartialEq, Clone)]
pub(crate) struct History {
    pub(crate) visited_nodes: Vec<(usize, usize)>,
    pub(crate) distance: i32,
    pub(crate) vector: Vector,
}

/// Explode the input map of strings into a hashmap of co-ordinates to nodes (which are
/// connected to other nodes via coordinate mapping).
pub(crate) fn explode_input_map(input_map: &Vec<String>) -> HashMap<(usize, usize), Node> {
    let mut tile_map: HashMap<(usize, usize), Node> = HashMap::new();
    // First, we generate a hashmap of coordinates against nodes
    for y_val in 0..input_map.len() {
        for x_val in 0..input_map[y_val].len() {
            let new_node = Node { coord: (x_val, y_val), connections: vec![], history: vec![] };
            tile_map.insert((x_val, y_val), new_node);
        }
    }

    // Then we go through the input map and figure out which node maps to which with which
    // weight. We only ever create connections to the right or below us as all connections
    // are bi-directional, meaning left and above maps are already created.
    for y_val in 0..input_map.len() {
        for x_val in 0..input_map[y_val].len() {
            let this_val = input_map[y_val].chars().nth(x_val).unwrap().to_string().parse::<i8>().unwrap();
            let mut right_val: i8 = -1;
            let mut down_val: i8 = -1;
            if x_val < input_map[y_val].len() - 1 {
                right_val = input_map[y_val].chars().nth(x_val + 1).unwrap().to_string().parse::<i8>().unwrap()
            }
            if y_val < input_map.len() - 1 {
                down_val = input_map[y_val + 1].chars().nth(x_val).unwrap().to_string().parse::<i8>().unwrap()
            }
            if right_val > -1 {
                let conn_right = Connection { to_node: (x_val + 1, y_val), weight: right_val, direction: (1, 0) };
                let conn_left = Connection { to_node: (x_val, y_val), weight: this_val, direction: (-1, 0) };
                tile_map.get_mut(&(x_val, y_val)).unwrap().connections.push(conn_right);
                tile_map.get_mut(&(x_val + 1, y_val)).unwrap().connections.push(conn_left);
            }
            if down_val > -1 {
                let conn_down = Connection { to_node: (x_val, y_val + 1), weight: down_val, direction: (0, 1) };
                let conn_up = Connection { to_node: (x_val, y_val), weight: this_val, direction: (0, -1) };
                tile_map.get_mut(&(x_val, y_val)).unwrap().connections.push(conn_down);
                tile_map.get_mut(&(x_val, y_val + 1)).unwrap().connections.push(conn_up);
            }
        }
    }

    return tile_map;
}

/// Given the ongoing history of a graph traversal and the next potential step in the traversal,
/// return whether that step is possible, alongside the updated history if that step were to be
/// performed.
///
/// This function has two modes of operation. If we are using ultra, then the minimum and maximum
/// distance travelled values change.
fn test_ongoing_direction(history: &History, next_step: Connection, ultra: bool) -> (bool, History) {
    let mut updated_history = History{ visited_nodes: history.visited_nodes.to_vec(), distance: history.distance, vector: history.vector };
    updated_history.visited_nodes.push(next_step.to_node);
    updated_history.distance = updated_history.distance + next_step.weight as i32;

    let mut outcome = !history.visited_nodes.contains(&next_step.to_node);
    let dist_min_val = if ultra {4} else {0};
    let dist_max_val = if ultra {11} else {4};

    if history.vector.distance == 0 {
        updated_history.vector.direction = next_step.direction;
        updated_history.vector.distance = 1;
    } else {
        if history.vector.direction == next_step.direction {
            updated_history.vector.distance = updated_history.vector.distance + 1;
            outcome = outcome && updated_history.vector.distance < dist_max_val;
        } else {
            if history.vector.distance < dist_min_val {
                outcome = false;
            } else {
                updated_history.vector.direction = next_step.direction;
                updated_history.vector.distance = 1;
            }
        }
    }

    return (outcome, updated_history);
}

/// Given a list of nodes to check and the history value that we want to inspect, calculate which
/// of them is the best next choice to visit, based on the distance values and the lowest vector
/// of straight line travel as a decider.
fn seek_lowest_unvisited(nodes_to_check: &HashSet<((usize, usize), usize)>, node_map: &HashMap<(usize, usize), Node>) -> ((usize, usize), usize) {
    let mut low_dist: i32 = i32::MAX;
    let mut low_dir: u8 = u8::MAX;
    let mut low_node: ((usize, usize), usize) = ((0, 0), 0);

    for tile in nodes_to_check {
        let local_history = node_map.get(&tile.0).unwrap().history.get(tile.1).unwrap();
        let test_dist = local_history.distance;
        if test_dist < low_dist {
            low_dist = test_dist;
            low_dir = local_history.vector.distance;
            low_node = *tile;
        } else if test_dist == low_dist {
            let test_dir = local_history.vector.distance;
            if test_dir < low_dir {
                low_dir = test_dir;
                low_node = *tile;
            }
        }
    }

    return low_node;
}

/// Find the lowest distance travelled in a given log of all possible paths to reach a
/// node.
fn find_lowest_history_dist(history_log: &Vec<History>) -> i32 {
    let mut low = i32::MAX;

    for hist in history_log {
        if hist.distance < low {
            low = hist.distance;
        }
    }

    return low;
}

/// Given a set of many options that we could visit, narrow down the list of options we will
/// consider to those with a distinct vector approaching any given coordinate. We remove
/// duplicate approaches to a given node based on which one has the lowest distance.
fn reduce_options(nodes_to_check: &HashSet<((usize, usize), usize)>, node_map: &HashMap<(usize, usize), Node>) -> HashSet<((usize, usize), usize)> {
    let mut seen_mapping: HashMap<((usize, usize), Vector), usize> = HashMap::new();

    for check_node in nodes_to_check {
        let local_history: &History = node_map.get(&check_node.0).unwrap().history.get(check_node.1).unwrap();
        if !seen_mapping.contains_key(&(check_node.0, local_history.vector)) {
            seen_mapping.insert((check_node.0, local_history.vector), check_node.1);
        }
        let winning_history: &History = &node_map.get(&check_node.0).unwrap().history.get(*seen_mapping.get(&(check_node.0, local_history.vector)).unwrap()).unwrap();
        if local_history.distance < winning_history.distance {
            seen_mapping.insert((check_node.0, local_history.vector), check_node.1);
        }
    }

    return seen_mapping.keys().fold(HashSet::new(), |mut acc, k| {
        acc.insert((k.0, *seen_mapping.get(k).unwrap()));
        return acc;
    });
}

/// Perform the modified algorithm on an input set of nodes when provided a starting point and
/// a goal node.
pub(crate) fn perform_algorithm<'a>(node_map: &'a mut HashMap<(usize, usize), Node>, start_node: &(usize, usize), goal_node: &(usize, usize), ultra: bool) -> &'a HashMap<(usize, usize), Node> {
    node_map.get_mut(start_node).unwrap().history.push(History {visited_nodes: vec![*start_node], distance: 0, vector: Vector { direction: (0, 0), distance: 0 }});
    println!("Starting from node {:?} with distance of 0", start_node);

    let min_travel_dist = if ultra {4} else {1};
    let waver_val = if ultra {20} else {5};

    // Start from the start node which has an almost blank history.
    let mut nodes_to_check: HashSet<((usize, usize), usize)> = HashSet::new();
    nodes_to_check.insert((*start_node, 0));

    // While we still have nodes to check...
    while !nodes_to_check.is_empty() {
        // Find the lowest unvisited node and remove it from our options
        let lowest_node = seek_lowest_unvisited(&nodes_to_check, node_map);
        nodes_to_check.remove(&lowest_node);

        // Print some basic information about our path so far
        let dist_from_goal = (goal_node.0 - lowest_node.0.0, goal_node.1 - lowest_node.0.1);
        println!("Inspecting node {:?} on approach {:?} - {:?} options remaining - {:?} coords from goal", lowest_node.0, lowest_node.1, nodes_to_check.len(), dist_from_goal);

        let local_history = node_map.get(&lowest_node.0).unwrap().history.get(lowest_node.1).unwrap().clone();

        // If we've reached the end, we're guaranteed to have already found our best candidate, so
        // quit now.
        if lowest_node.0 == *goal_node {
            if local_history.vector.distance >= min_travel_dist {
                break;
            }
        }

        // For each connection to this node...
        for conn in node_map.get(&lowest_node.0).unwrap().connections.to_vec() {
            let (can_continue, min_history) = test_ongoing_direction(&local_history, conn, ultra);

            // If we can take another step in this direction - and the weight of the node would be
            // improved by connecting to it now...
            if can_continue {
                let challenge_dist = find_lowest_history_dist(&node_map.get(&conn.to_node).unwrap().history.to_vec());

                if min_history.distance - waver_val < challenge_dist {
                    let hop_node = node_map.get_mut(&conn.to_node).unwrap();
                    hop_node.history.push(min_history);
                    nodes_to_check.insert((conn.to_node, hop_node.history.len() - 1));
                }
            }
        }

        // Reduce our options to only worthwhile potentials
        nodes_to_check = reduce_options(&nodes_to_check, node_map);
    }

    return node_map;
}
