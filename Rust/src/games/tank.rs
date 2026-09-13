use crate::common::*;
use macroquad::prelude::*;
use ::rand::Rng;
use std::collections::HashSet;

const W: f32 = 600.0;
const H: f32 = 600.0;
const CELL: f32 = 30.0;
const GRID: i32 = 20;
const BULLET_SPEED: i32 = 1;

#[derive(Clone)]
struct Tank {
    x: i32,
    y: i32,
    dir: (i32, i32),
    alive: bool,
    timer: i32,
    move_timer: i32,
}

struct Bullet {
    x: f32,
    y: f32,
    dir: (i32, i32),
    is_player: bool,
    alive: bool,
}

pub struct TankGame {
    walls: HashSet<(i32, i32)>,
    player: Tank,
    enemies: Vec<Tank>,
    bullets: Vec<Bullet>,
    score: i32,
    game_over: bool,
    enemy_spawn_timer: i32,
    enemy_spawn_delay: i32,
    move_delay: f32,
    enemy_move_delay: i32,
    offset: (f32, f32),
}

impl TankGame {
    pub fn new() -> Self {
        let mut g = Self {
            walls: HashSet::new(),
            player: Tank { x: 1, y: 1, dir: (1, 0), alive: true, timer: 0, move_timer: 0 },
            enemies: Vec::new(),
            bullets: Vec::new(),
            score: 0,
            game_over: false,
            enemy_spawn_timer: 0,
            enemy_spawn_delay: 120,
            move_delay: 0.2,
            enemy_move_delay: 250,
            offset: centered_offset(W, H),
        };
        g.reset();
        g
    }

    fn reset(&mut self) {
        self.walls = self.generate_walls();
        self.player = Tank { x: 1, y: 1, dir: (1, 0), alive: true, timer: 0, move_timer: 0 };
        self.enemies = vec![Tank {
            x: GRID - 2,
            y: GRID - 2,
            dir: (-1, 0),
            alive: true,
            timer: 0,
            move_timer: 0,
        }];
        self.bullets.clear();
        self.score = 0;
        self.game_over = false;
        self.enemy_spawn_timer = 0;
    }

    fn generate_walls(&self) -> HashSet<(i32, i32)> {
        let mut rng = ::rand::thread_rng();
        let mut walls = HashSet::new();
        for i in 0..GRID {
            walls.insert((i, 0));
            walls.insert((i, GRID - 1));
            walls.insert((0, i));
            walls.insert((GRID - 1, i));
        }
        for _ in 0..12 {
            let x = rng.gen_range(2..GRID - 2);
            let y = rng.gen_range(2..GRID - 2);
            if (x, y) != (1, 1) && (x, y) != (GRID - 2, GRID - 2) {
                walls.insert((x, y));
            }
        }
        walls
    }

    fn can_move(&self, x: i32, y: i32, is_player: bool) -> bool {
        if self.walls.contains(&(x, y)) {
            return false;
        }
        if is_player {
            for e in &self.enemies {
                if e.alive && e.x == x && e.y == y {
                    return false;
                }
            }
        } else {
            if self.player.alive && self.player.x == x && self.player.y == y {
                return false;
            }
            for e in &self.enemies {
                if e.alive && e.x == x && e.y == y {
                    return false;
                }
            }
        }
        x >= 0 && x < GRID && y >= 0 && y < GRID
    }

    fn shoot_bullet(&mut self, x: i32, y: i32, dir: (i32, i32), is_player: bool) {
        self.bullets.push(Bullet {
            x: (x + dir.0) as f32,
            y: (y + dir.1) as f32,
            dir,
            is_player,
            alive: true,
        });
    }
}

