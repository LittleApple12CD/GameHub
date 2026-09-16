use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;

const W: f32 = 400.0;
const H: f32 = 600.0;
const BIRD_X: f32 = 80.0;

struct Pipe {
    x: f32,
    gap_y: f32,
    passed: bool,
}

pub struct FlappyGame {
    bird_y: f32,
    bird_vy: f32,
    bird_radius: f32,
    gravity: f32,
    jump_strength: f32,
    pipes: Vec<Pipe>,
    pipe_width: f32,
    pipe_gap: f32,
    pipe_speed: f32,
    score: i32,
    game_over: bool,
    pipe_timer: i32,
    pipe_delay: i32,
    offset: (f32, f32),
}

impl FlappyGame {
    pub fn new() -> Self {
        let mut g = Self {
            bird_y: H / 2.0,
            bird_vy: 0.0,
            bird_radius: 15.0,
            gravity: 0.2,
            jump_strength: -7.0,
            pipes: Vec::new(),
            pipe_width: 60.0,
            pipe_gap: 160.0,
            pipe_speed: 2.0,
            score: 0,
            game_over: false,
            pipe_timer: 0,
            pipe_delay: 90,
            offset: centered_offset(W, H),
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.bird_y = H / 2.0;
        self.bird_vy = 0.0;
        self.pipes.clear();
        self.score = 0;
        self.game_over = false;
        self.pipe_timer = 0;
    }

    fn add_pipe(&mut self) {
        let mut rng = ::rand::thread_rng();
        let gap_y = rng.gen_range(80.0..=(H - 80.0 - self.pipe_gap));
        self.pipes.push(Pipe {
            x: W,
            gap_y,
            passed: false,
        });
    }
}

impl Game for FlappyGame {
    fn update(&mut self) -> bool {
        // ---- 事件 ----
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }
        if is_key_pressed(KeyCode::Space) && !self.game_over {
            self.bird_vy = self.jump_strength;
        }
        if is_mouse_button_pressed(MouseButton::Left) && !self.game_over {
            self.bird_vy = self.jump_strength;
        }

        if self.game_over {
            return true;
        }

        // ---- 逻辑 ----
        self.bird_vy += self.gravity;
        self.bird_y += self.bird_vy;

        if self.bird_y - self.bird_radius < 0.0 {
            self.bird_y = self.bird_radius;
            self.bird_vy = 0.0;
        } else if self.bird_y + self.bird_radius > H {
            self.game_over = true;
            return true;
        }

        self.pipe_timer += 1;
        if self.pipe_timer >= self.pipe_delay {
            self.pipe_timer = 0;
            self.add_pipe();
        }

        let mut remove_idx: Vec<usize> = Vec::new();
        for i in 0..self.pipes.len() {
            self.pipes[i].x -= self.pipe_speed;

            if !self.pipes[i].passed && self.pipes[i].x + self.pipe_width < BIRD_X {
                self.pipes[i].passed = true;
                self.score += 1;
            }

            // 碰撞
            let p = &self.pipes[i];
            if p.x < BIRD_X + self.bird_radius && p.x + self.pipe_width > BIRD_X - self.bird_radius {
                if self.bird_y - self.bird_radius < p.gap_y
                    || self.bird_y + self.bird_radius > p.gap_y + self.pipe_gap
                {
                    self.game_over = true;
                    return true;
                }
            }

            if p.x + self.pipe_width < 0.0 {
                remove_idx.push(i);
            }
        }
        for i in remove_idx.into_iter().rev() {
            self.pipes.remove(i);
        }

        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        // 天空
        draw_rectangle(ox, oy, W, H, Color::new(0.529, 0.808, 0.922, 1.0));

        // 鸟
        let bx = ox + BIRD_X;
        let by = oy + self.bird_y;
        draw_circle(bx, by, self.bird_radius, Color::new(1.0, 1.0, 0.196, 1.0));
        draw_circle(bx + 8.0, by - 5.0, 4.0, BLACK);
        draw_circle(bx + 10.0, by - 7.0, 2.0, WHITE);

        // 管道
        for p in &self.pipes {
            let top = Rect::new(ox + p.x, oy, self.pipe_width, p.gap_y);
            draw_rounded_rect(top.x, top.y, top.w, top.h, 6.0, Color::new(0.157, 0.784, 0.157, 1.0));
            draw_rounded_rect(
                ox + p.x - 6.0,
                oy + p.gap_y - 22.0,
                self.pipe_width + 12.0,
                22.0,
                6.0,
                Color::new(0.118, 0.627, 0.118, 1.0),
            );

            let bottom_y = p.gap_y + self.pipe_gap;
            draw_rounded_rect(
                ox + p.x,
                oy + bottom_y,
                self.pipe_width,
                H - bottom_y,
                6.0,
                Color::new(0.157, 0.784, 0.157, 1.0),
            );
            draw_rounded_rect(
                ox + p.x - 6.0,
                oy + bottom_y,
                self.pipe_width + 12.0,
                22.0,
                6.0,
                Color::new(0.118, 0.627, 0.118, 1.0),
            );
        }

        // 分数
        draw_text_centered(
            &self.score.to_string(),
            ox + W / 2.0,
            oy + 50.0,
            40,
            WHITE,
        );

        // Game Over
        if self.game_over {
            draw_overlay(ox, oy, W, H, 0.6, BLACK);
            draw_text_centered("GAME OVER", ox + W / 2.0, oy + H / 2.0 - 30.0, 40, WHITE);
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 20.0, 26, WHITE);
        }
    }
}
