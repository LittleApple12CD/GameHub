#include "Menu.hpp"
#include "games/SnakeGame.hpp"
#include "games/TetrisGame.hpp"
#include "games/MinesweeperGame.hpp"
#include "games/FlappyGame.hpp"
#include "games/TankGame.hpp"
#include "games/BreakoutGame.hpp"
#include "games/TicTacToeGame.hpp"
#include <iostream>

Menu::Menu() {
    if (!font.openFromFile("assets/fonts/arial.ttf"))
        std::cerr << "Warning: font not found\n";

    names = {"Snake", "Tetris", "Minesweeper", "Flappy Bird",
             "Tank", "Breakout", "Tic Tac Toe"};
    factories = {
        [] { return std::make_unique<SnakeGame>(); },
        [] { return std::make_unique<TetrisGame>(); },
        [] { return std::make_unique<MinesweeperGame>(); },
        [] { return std::make_unique<FlappyGame>(); },
        [] { return std::make_unique<TankGame>(); },
        [] { return std::make_unique<BreakoutGame>(); },
        [] { return std::make_unique<TicTacToeGame>(); },
    };
}

void Menu::drawMenu(sf::RenderWindow& window, int sel, bool quitSel) {
    window.clear(COLOR_BG);

    sf::Text title(font, "GAME HUB", 52);
    title.setFillColor(COLOR_WHITE);
    title.setPosition({SCREEN_WIDTH / 2.f - title.getLocalBounds().size.x / 2.f, 25.f});
    window.draw(title);

    for (size_t i = 0; i < names.size(); ++i) {
        float y = 130.f + i * 52.f;
        sf::FloatRect r({SCREEN_WIDTH / 2.f - 140.f, y}, {280.f, 46.f});
        sf::Color c = (static_cast<int>(i) == sel) ? COLOR_HOVER : COLOR_BUTTON;
        drawRect(window, r, c, 12.f, 2.f, COLOR_BORDER);

        sf::Text t(font, names[i], 30);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({SCREEN_WIDTH / 2.f - t.getLocalBounds().size.x / 2.f, y + 6.f});
        window.draw(t);
    }

    float qy = 130.f + names.size() * 52.f + 20.f;
    sf::FloatRect qr({SCREEN_WIDTH / 2.f - 60.f, qy}, {120.f, 40.f});
    drawRect(window, qr, quitSel ? COLOR_QUIT_HOVER : COLOR_QUIT, 10.f, 2.f, COLOR_BORDER);

    sf::Text qt(font, "EXIT", 30);
    qt.setFillColor(COLOR_WHITE);
    qt.setPosition({SCREEN_WIDTH / 2.f - qt.getLocalBounds().size.x / 2.f, qy + 5.f});
    window.draw(qt);

    window.display();
}

void Menu::run(sf::RenderWindow& window) {
    while (window.isOpen()) {
        while (const auto event = window.pollEvent()) {
            if (event->is<sf::Event::Closed>()) { window.close(); return; }

            if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
                using K = sf::Keyboard::Key;
                if (key->code == K::Escape) { window.close(); return; }
                else if (key->code == K::Up) {
                    if (quitSelected) { quitSelected = false; selected = (int)names.size() - 1; }
                    else selected = (selected - 1 + (int)names.size()) % (int)names.size();
                } else if (key->code == K::Down) {
                    if (selected == (int)names.size() - 1) quitSelected = true;
                    else { selected = (selected + 1) % (int)names.size(); quitSelected = false; }
                } else if (key->code == K::Enter) {
                    if (quitSelected) { window.close(); return; }
                    auto game = factories[selected]();
                    game->run(window);
                }
            }

            if (const auto* mm = event->getIf<sf::Event::MouseMoved>()) {
                sf::Vector2f p(mm->position.x, mm->position.y);
                float qy = 130.f + names.size() * 52.f + 20.f;
                sf::FloatRect qr({SCREEN_WIDTH / 2.f - 60.f, qy}, {120.f, 40.f});
                quitSelected = qr.contains(p);
                if (!quitSelected) {
                    for (size_t i = 0; i < names.size(); ++i) {
                        float y = 130.f + i * 52.f;
                        sf::FloatRect r({SCREEN_WIDTH / 2.f - 140.f, y}, {280.f, 46.f});
                        if (r.contains(p)) { selected = (int)i; break; }
                    }
                }
            }

            if (const auto* mb = event->getIf<sf::Event::MouseButtonPressed>()) {
                if (mb->button == sf::Mouse::Button::Left) {
                    sf::Vector2f p(mb->position.x, mb->position.y);
                    float qy = 130.f + names.size() * 52.f + 20.f;
                    sf::FloatRect qr({SCREEN_WIDTH / 2.f - 60.f, qy}, {120.f, 40.f});
                    if (qr.contains(p)) { window.close(); return; }
                    for (size_t i = 0; i < names.size(); ++i) {
                        float y = 130.f + i * 52.f;
                        sf::FloatRect r({SCREEN_WIDTH / 2.f - 140.f, y}, {280.f, 46.f});
                        if (r.contains(p)) {
                            auto game = factories[i]();
                            game->run(window);
                            break;
                        }
                    }
                }
            }
        }
        drawMenu(window, selected, quitSelected);
    }
}
