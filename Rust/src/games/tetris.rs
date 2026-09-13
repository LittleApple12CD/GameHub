use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;

const COLS: i32 = 10;
const ROWS: i32 = 20;
const CELL: f32 = 30.0;
const W: f32 = COLS as f32 * CELL; // 300
const H: f32 = ROWS as f32 * CELL; // 600

/// 方块形状：每种用 0/1 的矩形表示
fn shape_of(kind: char) -> Vec<Vec<u8>> {
    match kind {
        'I' => vec![vec![0, 0, 0, 0], vec![1, 1, 1, 1], vec![0, 0, 0, 0], vec![0, 0, 0, 0]],
        'O' => vec![vec![1, 1], vec![1, 1]],
        'T' => vec![vec![0, 1, 0], vec![1, 1, 1], vec![0, 0, 0]],
        'S' => vec![vec![0, 1, 1], vec![1, 1, 0], vec![0, 0, 0]],
        'Z' => vec![vec![1, 1, 0], vec![0, 1, 1], vec![0, 0, 0]],
        'L' => vec![vec![1, 0, 0], vec![1, 1, 1], vec![0, 0, 0]],
        'J' => vec![vec![0, 0, 1], vec![1, 1, 1], vec![0, 0, 0]],
        _ => vec![vec![1]],
    }
}

fn color_of(kind: char) -> Color {
    match kind {
        'I' => Color::new(0.0, 0.941, 0.941, 1.0),
        'O' => Color::new(0.941, 0.941, 0.0, 1.0),
        'T' => Color::new(0.706, 0.0, 0.941, 1.0),
        'S' => Color::new(0.0, 0.941, 0.0, 1.0),
        'Z' => Color::new(0.941, 0.0, 0.0, 1.0),
        'L' => Color::new(0.941, 0.627, 0.0, 1.0),
        'J' => Color::new(0.0, 0.0, 0.941, 1.0),
        _ => WHITE,
    }
}

fn all_kinds() -> [char; 7] {
    ['I', 'O', 'T', 'S', 'Z', 'L', 'J']
}

pub struct TetrisGame {
    board: Vec<Vec<Option<char>>>,
    score: i32,
    game_over: bool,
    fall_timer: f32,
    fall_delay: f32,

    current_kind: char,
    piece_shape: Vec<Vec<u8>>,
    piece_x: i32,
    piece_y: i32,

    offset: (f32, f32),
}

impl TetrisGame {
    pub fn new() -> Self {
        let mut g = Self {
            board: vec![vec![None; COLS as usize]; ROWS as usize],
            score: 0,
            game_over: false,
            fall_timer: 0.0,
            fall_delay: 0.5,
            current_kind: 'T',
            piece_shape: vec![vec![1]],
            piece_x: 0,
            piece_y: 0,
            offset: (0.0, 0.0),
        };
        // 居中：原公式 (SCREEN_WIDTH - W - 140)/2
        g.offset = ((SCREEN_WIDTH - W - 140.0) / 2.0, (SCREEN_HEIGHT - H) / 2.0);
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.board = vec![vec![None; COLS as usize]; ROWS as usize];
        self.score = 0;
        self.game_over = false;
        self.fall_timer = 0.0;
        self.spawn_piece();
    }

    fn spawn_piece(&mut self) {
        let kinds = all_kinds();
        let mut rng = ::rand::thread_rng();
        self.current_kind = kinds[rng.gen_range(0..kinds.len())];
        self.piece_shape = shape_of(self.current_kind);
        self.piece_x = COLS / 2 - self.piece_shape[0].len() as i32 / 2;
        self.piece_y = 0;
        if self.check_collision(&self.piece_shape, self.piece_x, self.piece_y) {
            self.game_over = true;
        }
    }

    fn check_collision(&self, shape: &[Vec<u8>], px: i32, py: i32) -> bool {
        for (r, row) in shape.iter().enumerate() {
            for (c, &cell) in row.iter().enumerate() {
                if cell != 0 {
                    let bx = px + c as i32;
                    let by = py + r as i32;
                    if bx < 0
                        || bx >= COLS
                        || by >= ROWS
                        || (by >= 0 && self.board[by as usize][bx as usize].is_some())
                    {
                        return true;
                    }
                }
            }
        }
        false
    }

