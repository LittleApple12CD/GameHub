use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;

const ROWS: i32 = 16;
const COLS: i32 = 16;
const MINES: i32 = 40;
const CELL: f32 = 30.0;
const W: f32 = COLS as f32 * CELL;
const H: f32 = ROWS as f32 * CELL;

pub struct MinesweeperGame {
    board: Vec<Vec<i32>>,
    revealed: Vec<Vec<bool>>,
    flagged: Vec<Vec<bool>>,
    game_over: bool,
    won: bool,
    first_click: bool,
    offset: (f32, f32),
}

impl MinesweeperGame {
    pub fn new() -> Self {
        let mut g = Self {
            board: vec![vec![0; COLS as usize]; ROWS as usize],
            revealed: vec![vec![false; COLS as usize]; ROWS as usize],
            flagged: vec![vec![false; COLS as usize]; ROWS as usize],
            game_over: false,
            won: false,
            first_click: true,
            offset: centered_offset(W, H),
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.board = vec![vec![0; COLS as usize]; ROWS as usize];
        self.revealed = vec![vec![false; COLS as usize]; ROWS as usize];
        self.flagged = vec![vec![false; COLS as usize]; ROWS as usize];
        self.game_over = false;
        self.won = false;
        self.first_click = true;
    }

    fn place_mines(&mut self, safe_x: i32, safe_y: i32) {
        let mut rng = ::rand::thread_rng();
        let mut placed = 0;
        while placed < MINES {
            let x = rng.gen_range(0..COLS);
            let y = rng.gen_range(0..ROWS);
            if self.board[y as usize][x as usize] == -1 {
                continue;
            }
            if (x - safe_x).abs() <= 1 && (y - safe_y).abs() <= 1 {
                continue;
            }
            self.board[y as usize][x as usize] = -1;
            placed += 1;
        }
        // 计周围雷数
        for y in 0..ROWS {
            for x in 0..COLS {
                if self.board[y as usize][x as usize] == -1 {
                    continue;
                }
                let mut count = 0;
                for dy in -1..=1 {
                    for dx in -1..=1 {
                        let nx = x + dx;
                        let ny = y + dy;
                        if nx >= 0 && nx < COLS && ny >= 0 && ny < ROWS {
                            if self.board[ny as usize][nx as usize] == -1 {
                                count += 1;
                            }
                        }
                    }
                }
                self.board[y as usize][x as usize] = count;
            }
        }
    }

    fn reveal(&mut self, x: i32, y: i32) {
        if x < 0 || x >= COLS || y < 0 || y >= ROWS {
            return;
        }
        if self.revealed[y as usize][x as usize] || self.flagged[y as usize][x as usize] {
            return;
        }
        self.revealed[y as usize][x as usize] = true;

        if self.board[y as usize][x as usize] == -1 {
            self.game_over = true;
            return;
        }

        if self.board[y as usize][x as usize] == 0 {
            for dy in -1..=1 {
                for dx in -1..=1 {
                    self.reveal(x + dx, y + dy);
                }
            }
        }

        let revealed_count: i32 = self
            .revealed
            .iter()
            .map(|row| row.iter().filter(|&&r| r).count() as i32)
            .sum();
        if revealed_count == ROWS * COLS - MINES {
            self.won = true;
        }
    }
}

impl Game for MinesweeperGame {
    fn update(&mut self) -> bool {
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }

        if !self.game_over && !self.won {
            let (mx, my) = mouse_position();
            let (ox, oy) = self.offset;
            let x = ((mx - ox) / CELL).floor() as i32;
            let y = ((my - oy) / CELL).floor() as i32;

            if x >= 0 && x < COLS && y >= 0 && y < ROWS {
                // 左键：翻开
                if is_mouse_button_pressed(MouseButton::Left) {
                    if self.first_click {
                        self.place_mines(x, y);
                        self.first_click = false;
                    }
                    if !self.flagged[y as usize][x as usize] {
                        self.reveal(x, y);
                    }
                }
                // 右键：插旗
                if is_mouse_button_pressed(MouseButton::Right) {
                    if !self.revealed[y as usize][x as usize] {
                        self.flagged[y as usize][x as usize] =
                            !self.flagged[y as usize][x as usize];
                    }
                }
            }
        }
        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        draw_rectangle(ox, oy, W, H, Color::new(0.118, 0.118, 0.165, 1.0));

        for y in 0..ROWS {
            for x in 0..COLS {
                let rx = ox + x as f32 * CELL + 1.0;
                let ry = oy + y as f32 * CELL + 1.0;
                let rw = CELL - 2.0;
                let rh = CELL - 2.0;

                if self.revealed[y as usize][x as usize] {
                    if self.board[y as usize][x as usize] == -1 {
                        draw_rounded_rect(rx, ry, rw, rh, 4.0, Color::new(0.784, 0.196, 0.196, 1.0));
                        draw_circle(rx + CELL / 2.0, ry + CELL / 2.0, 9.0, BLACK);
                    } else {
                        draw_rounded_rect(rx, ry, rw, rh, 4.0, Color::new(0.745, 0.745, 0.784, 1.0));
                        let n = self.board[y as usize][x as usize];
                        if n > 0 {
                            let colors = [
                                BLACK,
                                Color::new(0.0, 0.0, 1.0, 1.0),
                                Color::new(0.0, 0.588, 0.0, 1.0),
                                Color::new(1.0, 0.0, 0.0, 1.0),
                                Color::new(0.0, 0.0, 0.706, 1.0),
                                Color::new(0.588, 0.0, 0.0, 1.0),
                                Color::new(0.0, 0.588, 0.588, 1.0),
                                BLACK,
                            ];
                            let t = n.to_string();
                            let dims = measure_text(&t, None, 22, 1.0);
                            draw_text(
                                &t,
                                rx + (CELL - dims.width) / 2.0,
                                ry + CELL / 2.0 + dims.height / 2.0 - 2.0,
                                22.0,
                                colors[n.min(7) as usize],
                            );
                        }
                    }
                } else {
                    let flagged = self.flagged[y as usize][x as usize];
                    let color = if flagged {
                        Color::new(0.824, 0.824, 0.235, 1.0)
                    } else {
                        Color::new(0.314, 0.314, 0.412, 1.0)
                    };
                    draw_rounded_rect(rx, ry, rw, rh, 4.0, color);
                    if flagged {
                        let dims = measure_text("F", None, 22, 1.0);
                        draw_text(
                            "F",
                            rx + (CELL - dims.width) / 2.0,
                            ry + CELL / 2.0 + dims.height / 2.0 - 2.0,
                            22.0,
                            Color::new(1.0, 0.196, 0.196, 1.0),
                        );
                    }
                }
            }
        }

        if self.game_over {
            draw_overlay(ox, oy, W, H, 0.6, BLACK);
            draw_text_centered("GAME OVER", ox + W / 2.0, oy + H / 2.0 - 20.0, 40, WHITE);
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 30.0, 22, WHITE);
        } else if self.won {
            draw_overlay(ox, oy, W, H, 0.6, BLACK);
            draw_text_centered("YOU WIN!", ox + W / 2.0, oy + H / 2.0 - 20.0, 40, Color::new(0.0, 1.0, 0.0, 1.0));
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 30.0, 22, WHITE);
        }
    }
}