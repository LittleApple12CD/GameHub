use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;

const W: f32 = 600.0;
const H: f32 = 600.0;

struct Brick {
    rect: Rect,
    color: Color,
    alive: bool,
}

pub struct BreakoutGame {
    paddle: Rect,
    paddle_speed: f32,
    ball: Rect,
    ball_dx: f32,
    ball_dy: f32,
    ball_speed: f32,
    bricks: Vec<Brick>,
    score: i32,
    lives: i32,
    game_over: bool,
    waiting: bool,
    offset: (f32, f32),
}

impl BreakoutGame {
    pub fn new() -> Self {
        let mut g = Self {
            paddle: Rect::new(0.0, 0.0, 120.0, 16.0),
            paddle_speed: 6.0,
            ball: Rect::new(0.0, 0.0, 20.0, 20.0),
            ball_dx: 5.0,
            ball_dy: -6.0,
            ball_speed: 4.0,
            bricks: Vec::new(),
            score: 0,
            lives: 3,
            game_over: false,
            waiting: true,
            offset: centered_offset(W, H),
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.paddle = Rect::new(W / 2.0 - 60.0, H - 40.0, 120.0, 16.0);
        self.ball = Rect::new(W / 2.0 - 10.0, H - 70.0, 20.0, 20.0);
        self.ball_dx = 5.0;
        self.ball_dy = -6.0;
        self.score = 0;
        self.lives = 3;
        self.game_over = false;
        self.waiting = true;

        self.bricks.clear();
        let rows = 5;
        let cols = 8;
        let brick_w = (W - 20.0) / cols as f32 - 4.0;
        let brick_h = 22.0;
        let colors = [
            Color::new(0.902, 0.196, 0.196, 1.0),
            Color::new(0.902, 0.588, 0.196, 1.0),
            Color::new(0.902, 0.902, 0.196, 1.0),
            Color::new(0.196, 0.902, 0.196, 1.0),
            Color::new(0.196, 0.588, 0.902, 1.0),
        ];
        for row in 0..rows {
            for col in 0..cols {
                let x = 10.0 + col as f32 * (brick_w + 4.0);
                let y = 40.0 + row as f32 * (brick_h + 4.0);
                self.bricks.push(Brick {
                    rect: Rect::new(x, y, brick_w, brick_h),
                    color: colors[row % colors.len()],
                    alive: true,
                });
            }
        }
    }
}

impl Game for BreakoutGame {
    fn update(&mut self) -> bool {
        // ---- 事件 ----
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }
        if is_key_pressed(KeyCode::Space) && self.waiting && !self.game_over {
            self.waiting = false;
        }

        if !self.game_over && !self.waiting {
            if is_key_down(KeyCode::Left) && self.paddle.x > 0.0 {
                self.paddle.x -= self.paddle_speed;
            }
            if is_key_down(KeyCode::Right) && self.paddle.x + self.paddle.w < W {
                self.paddle.x += self.paddle_speed;
            }
        }

        if self.game_over || self.waiting {
            return true;
        }

        // ---- 逻辑 ----
        self.ball.x += self.ball_dx;
        self.ball.y += self.ball_dy;

        if self.ball.x <= 0.0 || self.ball.x + self.ball.w >= W {
            self.ball_dx = -self.ball_dx;
        }
        if self.ball.y <= 0.0 {
            self.ball_dy = -self.ball_dy;
        }

        // 掉落
        if self.ball.y + self.ball.h >= H {
            self.lives -= 1;
            if self.lives <= 0 {
                self.game_over = true;
            } else {
                self.waiting = true;
                self.ball.x = W / 2.0 - 10.0;
                self.ball.y = H - 70.0;
                let mut rng = ::rand::thread_rng();
                let sign = if rng.gen_bool(0.5) { 1.0 } else { -1.0 };
                self.ball_dx = self.ball_speed * sign;
                self.ball_dy = -self.ball_speed;
            }
            return true;
        }

