mod fileio;
mod digger;

use clap::Parser;

/// Describe the args that are going to be fed into this program.
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct ProgArgs {
    #[arg(short, long)]
    input: String,
}

/// Kick off the main program to perform the operations described in the readme.
fn main() {
    let args = ProgArgs::parse();
    let dig_instructions = fileio::read_file(args.input);

    println!("{:?} instructions found", dig_instructions.len());

    let points = digger::generate_polygon(&dig_instructions);
    let trench_len = digger::calculate_trench_len(&dig_instructions);
    let graph = digger::generate_graph(&points);

    for ln in graph.to_vec() {
        println!("{:?}", ln);
    }

    println!("Total area - {:?}", digger::count_filled_squares(&graph));

    // And print out the number of lit tiles that was the maximum seen.
    // println!("Final area was {:?} + {:?} = {:?}", area, trench_len, area + trench_len)
}