impl Game for TankGame {
    fn update(&mut self) -> bool {
        if is_key_pressed(KeyCode::Escape) {
            return false;
        }
        if is_key_pressed(KeyCode::R) {
            self.reset();
        }
        if is_key_pressed(KeyCode::Space) && !self.game_over && self.player.alive {
            let d = self.player.dir;
            let (px, py) = (self.player.x, self.player.y);
            self.shoot_bullet(px, py, d, true);
        }

        // 玩家移动
        if !self.game_over && self.player.alive {
            let mut dx = 0;
            let mut dy = 0;
            if is_key_down(KeyCode::W) {
                dy = -1;
            } else if is_key_down(KeyCode::S) {
                dy = 1;
            } else if is_key_down(KeyCode::A) {
                dx = -1;
            } else if is_key_down(KeyCode::D) {
                dx = 1;
            }

            if dx != 0 || dy != 0 {
                self.player.dir = (dx, dy);
                self.player.move_timer += (get_frame_time() * 1000.0) as i32;
                if self.player.move_timer as f32 >= self.move_delay * 1000.0 {
                    self.player.move_timer = 0;
                    let nx = self.player.x + dx;
                    let ny = self.player.y + dy;
                    if self.can_move(nx, ny, true) {
                        self.player.x = nx;
                        self.player.y = ny;
                    }
                }
            } else {
                self.player.move_timer = 0;
            }
        }

        if self.game_over {
            return true;
        }

        // 子弹
        let mut new_bullets: Vec<Bullet> = Vec::new();
        for mut b in self.bullets.drain(..) {
            if !b.alive {
                continue;
            }
            let dx = b.dir.0 * BULLET_SPEED;
            let dy = b.dir.1 * BULLET_SPEED;
            let steps = dx.abs().max(dy.abs()) + 1;
            let step_x = dx as f32 / steps as f32;
            let step_y = dy as f32 / steps as f32;

            let mut hit = false;
            for _ in 0..steps {
                b.x += step_x;
                b.y += step_y;
                let bx = b.x.round() as i32;
                let by = b.y.round() as i32;

                if bx < 0 || bx >= GRID || by < 0 || by >= GRID {
                    hit = true;
                    break;
                }
                if self.walls.contains(&(bx, by)) {
                    hit = true;
                    break;
                }
                if b.is_player {
                    let mut enemy_hit = false;
                    for e in &mut self.enemies {
                        if e.alive && e.x == bx && e.y == by {
                            e.alive = false;
                            self.score += 10;
                            enemy_hit = true;
                            break;
                        }
                    }
                    if enemy_hit {
                        hit = true;
                        break;
                    }
                } else {
                    if self.player.alive && self.player.x == bx && self.player.y == by {
                        self.player.alive = false;
                        self.game_over = true;
                        hit = true;
                        break;
                    }
                }
            }
            if !hit {
                new_bullets.push(b);
            }
        }
        self.bullets = new_bullets;

        // 敌人 AI
        let mut rng = ::rand::thread_rng();
        let dirs = [(1, 0), (-1, 0), (0, 1), (0, -1)];
        let mut new_enemies_shots: Vec<(i32, i32, (i32, i32))> = Vec::new();
        for i in 0..self.enemies.len() {
            if !self.enemies[i].alive {
                continue;
            }
            self.enemies[i].timer += 1;
            if self.enemies[i].timer >= 20 {
                self.enemies[i].timer = 0;
                if rng.gen_bool(0.25) {
                    self.enemies[i].dir = dirs[rng.gen_range(0..dirs.len())];
                }
                self.enemies[i].move_timer += 20;
                if self.enemies[i].move_timer >= self.enemy_move_delay {
                    self.enemies[i].move_timer = 0;
                    let nx = self.enemies[i].x + self.enemies[i].dir.0;
                    let ny = self.enemies[i].y + self.enemies[i].dir.1;
                    // 注意：can_move 里对敌检查不区分"自己"，简化处理即可
                    let mut blocked = self.walls.contains(&(nx, ny));
                    if !blocked {
                        if self.player.alive && self.player.x == nx && self.player.y == ny {
                            blocked = true;
                        }
                        for (j, e2) in self.enemies.iter().enumerate() {
                            if j != i && e2.alive && e2.x == nx && e2.y == ny {
                                blocked = true;
                                break;
                            }
                        }
                    }
                    if !blocked {
                        self.enemies[i].x = nx;
                        self.enemies[i].y = ny;
                    }
                }
            }
            if rng.gen_bool(0.015) {
                let (ex, ey, ed) = (self.enemies[i].x, self.enemies[i].y, self.enemies[i].dir);
                new_enemies_shots.push((ex, ey, ed));
            }
        }
        for (ex, ey, ed) in new_enemies_shots {
            self.shoot_bullet(ex, ey, ed, false);
        }

        // 刷敌
        self.enemy_spawn_timer += 1;
        if self.enemy_spawn_timer >= self.enemy_spawn_delay {
            self.enemy_spawn_timer = 0;
            let alive_count = self.enemies.iter().filter(|e| e.alive).count();
            if alive_count < 3 {
                self.enemies.push(Tank {
                    x: GRID - 2,
                    y: GRID - 2,
                    dir: (-1, 0),
                    alive: true,
                    timer: 0,
                    move_timer: 0,
                });
            }
        }

        true
    }

