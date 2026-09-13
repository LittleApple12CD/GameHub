use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;

const W: f32 = 600.0;
const H: f32 = 600.0;
const CELL: f32 = 20.0;

pub struct SnakeGame {
    snake: Vec<(i32, i32)>,
    direction: (i32, i32),
    next_direction: (i32, i32),
    food: (i32, i32),
    score: i32,
    game_over: bool,
    move_timer: f32,
    move_delay: f32,
    offset: (f32, f32),
    gw: i32,
    gh: i32,
}

impl SnakeGame {
    pub fn new() -> Self {
        let mut g = Self {
            snake: vec![],
            direction: (1, 0),
            next_direction: (1, 0),
            food: (0, 0),
            score: 0,
            game_over: false,
            move_timer: 0.0,
            move_delay: 0.15,
            offset: centered_offset(W, H),
            gw: (W / CELL) as i32,
            gh: (H / CELL) as i32,
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.snake = vec![(self.gw / 2, self.gh / 2)];
        self.direction = (1, 0);
        self.next_direction = (1, 0);
        self.food = self.spawn_food();
        self.score = 0;
        self.game_over = false;
        self.move_timer = 0.0;
    }

    fn spawn_food(&self) -> (i32, i32) {
        let mut rng = ::rand::thread_rng();
        loop {
            let p = (rng.gen_range(0..self.gw), rng.gen_range(0..self.gh));
            if !self.snake.contains(&p) {
                return p;
            }
        }
    }
}

impl Game for SnakeGame {
    fn update(&mut self) -> bool {
        // 输入
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }
        if is_key_pressed(KeyCode::Up)    && self.direction != (0, 1)  { self.next_direction = (0, -1); }
        if is_key_pressed(KeyCode::Down)  && self.direction != (0, -1) { self.next_direction = (0, 1); }
        if is_key_pressed(KeyCode::Left)  && self.direction != (1, 0)  { self.next_direction = (-1, 0); }
        if is_key_pressed(KeyCode::Right) && self.direction != (-1, 0) { self.next_direction = (1, 0); }

        if self.game_over {
            return true;
        }

        // 逻辑
        self.move_timer += get_frame_time();
        if self.move_timer >= self.move_delay {
            self.move_timer = 0.0;
            self.direction = self.next_direction;
            let head = self.snake[0];
            let new_head = (head.0 + self.direction.0, head.1 + self.direction.1);

            if new_head.0 < 0 || new_head.0 >= self.gw
                || new_head.1 < 0 || new_head.1 >= self.gh
                || self.snake.contains(&new_head)
            {
                self.game_over = true;
                return true;
            }

            self.snake.insert(0, new_head);
            if new_head == self.food {
                self.score += 10;
                self.food = self.spawn_food();
            } else {
                self.snake.pop();
            }
        }
        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        // 游戏面板
        draw_rectangle(ox, oy, W, H, Color::new(0.078, 0.078, 0.117, 1.0));

        // 蛇身
        for (i, seg) in self.snake.iter().enumerate() {
            let color = if i == 0 {
                Color::new(0.235, 0.862, 0.235, 1.0)
            } else {
                Color::new(0.156, 0.706, 0.156, 1.0)
            };
            draw_rounded_rect(
                ox + seg.0 as f32 * CELL + 1.0,
                oy + seg.1 as f32 * CELL + 1.0,
                CELL - 2.0,
                CELL - 2.0,
                4.0,
                color,
            );
        }

        // 食物
        draw_rounded_rect(
            ox + self.food.0 as f32 * CELL + 1.0,
            oy + self.food.1 as f32 * CELL + 1.0,
            CELL - 2.0,
            CELL - 2.0,
            6.0,
            Color::new(1.0, 0.235, 0.235, 1.0),
        );

        // 分数
        draw_text(&format!("Score: {}", self.score), ox + 10.0, oy + 30.0, 36.0, WHITE);

        // Game Over
        if self.game_over {
            draw_overlay(ox, oy, W, H, 0.7, BLACK);
            draw_text_centered(
                "GAME OVER - Press R",
                ox + W / 2.0,
                oy + H / 2.0,
                36,
                WHITE,
            );
        }
    }
}