use crate::common::*;
use macroquad::prelude::*;

const W: f32 = 500.0;
const H: f32 = 500.0;
const CELL: f32 = W / 3.0;

pub struct TicTacToeGame {
    board: [[char; 3]; 3],
    current_player: char,
    winner: char,
    game_over: bool,
    move_count: i32,
    offset: (f32, f32),
}

impl TicTacToeGame {
    pub fn new() -> Self {
        let mut g = Self {
            board: [['\0'; 3]; 3],
            current_player: 'X',
            winner: '\0',
            game_over: false,
            move_count: 0,
            offset: centered_offset(W, H),
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.board = [['\0'; 3]; 3];
        self.current_player = 'X';
        self.winner = '\0';
        self.game_over = false;
        self.move_count = 0;
    }

    fn check_winner(&self) -> char {
        for r in 0..3 {
            if self.board[r][0] != '\0'
                && self.board[r][0] == self.board[r][1]
                && self.board[r][1] == self.board[r][2]
            {
                return self.board[r][0];
            }
        }
        for c in 0..3 {
            if self.board[0][c] != '\0'
                && self.board[0][c] == self.board[1][c]
                && self.board[1][c] == self.board[2][c]
            {
                return self.board[0][c];
            }
        }
        if self.board[0][0] != '\0'
            && self.board[0][0] == self.board[1][1]
            && self.board[1][1] == self.board[2][2]
        {
            return self.board[0][0];
        }
        if self.board[0][2] != '\0'
            && self.board[0][2] == self.board[1][1]
            && self.board[1][1] == self.board[2][0]
        {
            return self.board[0][2];
        }
        '\0'
    }
}

impl Game for TicTacToeGame {
    fn update(&mut self) -> bool {
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }

        if !self.game_over && is_mouse_button_pressed(MouseButton::Left) {
            let (mx, my) = mouse_position();
            let (ox, oy) = self.offset;
            let x = ((mx - ox) / CELL).floor() as i32;
            let y = ((my - oy) / CELL).floor() as i32;
            if x >= 0 && x < 3 && y >= 0 && y < 3
                && self.board[y as usize][x as usize] == '\0'
            {
                self.board[y as usize][x as usize] = self.current_player;
                self.move_count += 1;
                self.winner = self.check_winner();
                if self.winner != '\0' {
                    self.game_over = true;
                } else if self.move_count == 9 {
                    self.game_over = true;
                } else {
                    self.current_player = if self.current_player == 'X' { 'O' } else { 'X' };
                }
            }
        }
        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        draw_rectangle(ox, oy, W, H, Color::new(0.118, 0.118, 0.165, 1.0));

        // 格子 + 落子
        for r in 0..3 {
            for c in 0..3 {
                let rx = ox + c as f32 * CELL + 2.0;
                let ry = oy + r as f32 * CELL + 2.0;
                let rw = CELL - 4.0;
                let rh = CELL - 4.0;
                draw_rounded_rect(rx, ry, rw, rh, 8.0, Color::new(0.196, 0.196, 0.255, 1.0));

                let ch = self.board[r][c];
                if ch != '\0' {
                    let color = if ch == 'X' {
                        Color::new(0.235, 0.784, 1.0, 1.0)
                    } else {
                        Color::new(1.0, 0.784, 0.235, 1.0)
                    };
                    let s = ch.to_string();
                    draw_text_centered(
                        &s,
                        rx + CELL / 2.0 - 2.0,
                        ry + CELL / 2.0 - 2.0,
                        80,
                        color,
                    );
                }
            }
        }

        // 网格线
        for i in 1..3 {
            draw_line(
                ox + i as f32 * CELL,
                oy + 5.0,
                ox + i as f32 * CELL,
                oy + H - 5.0,
                3.0,
                Color::new(0.314, 0.314, 0.392, 1.0),
            );
            draw_line(
                ox + 5.0,
                oy + i as f32 * CELL,
                ox + W - 5.0,
                oy + i as f32 * CELL,
                3.0,
                Color::new(0.314, 0.314, 0.392, 1.0),
            );
        }

        // 状态文本
        let status_y = oy + H + 20.0;
        if self.game_over {
            if self.winner != '\0' {
                let msg = format!("Player {} Wins!", self.winner);
                draw_text_centered(
                    &msg,
                    SCREEN_WIDTH / 2.0,
                    status_y,
                    48,
                    Color::new(0.0, 1.0, 0.0, 1.0),
                );
            } else {
                draw_text_centered(
                    "Draw!",
                    SCREEN_WIDTH / 2.0,
                    status_y,
                    48,
                    Color::new(1.0, 1.0, 0.392, 1.0),
                );
            }
            draw_text_centered(
                "Press R to restart",
                SCREEN_WIDTH / 2.0,
                status_y + 40.0,
                28,
                Color::new(0.784, 0.784, 0.784, 1.0),
            );
        } else {
            let msg = format!("Player {}'s Turn", self.current_player);
            draw_text_centered(&msg, SCREEN_WIDTH / 2.0, status_y, 48, WHITE);
        }
    }
}