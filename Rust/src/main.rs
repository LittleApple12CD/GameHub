#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod common;
mod games;
mod ui;

use common::*;
use macroquad::prelude::*;

struct GameEntry {
    name: &'static str,
    make: fn() -> Box<dyn Game>,
}

fn game_registry() -> Vec<GameEntry> {
    vec![
        GameEntry { name: "SnakeGame",        make: || Box::new(games::snake::SnakeGame::new()) },
        GameEntry { name: "TetrisGame",       make: || Box::new(games::tetris::TetrisGame::new()) },
        GameEntry { name: "Minesweeper",  make: || Box::new(games::minesweeper::MinesweeperGame::new()) },
        GameEntry { name: "Flappy Bird",  make: || Box::new(games::flappy::FlappyGame::new()) },
        GameEntry { name: "TankBattle",         make: || Box::new(games::tank::TankGame::new()) },
        GameEntry { name: "Breakout",     make: || Box::new(games::breakout::BreakoutGame::new()) },
        GameEntry { name: "Tic Tac Toe",  make: || Box::new(games::tic_tac_toe::TicTacToeGame::new()) },
    ]
}

enum AppState {
    Menu,
    Playing(Box<dyn Game>),
}

#[macroquad::main(window_conf)]
async fn main() {
    common::init_font();
    
    let registry = game_registry();
    let mut state = AppState::Menu;
    let mut selected: usize = 0;
    let mut quit_selected = false;

    loop {
        match &mut state {
            AppState::Menu => {
                if let Some(action) = ui::menu_update(&registry, &mut selected, &mut quit_selected) {
                    match action {
                        ui::MenuAction::Play(idx) => {
                            let g = (registry[idx].make)();
                            state = AppState::Playing(g);
                        }
                        ui::MenuAction::Quit => break,
                    }
                }
                ui::draw_menu(&registry, selected, quit_selected);
            }
            AppState::Playing(game) => {
                let keep_going = game.update();
                game.draw();
                if !keep_going {
                    state = AppState::Menu;
                }
            }
        }
        next_frame().await;
    }
}

fn window_conf() -> Conf {
    Conf {
        window_title: "Game Hub".to_owned(),
        window_width: SCREEN_WIDTH as i32,
        window_height: SCREEN_HEIGHT as i32,
        window_resizable: false,
        ..Default::default()
    }
}
