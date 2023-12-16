mod fileio;
mod light_line;

use clap::Parser;

/// Describe the args that are going to be fed into this program.
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct ProgArgs {
    #[arg(short, long)]
    input: String,

    #[arg(long, default_value_t = false)]
    explore: bool,

    #[arg(long, default_value_t = false)]
    debug: bool,
}

/// Kick off the main program to perform the operations described in the readme.
fn main() {
    let args = ProgArgs::parse();
    let map_contents = fileio::read_file(args.input);

    let mut max_seen: u32 = 0;
    let mut starting_points: Vec<light_line::TileApproach> = vec![];

    // Run the program in one of two modes. If we're running on explore, then we generate a list
    // of all edge nodes with a vector pointing inwards. Otherwise just add a single starting
    // point of [0,0], going right.
    if args.explore {
        let max_height = map_contents.len();
        let max_width = if max_height > 0 {map_contents[0].len()} else {0};
        for y_start in 0..max_height {
            starting_points.push(light_line::TileApproach{ tile: (0, y_start as i32), vector: (1, 0)});
            starting_points.push(light_line::TileApproach{ tile: ((max_width - 1) as i32, y_start as i32), vector: (-1, 0)});
        }
        for x_start in 0..max_width {
            starting_points.push(light_line::TileApproach{ tile: (x_start as i32, 0), vector: (0, 1)});
            starting_points.push(light_line::TileApproach{ tile: (x_start as i32, (max_height - 1) as i32), vector: (0, -1)});
        }
    } else {
        starting_points.push(light_line::TileApproach { tile: (0, 0), vector: (1, 0)});
    }

    // For each starting point, we calculate the path of light and pull back the number of tiles
    // that were lit as part of that path.
    for starting_point in starting_points {
        let seen_tiles = light_line::calculate_light_line(starting_point, &map_contents);
        println!("Starting point {:?} found {:?} lit tiles", starting_point, seen_tiles);
        if seen_tiles > max_seen {
            max_seen = seen_tiles;
        }
    }

    // And print out the number of lit tiles that was the maximum seen.
    println!("Found {:?} lit tiles as maximum", max_seen);
}
