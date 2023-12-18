mod fileio;
mod dijkstra;

use clap::Parser;
use crate::dijkstra::{History, Vector};

/// Describe the args that are going to be fed into this program.
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct ProgArgs {
    #[arg(short, long)]
    input: String,
}

fn lookup_minimum_history_val(histories: &Vec<dijkstra::History>) -> dijkstra::History {
    let mut min_hist = History{ distance: 999999999, visited_nodes: vec![], vector: Vector{ direction: (0, 0), distance: 0 }};
    for history in histories {
        if history.distance < min_hist.distance {
            min_hist = history.clone();
        }
    }
    return min_hist;
}

/// Kick off the main program to perform the operations described in the readme.
fn main() {
    let args = ProgArgs::parse();
    let map_contents = fileio::read_file(args.input);

    let start_node: (usize, usize) = (0, 0);
    let goal_node = (map_contents[0].len() - 1, map_contents.len() - 1);

    let mut tile_connections = dijkstra::explode_input_map(&map_contents);
    dijkstra::perform_algorithm(&mut tile_connections, &start_node, &goal_node);

    let min_hist = lookup_minimum_history_val(&tile_connections.get(&goal_node).unwrap().history);

    println!("Minimum history value found was {:?}", min_hist);

    println!();
    for y_val in 0..map_contents.len() {
        let mut working_line: String = "".to_owned();
        for x_val in 0..map_contents[0].len() {
            if min_hist.visited_nodes.contains(&(x_val, y_val)) {
                working_line.push_str("*");
            } else {
                working_line.push_str(&*map_contents[y_val].chars().nth(x_val).unwrap().to_string())
            }
        }
        println!("{:?}", working_line);
    }
}