        // 挡板碰撞
        if self.ball.overlaps(&self.paddle) {
            self.ball_dy = -self.ball_dy.abs();
            let hit_pos = (self.ball.x + self.ball.w / 2.0
                - (self.paddle.x + self.paddle.w / 2.0))
                / (self.paddle.w / 2.0);
            self.ball_dx = hit_pos * self.ball_speed * 0.9;
            if self.ball_dx.abs() < 1.5 {
                self.ball_dx = if self.ball_dx >= 0.0 { 3.0 } else { -3.0 };
            }
        }

        // 砖块碰撞
        for i in 0..self.bricks.len() {
            if !self.bricks[i].alive {
                continue;
            }
            let br = self.bricks[i].rect;
            if self.ball.overlaps(&br) {
                self.bricks[i].alive = false;
                self.score += 10;

                let overlap_top = (self.ball.y + self.ball.h) - br.y;
                let overlap_bottom = (br.y + br.h) - self.ball.y;
                let overlap_left = (self.ball.x + self.ball.w) - br.x;
                let overlap_right = (br.x + br.w) - self.ball.x;
                let min_overlap = overlap_top
                    .min(overlap_bottom)
                    .min(overlap_left)
                    .min(overlap_right);

                if min_overlap == overlap_top || min_overlap == overlap_bottom {
                    self.ball_dy = -self.ball_dy;
                } else {
                    self.ball_dx = -self.ball_dx;
                }
                break;
            }
        }

        if self.bricks.iter().all(|b| !b.alive) {
            self.game_over = true;
        }

        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        draw_rectangle(ox, oy, W, H, Color::new(0.078, 0.078, 0.117, 1.0));

        // 砖块
        for b in &self.bricks {
            if b.alive {
                draw_rounded_rect(ox + b.rect.x, oy + b.rect.y, b.rect.w, b.rect.h, 8.0, b.color);
                draw_rounded_rect_lines(ox + b.rect.x, oy + b.rect.y, b.rect.w, b.rect.h, 8.0, 1.0, WHITE);
            }
        }

        // 挡板
        draw_rounded_rect(ox + self.paddle.x, oy + self.paddle.y, self.paddle.w, self.paddle.h, 8.0, WHITE);
        draw_rounded_rect_lines(ox + self.paddle.x, oy + self.paddle.y, self.paddle.w, self.paddle.h, 8.0, 2.0,
                                Color::new(0.784, 0.784, 0.863, 1.0));

        // 球
        draw_circle(
            ox + self.ball.x + self.ball.w / 2.0,
            oy + self.ball.y + self.ball.h / 2.0,
            self.ball.w / 2.0,
            Color::new(1.0, 1.0, 0.471, 1.0),
        );

        // 分数 / 生命
        draw_text(&format!("Score: {}", self.score), ox + 10.0, oy + 30.0, 36.0, WHITE);
        let lives_text = format!("Lives: {}", self.lives);
        let dims = measure_text(&lives_text, None, 36, 1.0);
        draw_text(&lives_text, ox + W - dims.width - 10.0, oy + 30.0, 36.0, WHITE);

        // 提示
        if self.waiting && !self.game_over {
            draw_text_centered("Press SPACE", ox + W / 2.0, oy + H / 2.0 + 30.0, 48, WHITE);
            draw_text_centered("to start", ox + W / 2.0, oy + H / 2.0 + 70.0, 36, WHITE);
        } else if self.game_over {
            draw_overlay(ox, oy, W, H, 0.6, BLACK);
            let msg = if self.lives <= 0 { "GAME OVER" } else { "YOU WIN!" };
            let col = if self.lives <= 0 {
                WHITE
            } else {
                Color::new(0.0, 1.0, 0.0, 1.0)
            };
            draw_text_centered(msg, ox + W / 2.0, oy + H / 2.0 - 20.0, 48, col);
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 30.0, 36, WHITE);
        }
    }
}