/// Describe a simple data object for how we approach a tile in the reflection map,
/// made up of the tile we're inspecting at the moment and the vector that we're
/// approaching the tile on.
#[derive(Debug, PartialEq, Clone, Copy)]
pub(crate) struct TileApproach {
    pub(crate) tile: (i32, i32),
    pub(crate) vector: (i32, i32)
}

/// Test whether a given tile approach is something valid. We treat validity as a two-case
/// situation. If the new tile approach goes out of bounds of the map, then it is not valid.
/// In addition, if we have already seen the exact approach before, there's no point
/// attempting the approach again, so we strike it from the list.
///
/// # Arguments
///
/// * `test_opt` - The test approach we're going to check for validity
/// * `map` - The current map of tiles that we're working on - predominantly used for bounds
///   checking
/// * `witnessed` - A list of tile approaches that we have already seen and that we can
///   discount if we see again.
fn test_valid_option(test_opt: &TileApproach, map: &Vec<String>, witnessed: &Vec<TileApproach>) -> bool {
    let max_height = map.len() as i32;
    let max_width = if max_height > 0 {map[0].len() as i32} else {0i32};
    if test_opt.tile.0 < 0 || test_opt.tile.1 < 0 || test_opt.tile.0 >= max_width || test_opt.tile.1 >= max_height {
        return false;
    }
    return !witnessed.contains(&test_opt);
}

/// Generate a list of options that can be taken from a given approach. This approach will inspect
/// the tile described by the input tile, then, depending on the tile and the input vector, will
/// generate a list of possible options that can be travelled on this run of the light run. This
/// will not return options that go out of bounds or that have already been witnessed.
///
/// # Arguments
///
/// * `input_opt` - The tile approach that we're going to generate options from, containing the
///   tile we're looking at right now, and the direction of travel we're moving in.
/// * `map` - The map of mirrors that will contain the tiles we're inspecting.
/// * `witnessed` - The list of tile approaches that we've already seen and can discount safely.
fn generate_options(input_opt: TileApproach, map: &Vec<String>, witnessed: &Vec<TileApproach>) -> Vec<TileApproach> {
    let mut ret_options: Vec<TileApproach> = vec![];
    let test_chr = map[input_opt.tile.1 as usize].chars().nth(input_opt.tile.0 as usize).unwrap();
    if test_chr == '.' {
        ret_options.push(TileApproach{ tile: (input_opt.tile.0 + input_opt.vector.0, input_opt.tile.1 + input_opt.vector.1), vector: input_opt.vector })
    } else if test_chr == '/' {
        let mod_vec: (i32, i32);
        if input_opt.vector == (0, 1) {
            mod_vec = (-1, 0)
        } else if input_opt.vector == (0, -1) {
            mod_vec = (1, 0)
        } else if input_opt.vector == (1, 0) {
            mod_vec = (0, -1)
        } else {
            mod_vec = (0, 1)
        }
        ret_options.push(TileApproach{ tile: (input_opt.tile.0 + mod_vec.0, input_opt.tile.1 + mod_vec.1), vector: mod_vec })
    } else if test_chr == '\\' {
        let mod_vec: (i32, i32);
        if input_opt.vector == (0, 1) {
            mod_vec = (1, 0)
        } else if input_opt.vector == (0, -1) {
            mod_vec = (-1, 0)
        } else if input_opt.vector == (1, 0) {
            mod_vec = (0, 1)
        } else {
            mod_vec = (0, -1)
        }
        ret_options.push(TileApproach{ tile: (input_opt.tile.0 + mod_vec.0, input_opt.tile.1 + mod_vec.1), vector: mod_vec })
    } else if test_chr == '-' {
        if input_opt.vector == (0, 1) || input_opt.vector == (0, -1) {
            ret_options.push(TileApproach{ tile: (input_opt.tile.0 + 1, input_opt.tile.1), vector: (1, 0) });
            ret_options.push(TileApproach{ tile: (input_opt.tile.0 - 1, input_opt.tile.1), vector: (-1, 0) });
        } else {
            ret_options.push(TileApproach{ tile: (input_opt.tile.0 + input_opt.vector.0, input_opt.tile.1 + input_opt.vector.1), vector: input_opt.vector });
        }
    } else if test_chr == '|' {
        if input_opt.vector == (1, 0) || input_opt.vector == (-1, 0) {
            ret_options.push(TileApproach{ tile: (input_opt.tile.0, input_opt.tile.1 + 1), vector: (0, 1) });
            ret_options.push(TileApproach{ tile: (input_opt.tile.0, input_opt.tile.1 - 1), vector: (0, -1) });
        } else {
            ret_options.push(TileApproach{ tile: (input_opt.tile.0 + input_opt.vector.0, input_opt.tile.1 + input_opt.vector.1), vector: input_opt.vector });
        }
    }
    return ret_options.into_iter()
        .filter(|opt| test_valid_option(opt, map, witnessed))
        .collect();
}

/// Count the number of seen tiles when given a list of approaches that we've seen. The two
/// are not necessarily the same thing as we may have witnessed the same tile more than once
/// from different directions.
///
/// # Arguments
///
/// * `witnessed_tiles` - A list of all witnessed tile approaches that we're going to simplify
///   down into a tile set.
fn count_seen_tiles(witnessed_tiles: Vec<TileApproach>) -> u32 {
    let mut found_tiles: Vec<(i32, i32)> = vec![];

    for tile_seen in witnessed_tiles {
        if !found_tiles.contains(&tile_seen.tile) {
            found_tiles.push(tile_seen.tile)
        }
    }

    return found_tiles.len() as u32;
}

/// Drag out the function to calculate the number of lit tiles when provided the starting
/// approach and the map of mirrored surfaces that we're going to travel through. This
/// calculation is two-step:
///
/// * Calculate all paths of light from a starting point, dropping options when we leave
///   the map, or see a given tile approach more than once.
/// * Calculate the set of tiles that we've seen that have been lit at least once and
///   count them.
///
/// # Arguments
///
/// * `start_approach` - The approach we take to start on the map.
/// * `map` - The map of mirrors that we're exploring.
pub(crate) fn calculate_light_line(start_approach: TileApproach, map: &Vec<String>) -> u32 {
    let mut seen: Vec<TileApproach> = vec![];
    let mut options: Vec<TileApproach> = vec![];
    options.push(start_approach);

    while !options.is_empty() {
        let working_opt = options.pop().unwrap();
        let new_opts = generate_options(working_opt, map, &seen);
        // println!("Inspecting tile at {:?} discovered options {:?}", working_opt, new_opts);
        for tile_opt in new_opts {
            options.push(tile_opt)
        }
        seen.push(working_opt)
    }

    return count_seen_tiles(seen);
}
