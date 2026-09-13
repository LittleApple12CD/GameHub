use crate::common::*;
use macroquad::prelude::*;
use crate::GameEntry;

pub enum MenuAction {
    Play(usize),
    Quit,
}

const BTN_W: f32 = 280.0;
const BTN_H: f32 = 46.0;
const BTN_SPACING: f32 = 52.0;
const BTN_START_Y: f32 = 130.0;

fn button_rect(i: usize) -> Rect {
    Rect::new(
        SCREEN_WIDTH / 2.0 - BTN_W / 2.0,
        BTN_START_Y + i as f32 * BTN_SPACING,
        BTN_W,
        BTN_H,
    )
}

fn quit_rect(n: usize) -> Rect {
    Rect::new(
        SCREEN_WIDTH / 2.0 - 60.0,
        BTN_START_Y + n as f32 * BTN_SPACING + 20.0,
        120.0,
        40.0,
    )
}

pub fn menu_update(
    registry: &[GameEntry],
    selected: &mut usize,
    quit_selected: &mut bool,
) -> Option<MenuAction> {
    let (mx, my) = mouse_position();

    if is_key_pressed(KeyCode::Escape) {
        return Some(MenuAction::Quit);
    }
    if is_key_pressed(KeyCode::Up) {
        if *quit_selected {
            *quit_selected = false;
            *selected = registry.len() - 1;
        } else {
            *selected = if *selected == 0 { registry.len() - 1 } else { *selected - 1 };
        }
    }
    if is_key_pressed(KeyCode::Down) {
        if *selected == registry.len() - 1 {
            *quit_selected = true;
        } else {
            *selected += 1;
            *quit_selected = false;
        }
    }
    if is_key_pressed(KeyCode::Enter) {
        if *quit_selected {
            return Some(MenuAction::Quit);
        }
        return Some(MenuAction::Play(*selected));
    }

    *quit_selected = quit_rect(registry.len()).contains(vec2(mx, my));
    if !*quit_selected {
        for i in 0..registry.len() {
            if button_rect(i).contains(vec2(mx, my)) {
                *selected = i;
                break;
            }
        }
    }

    if is_mouse_button_pressed(MouseButton::Left) {
        if quit_rect(registry.len()).contains(vec2(mx, my)) {
            return Some(MenuAction::Quit);
        }
        for i in 0..registry.len() {
            if button_rect(i).contains(vec2(mx, my)) {
                return Some(MenuAction::Play(i));
            }
        }
    }

    None
}

pub fn draw_menu(registry: &[GameEntry], selected: usize, quit_selected: bool) {
    clear_background(COLOR_BG);

    draw_text_centered("GAME HUB", SCREEN_WIDTH / 2.0, 55.0, 52, COLOR_WHITE);

    for (i, entry) in registry.iter().enumerate() {
        let r = button_rect(i);
        let color = if i == selected { COLOR_HOVER } else { COLOR_BUTTON };
        draw_rounded_rect(r.x, r.y, r.w, r.h, 12.0, color);
        draw_rounded_rect_lines(r.x, r.y, r.w, r.h, 12.0, 2.0, COLOR_BORDER);
        draw_text_centered(
            entry.name,
            SCREEN_WIDTH / 2.0,
            r.y + r.h / 2.0,
            38,
            COLOR_WHITE,
        );
    }

    let qr = quit_rect(registry.len());
    let qc = if quit_selected { COLOR_QUIT_HOVER } else { COLOR_QUIT };
    draw_rounded_rect(qr.x, qr.y, qr.w, qr.h, 10.0, qc);
    draw_rounded_rect_lines(qr.x, qr.y, qr.w, qr.h, 10.0, 2.0, COLOR_BORDER);
    draw_text_centered("EXIT", SCREEN_WIDTH / 2.0, qr.y + qr.h / 2.0, 38, COLOR_WHITE);
}