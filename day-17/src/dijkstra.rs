use std::collections::{HashMap, HashSet};

/// Describe a node in the system, along with its connections to other nodes.
#[derive(Debug, PartialEq, Clone)]
pub(crate) struct Node {
    pub(crate) coord: (usize, usize),
    pub(crate) connections: Vec<Connection>,
    pub(crate) history: Vec<History>,
}

#[derive(Debug, PartialEq, Clone, Copy)]
pub(crate) struct Connection {
    pub(crate) to_node: (usize, usize),
    pub(crate) weight: i8,
    pub(crate) direction: (i8, i8),
}

#[derive(Debug, PartialEq, Clone, Copy, Eq, Hash)]
pub(crate) struct Vector {
    pub(crate) direction: (i8, i8),
    pub(crate) distance: u8,
}

#[derive(Debug, PartialEq, Clone)]
pub(crate) struct History {
    pub(crate) visited_nodes: Vec<(usize, usize)>,
    pub(crate) distance: i32,
    pub(crate) vector: Vector,
}

pub(crate) fn explode_input_map(input_map: &Vec<String>) -> HashMap<(usize, usize), Node> {
    let mut tile_map: HashMap<(usize, usize), Node> = HashMap::new();
    for y_val in 0..input_map.len() {
        for x_val in 0..input_map[y_val].len() {
            let new_node = Node { coord: (x_val, y_val), connections: vec![], history: vec![] };
            tile_map.insert((x_val, y_val), new_node);
        }
    }

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

fn test_ongoing_direction(history: &History, next_step: Connection) -> (bool, History) {
    let mut updated_history = History{ visited_nodes: history.visited_nodes.to_vec(), distance: history.distance, vector: history.vector };
    updated_history.visited_nodes.push(next_step.to_node);
    updated_history.distance = updated_history.distance + next_step.weight as i32;

    if history.vector.direction.0 == next_step.direction.0 && history.vector.direction.1 == next_step.direction.1 {
        updated_history.vector.distance = updated_history.vector.distance + 1;
    } else {
        updated_history.vector.direction = next_step.direction;
        updated_history.vector.distance = 1;
    }

    let outcome = !history.visited_nodes.contains(&next_step.to_node) && updated_history.vector.distance < 4;

    return (outcome, updated_history);
}

fn seek_lowest_unvisited(nodes_to_check: &mut HashSet<((usize, usize), usize)>, node_map: &HashMap<(usize, usize), Node>) -> ((usize, usize), usize) {
    let mut low_dist: i32 = i32::MAX;
    let mut low_dir: u8 = u8::MAX;
    let mut low_node: ((usize, usize), usize) = ((0, 0), 0);

    let immutable_copy: &HashSet<((usize, usize), usize)> = nodes_to_check;
    for tile in immutable_copy {
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

    nodes_to_check.remove(&low_node);
    return low_node;
}

fn find_lowest_history_dist(history_log: &Vec<History>) -> i32 {
    let mut low = i32::MAX;

    for hist in history_log {
        if hist.distance < low {
            low = hist.distance;
        }
    }

    return low;
}

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

pub(crate) fn perform_algorithm<'a>(node_map: &'a mut HashMap<(usize, usize), Node>, start_node: &(usize, usize), goal_node: &(usize, usize)) -> &'a HashMap<(usize, usize), Node> {
    node_map.get_mut(start_node).unwrap().history.push(History {visited_nodes: vec![*start_node], distance: 0, vector: Vector { direction: (0, 0), distance: 0 }});
    println!("Starting from node {:?} with distance of 0", start_node);

    let mut nodes_to_check: HashSet<((usize, usize), usize)> = HashSet::new();
    nodes_to_check.insert((*start_node, 0));

    while !nodes_to_check.is_empty() {
        let lowest_node = seek_lowest_unvisited(&mut nodes_to_check, node_map);
        let dist_from_goal = (goal_node.0 - lowest_node.0.0, goal_node.1 - lowest_node.0.1);
        println!("Inspecting node {:?} on approach {:?} - {:?} options remaining - {:?} coords from goal", lowest_node.0, lowest_node.1, nodes_to_check.len(), dist_from_goal);

        if lowest_node.0 == *goal_node {
            break;
        }

        let local_history = node_map.get(&lowest_node.0).unwrap().history.get(lowest_node.1).unwrap().clone();

        // For each connection to this node...
        for conn in node_map.get(&lowest_node.0).unwrap().connections.to_vec() {
            let (can_continue, min_history) = test_ongoing_direction(&local_history, conn);

            // If we can take another step in this direction - and the weight of the node would be
            // improved by connecting to it now...
            if can_continue {
                let challenge_dist = find_lowest_history_dist(&node_map.get(&conn.to_node).unwrap().history.to_vec());

                if min_history.distance -5 < challenge_dist {
                    let hop_node = node_map.get_mut(&conn.to_node).unwrap();
                    hop_node.history.push(min_history);
                    nodes_to_check.insert((conn.to_node, hop_node.history.len() - 1));
                }
            }
        }

        nodes_to_check = reduce_options(&nodes_to_check, node_map);
    }

    return node_map;
}
