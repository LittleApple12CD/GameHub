use macroquad::prelude::*;
use macroquad::text::{load_ttf_font_from_bytes, TextParams};
use std::sync::OnceLock;

static FONT: OnceLock<Font> = OnceLock::new();

pub fn init_font() {
    let bytes = include_bytes!("../assets/fonts/arial.ttf");
    match load_ttf_font_from_bytes(bytes) {
        Ok(f) => { let _ = FONT.set(f); }
        Err(e) => eprintln!("Error: {e}"),
    }
}

fn font() -> Option<&'static Font> {
    FONT.get()
}

// ---------- 窗口尺寸 ----------
pub const SCREEN_WIDTH: f32 = 850.0;
pub const SCREEN_HEIGHT: f32 = 650.0;

// ---------- 主菜单配色 ----------
pub const COLOR_BG: Color         = Color::new(0.098, 0.098, 0.137, 1.0);
pub const COLOR_WHITE: Color      = WHITE;
pub const COLOR_HOVER: Color      = Color::new(0.314, 0.784, 1.0, 1.0);
pub const COLOR_BUTTON: Color     = Color::new(0.216, 0.216, 0.294, 1.0);
pub const COLOR_BORDER: Color     = Color::new(0.471, 0.471, 0.627, 1.0);
pub const COLOR_QUIT: Color       = Color::new(0.784, 0.196, 0.196, 1.0);
pub const COLOR_QUIT_HOVER: Color = Color::new(1.0, 0.314, 0.314, 1.0);

pub trait Game {
    fn update(&mut self) -> bool;
    fn draw(&mut self);
}

pub fn centered_offset(w: f32, h: f32) -> (f32, f32) {
    ((SCREEN_WIDTH - w) / 2.0, (SCREEN_HEIGHT - h) / 2.0)
}

pub fn draw_text_centered(
    text: &str,
    center_x: f32,
    center_y: f32,
    font_size: u16,
    color: Color,
) {
    let dims = measure_text(text, font(), font_size, 1.0);
    draw_text_ex(
        text,
        center_x - dims.width / 2.0,
        center_y + dims.height / 2.0,
        TextParams {
            font: font(),
            font_size,
            color,
            ..Default::default()
        },
    );
}

pub fn draw_overlay(x: f32, y: f32, w: f32, h: f32, alpha: f32, color: Color) {
    let mut c = color;
    c.a = alpha;
    draw_rectangle(x, y, w, h, c);
}

use std::f32::consts::PI;

pub fn draw_rounded_rect(x: f32, y: f32, w: f32, h: f32, r: f32, color: Color) {
    let r = r.min(w / 2.0).min(h / 2.0).max(0.0);

    draw_rectangle(x, y + r, w, h - 2.0 * r, color);
    draw_rectangle(x + r, y, w - 2.0 * r, r, color);
    draw_rectangle(x + r, y + h - r, w - 2.0 * r, r, color);
    draw_quarter_circle(x + r,         y + r,         r, PI,          PI * 1.5, color); // 左上
    draw_quarter_circle(x + w - r,     y + r,         r, PI * 1.5,    PI * 2.0, color); // 右上
    draw_quarter_circle(x + w - r,     y + h - r,     r, 0.0,         PI * 0.5, color); // 右下
    draw_quarter_circle(x + r,         y + h - r,     r, PI * 0.5,    PI,       color); // 左下
}

pub fn draw_rounded_rect_lines(x: f32, y: f32, w: f32, h: f32, r: f32, thickness: f32, color: Color) {
    let r = r.min(w / 2.0).min(h / 2.0).max(0.0);

    draw_line(x + r, y,         x + w - r, y,         thickness, color); // 上
    draw_line(x + r, y + h,     x + w - r, y + h,     thickness, color); // 下
    draw_line(x,     y + r,     x,         y + h - r, thickness, color); // 左
    draw_line(x + w, y + r,     x + w,     y + h - r, thickness, color); // 右

    draw_arc_lines(x + r,     y + r,     r, PI,       PI * 1.5, thickness, color); // 左上
    draw_arc_lines(x + w - r, y + r,     r, PI * 1.5, PI * 2.0, thickness, color); // 右上
    draw_arc_lines(x + w - r, y + h - r, r, 0.0,      PI * 0.5, thickness, color); // 右下
    draw_arc_lines(x + r,     y + h - r, r, PI * 0.5, PI,       thickness, color); // 左下
}

fn draw_quarter_circle(cx: f32, cy: f32, r: f32, start: f32, end: f32, color: Color) {
    const SEGMENTS: usize = 8;
    let step = (end - start) / SEGMENTS as f32;
    for i in 0..SEGMENTS {
        let a0 = start + step * i as f32;
        let a1 = a0 + step;
        let p0 = vec2(cx + r * a0.cos(), cy + r * a0.sin());
        let p1 = vec2(cx + r * a1.cos(), cy + r * a1.sin());
        draw_triangle(vec2(cx, cy), p0, p1, color);
    }
}

fn draw_arc_lines(cx: f32, cy: f32, r: f32, start: f32, end: f32, thickness: f32, color: Color) {
    const SEGMENTS: usize = 8;
    let step = (end - start) / SEGMENTS as f32;
    let mut prev = vec2(cx + r * start.cos(), cy + r * start.sin());
    for i in 1..=SEGMENTS {
        let a = start + step * i as f32;
        let p = vec2(cx + r * a.cos(), cy + r * a.sin());
        draw_line(prev.x, prev.y, p.x, p.y, thickness, color);
        prev = p;
    }
}
