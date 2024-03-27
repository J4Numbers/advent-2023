use std::cmp::{max, min};

/// Describe a digger instruction that can be found in our input file.
#[derive(Debug, PartialEq, Clone)]
pub(crate) struct DigInstruction {
    pub(crate) direction: char,
    pub(crate) count: i32,
    pub(crate) colour: String,
}

pub(crate) fn generate_polygon(instructions: &Vec<DigInstruction>) -> Vec<(i32, i32)> {
    let start_node: (i32, i32) = (0, 0);
    let mut path = vec![start_node];

    let mut working_node = start_node;
    for instr in instructions {
        if instr.direction == 'U' {
            working_node = (working_node.0, working_node.1 - instr.count);
        } else if instr.direction == 'D' {
            working_node = (working_node.0, working_node.1 + instr.count);
        } else if instr.direction == 'L' {
            working_node = (working_node.0 - instr.count, working_node.1);
        } else if instr.direction == 'R' {
            working_node = (working_node.0 + instr.count, working_node.1);
        }
        path.push(working_node);
    }

    return path;
}

pub(crate) fn offset_graph(points: &Vec<(i32, i32)>, offset: (i32, i32)) -> Vec<(i32, i32)> {
    let mut offset_points: Vec<(i32, i32)> = vec![];
    for point in points {
        offset_points.push((point.0 - offset.0, point.1 - offset.1));
    }
    return offset_points;
}

pub(crate) fn between_points(points: &Vec<(i32, i32)>, point_under_test: (i32, i32)) -> bool {
    let mut between = false;
    let mut prev_point = (0, 0);
    for point in points {
        between = between || (point_under_test.0 <= max(prev_point.0, point.0) && point_under_test.0 >= min(prev_point.0, point.0) && point_under_test.1 == prev_point.1 && point_under_test.1 == point.1);
        between = between || (point_under_test.1 <= max(prev_point.1, point.1) && point_under_test.1 >= min(prev_point.1, point.1) && point_under_test.0 == prev_point.0 && point_under_test.0 == point.0);
        if between {
            break;
        }
        prev_point = *point;
    }
    return between;
}

pub(crate) fn get_intersecting_lines(points: &Vec<(i32, i32)>, point_under_test: (i32, i32)) -> Vec<((i32, i32), (i32, i32))> {
    let mut all_intersecting_lines: Vec<((i32, i32), (i32, i32))> = vec![];
    let mut prev_point = (0, 0);

    for point in points {
        if (point_under_test.0 <= max(prev_point.0, point.0) && point_under_test.0 >= min(prev_point.0, point.0) && point_under_test.1 == prev_point.1 && point_under_test.1 == point.1) || (point_under_test.1 <= max(prev_point.1, point.1) && point_under_test.1 >= min(prev_point.1, point.1) && point_under_test.0 == prev_point.0 && point_under_test.0 == point.0) {
            all_intersecting_lines.push((*point, prev_point));
        }
        prev_point = *point;
    }
    return all_intersecting_lines;
}

pub(crate) fn within_shape(points: &Vec<(i32, i32)>, point_under_test: (i32, i32)) -> bool {
    let mut l_count = 0;
    let mut red_x = 0;
    while red_x < point_under_test.0 {
        if between_points(points, (red_x, point_under_test.1)) {
            let valid_lines = get_intersecting_lines(points, (red_x, point_under_test.1));
            for (pt_a, pt_b) in valid_lines.to_vec() {
                if max(pt_a.0, pt_b.0) > red_x {
                    red_x = max(pt_a.0, pt_b.0)
                }
            }
            if valid_lines.len() == 1 {
                l_count += 1;
            }
        }
        red_x = red_x + 1;
    }

    println!("{:?} found {:?} intersecting lines to the left", point_under_test, l_count);
    return l_count % 2 == 1;
}

pub(crate) fn generate_graph(points: &Vec<(i32, i32)>) -> Vec<Vec<char>> {
    let (top_left, bottom_right) = generate_bounds(points);

    let offset_points = offset_graph(points, top_left);

    println!("{:?} with offset of {:?} generates {:?}", points, top_left, offset_points);

    let y_diff = bottom_right.1 - top_left.1;
    let x_diff = bottom_right.0 - top_left.0;

    let mut char_graph: Vec<Vec<char>> = vec![];

    for y in 0..y_diff+1 {
        let mut y_row: Vec<char> = vec![];
        for x in 0..x_diff+1 {
            if between_points(&offset_points, (x, y)) || within_shape(&offset_points,(x, y)) {
                y_row.push('#');
            } else {
                y_row.push('.');
            }
        }
        char_graph.push(y_row);
    }

    return char_graph;
}

pub(crate) fn count_filled_squares(dig_graph: &Vec<Vec<char>>) -> i32 {
    let mut shaded = 0;
    for line in dig_graph {
        let mut started = false;
        let mut working_count = 0;
        for char in line {
            if started {
                working_count = working_count + 1;
            }
            if *char == '#' {
                if started {
                    shaded = shaded + working_count;
                    working_count = 0;
                } else {
                    shaded = shaded + 1;
                    started = true;
                }
            }
        }
    }
    return shaded;
}

pub(crate) fn generate_bounds(points: &Vec<(i32, i32)>) -> ((i32, i32), (i32, i32)) {
    let mut top_left = (0, 0);
    let mut bottom_right = (0, 0);

    for point in points {
        if point.0 > bottom_right.0 {
            bottom_right.0 = point.0;
        } else if point.0 < top_left.0 {
            top_left.0 = point.0;
        }
        if point.1 > bottom_right.1 {
            bottom_right.1 = point.1;
        } else if point.1 < top_left.1 {
            top_left.1 = point.1;
        }
    }

    return (top_left, bottom_right);
}

pub(crate) fn calculate_trench_len(instructions: &Vec<DigInstruction>) -> i32 {
    let mut trench_len = 0;
    for instr in instructions {
        trench_len = trench_len + instr.count;
    }
    return trench_len;
}

pub(crate) fn calculate_area(points: &Vec<(i32, i32)>) -> i32 {
    let mut area = 0;
    let mut j = points.len() - 1;

    for i in 0..points.len() {
        println!("Comparing {:?} against {:?}", points.get(j).unwrap(), points.get(i).unwrap());
        area += (points.get(j).unwrap().0 + points.get(i).unwrap().0) * (points.get(j).unwrap().1 - points.get(i).unwrap().1);
        j = i;
    }

    return area.abs() / 2;
}
