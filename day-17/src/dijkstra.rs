use std::collections::{HashMap, HashSet};

/// Describe a node in the system, along with its connections to other nodes.
#[derive(Debug, PartialEq, Clone)]
pub(crate) struct Node {
    pub(crate) coord: (usize, usize),
    pub(crate) min_dist: i32,
    pub(crate) connections: Vec<Connection>,
    pub(crate) history: History,
}

#[derive(Debug, PartialEq, Clone, Copy)]
pub(crate) struct Connection {
    pub(crate) to_node: (usize, usize),
    pub(crate) weight: i8,
    pub(crate) direction: (i8, i8),
}

#[derive(Debug, PartialEq, Clone)]
pub(crate) struct History {
    pub(crate) visited_nodes: Vec<(usize, usize)>,
    pub(crate) current_line: (i8, i8),
}

pub(crate) fn explode_input_map(input_map: &Vec<String>) -> HashMap<(usize, usize), Node> {
    let mut tile_map: HashMap<(usize, usize), Node> = HashMap::new();
    for y_val in 0..input_map.len() {
        for x_val in 0..input_map[y_val].len() {
            let new_node = Node { coord: (x_val, y_val), connections: vec![], min_dist: 999999, history: History{ visited_nodes: vec![], current_line: (0, 0) } };
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

fn test_ongoing_direction(history: &History, next_step: Connection) -> (bool, (i8, i8)) {
    let new_vector: (i8, i8);
    if (history.current_line.0 > 0 && next_step.direction.0 > 0) || (history.current_line.0 < 0 && next_step.direction.0 < 0) || (history.current_line.1 > 0 && next_step.direction.1 > 0) || (history.current_line.1 < 0 && next_step.direction.1 < 0) {
        new_vector = (history.current_line.0 + next_step.direction.0, history.current_line.1 + next_step.direction.1);
    } else {
        new_vector = next_step.direction;
    }
    let outcome = !history.visited_nodes.contains(&next_step.to_node) && new_vector.0.abs() < 4 && new_vector.1.abs() < 4;

    return (outcome, new_vector);
}

fn seek_lowest_unvisited(unvisited_nodes: &HashSet<(usize, usize)>, node_map: &HashMap<(usize, usize), Node>) -> (usize, usize) {
    let mut low_val: i32 = 999999999;
    let mut low_node: (usize, usize) = (0, 0);

    for tile in unvisited_nodes {
        if node_map.get(&tile).unwrap().min_dist < low_val {
            low_val = node_map.get(&tile).unwrap().min_dist;
            low_node = (tile.0, tile.1);
        }
    }

    return low_node;
}

pub(crate) fn perform_algorithm<'a>(node_map: &'a mut HashMap<(usize, usize), Node>, start_node: &(usize, usize)) -> &'a HashMap<(usize, usize), Node> {
    node_map.get_mut(start_node).unwrap().min_dist = 0;
    println!("Starting from node {:?} with distance of 0", start_node);

    let mut visited_nodes: HashSet<(usize, usize)> = node_map.keys().cloned().collect();

    while !visited_nodes.is_empty() {
        let lowest_node = seek_lowest_unvisited(&visited_nodes, node_map);
        println!("Inspecting node {:?} with current distance of {:?} :: {:?}", lowest_node, node_map.get(&lowest_node).unwrap().min_dist, node_map.get(&lowest_node).unwrap().history);

        for conn in node_map.get(&lowest_node).unwrap().connections.to_vec() {
            println!("Comparing connection from {:?} to {:?} with weight of {:?}, resulting in a weight of {:?}", lowest_node, conn.to_node, conn.weight, node_map.get(&lowest_node).unwrap().min_dist + conn.weight as i32);
            let (valid_history, updated_vector) = test_ongoing_direction(&node_map.get(&lowest_node).unwrap().history, conn);
            if valid_history && (node_map.get(&lowest_node).unwrap().min_dist + conn.weight as i32) < node_map.get(&conn.to_node).unwrap().min_dist {
                let mut newly_visited_node = node_map.get(&lowest_node).unwrap().history.visited_nodes.to_vec();
                let new_weight = node_map.get(&lowest_node).unwrap().min_dist + conn.weight as i32;
                let hop_node = node_map.get_mut(&conn.to_node).unwrap();
                newly_visited_node.push(conn.to_node);
                hop_node.min_dist = new_weight;
                hop_node.history = History{ visited_nodes: newly_visited_node, current_line: updated_vector };
            }
        }
        visited_nodes.remove(&lowest_node);
    }

    return node_map;
}