    fn draw(&mut self) {
        clear_background(COLOR_BG);
        let (ox, oy) = self.offset;

        draw_rectangle(ox, oy, W, H, Color::new(0.118, 0.118, 0.165, 1.0));

        // 墙
        for &(x, y) in &self.walls {
            draw_rounded_rect(
                ox + x as f32 * CELL + 1.0,
                oy + y as f32 * CELL + 1.0,
                CELL - 2.0,
                CELL - 2.0,
                8.0,
                Color::new(0.392, 0.392, 0.510, 1.0),
            );
        }

        // 玩家坦克
        if self.player.alive {
            draw_tank(ox, oy, self.player.x, self.player.y,
                      Color::new(0.196, 0.902, 0.196, 1.0),
                      Color::new(0.118, 0.784, 0.118, 1.0),
                      self.player.dir);
        }
        // 敌人
        for e in &self.enemies {
            if e.alive {
                draw_tank(ox, oy, e.x, e.y,
                          Color::new(0.902, 0.196, 0.196, 1.0),
                          Color::new(0.784, 0.118, 0.118, 1.0),
                          e.dir);
            }
        }

        // 子弹
        for b in &self.bullets {
            let color = if b.is_player {
                Color::new(1.0, 1.0, 0.314, 1.0)
            } else {
                Color::new(1.0, 0.588, 0.196, 1.0)
            };
            draw_rounded_rect(
                ox + b.x * CELL + 10.0,
                oy + b.y * CELL + 10.0,
                10.0,
                10.0,
                5.0,
                color,
            );
        }

        // 分数
        draw_text(&format!("Score: {}", self.score), ox + 10.0, oy + 30.0, 36.0, WHITE);

        // Game Over
        if self.game_over {
            draw_overlay(ox, oy, W, H, 0.6, BLACK);
            draw_text_centered("GAME OVER", ox + W / 2.0, oy + H / 2.0 - 20.0, 48, WHITE);
            draw_text_centered("Press R to restart", ox + W / 2.0, oy + H / 2.0 + 30.0, 36, WHITE);
        }
    }
}

fn draw_tank(ox: f32, oy: f32, x: i32, y: i32, color: Color, gun_color: Color, dir: (i32, i32)) {
    let cx = ox + x as f32 * CELL + CELL / 2.0;
    let cy = oy + y as f32 * CELL + CELL / 2.0;

    draw_rounded_rect(
        ox + x as f32 * CELL + 2.0,
        oy + y as f32 * CELL + 2.0,
        CELL - 4.0,
        CELL - 4.0,
        8.0,
        color,
    );

    let (dx, dy) = dir;
    let gun = if dx == 1 {
        Rect::new(cx + 4.0, cy - 4.0, 14.0, 8.0)
    } else if dx == -1 {
        Rect::new(cx - 18.0, cy - 4.0, 14.0, 8.0)
    } else if dy == -1 {
        Rect::new(cx - 4.0, cy - 18.0, 8.0, 14.0)
    } else {
        Rect::new(cx - 4.0, cy + 4.0, 8.0, 14.0)
    };
    draw_rounded_rect(gun.x, gun.y, gun.w, gun.h, 4.0, gun_color);
    draw_circle(cx, cy, 6.0, color);
}