    fn rotate_piece(&mut self) {
        let shape = &self.piece_shape;
        let n = shape.len();
        let m = shape[0].len();
        let mut rotated = vec![vec![0u8; n]; m];
        for y in 0..n {
            for x in 0..m {
                rotated[x][n - 1 - y] = shape[y][x];
            }
        }
        if !self.check_collision(&rotated, self.piece_x, self.piece_y) {
            self.piece_shape = rotated;
        }
    }

    fn lock_piece(&mut self) {
        for (r, row) in self.piece_shape.iter().enumerate() {
            for (c, &cell) in row.iter().enumerate() {
                if cell != 0 {
                    let by = self.piece_y + r as i32;
                    let bx = self.piece_x + c as i32;
                    if by >= 0 {
                        self.board[by as usize][bx as usize] = Some(self.current_kind);
                    }
                }
            }
        }
        self.clear_lines();
        self.spawn_piece();
    }

    fn clear_lines(&mut self) {
        let mut cleared = 0;
        let mut row = ROWS - 1;
        while row >= 0 {
            if self.board[row as usize].iter().all(|c| c.is_some()) {
                self.board.remove(row as usize);
                self.board.insert(0, vec![None; COLS as usize]);
                cleared += 1;
            } else {
                row -= 1;
            }
        }
        if cleared > 0 {
            let table = [0, 100, 300, 500, 800];
            self.score += table[cleared.min(4)];
        }
    }
}

impl Game for TetrisGame {
    fn update(&mut self) -> bool {
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }

        if !self.game_over {
            if is_key_pressed(KeyCode::Left) {
                if !self.check_collision(&self.piece_shape, self.piece_x - 1, self.piece_y) {
                    self.piece_x -= 1;
                }
            }
            if is_key_pressed(KeyCode::Right) {
                if !self.check_collision(&self.piece_shape, self.piece_x + 1, self.piece_y) {
                    self.piece_x += 1;
                }
            }
            if is_key_pressed(KeyCode::Down) {
                if !self.check_collision(&self.piece_shape, self.piece_x, self.piece_y + 1) {
                    self.piece_y += 1;
                }
            }
            if is_key_pressed(KeyCode::Space) || is_key_pressed(KeyCode::Up) {
                self.rotate_piece();
            }
        }

        if self.game_over {
            return true;
        }

        self.fall_timer += get_frame_time();
        if self.fall_timer >= self.fall_delay {
            self.fall_timer = 0.0;
            if !self.check_collision(&self.piece_shape, self.piece_x, self.piece_y + 1) {
                self.piece_y += 1;
            } else {
                self.lock_piece();
            }
        }

        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        draw_rectangle(ox, oy, W, H, Color::new(0.078, 0.078, 0.117, 1.0));

        // 已锁定的方块
        for y in 0..ROWS as usize {
            for x in 0..COLS as usize {
                if let Some(kind) = self.board[y][x] {
                    let color = color_of(kind);
                    draw_rounded_rect(
                        ox + x as f32 * CELL + 1.0,
                        oy + y as f32 * CELL + 1.0,
                        CELL - 2.0,
                        CELL - 2.0,
                        4.0,
                        color,
                    );
                }
            }
        }

        // 当前方块
        if !self.game_over {
            let color = color_of(self.current_kind);
            for (r, row) in self.piece_shape.iter().enumerate() {
                for (c, &cell) in row.iter().enumerate() {
                    if cell != 0 {
                        let px = self.piece_x + c as i32;
                        let py = self.piece_y + r as i32;
                        draw_rounded_rect(
                            ox + px as f32 * CELL + 1.0,
                            oy + py as f32 * CELL + 1.0,
                            CELL - 2.0,
                            CELL - 2.0,
                            4.0,
                            color,
                        );
                    }
                }
            }
        }

        // 分数（右侧）
        draw_text(
            &format!("Score: {}", self.score),
            ox + W + 20.0,
            oy + 50.0,
            32.0,
            WHITE,
        );

        // Game Over
        if self.game_over {
            draw_overlay(ox, oy, W, H, 0.7, BLACK);
            draw_text_centered("GAME OVER", ox + W / 2.0, oy + H / 2.0 - 20.0, 48, WHITE);
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 30.0, 32, WHITE);
        }
    }
}