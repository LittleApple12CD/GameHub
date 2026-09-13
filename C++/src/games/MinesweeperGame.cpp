#include "MinesweeperGame.hpp"
#include <random>
#include <iostream>

static std::mt19937 rng(std::random_device{}());

bool MinesweeperGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf")) std::cerr << "font missing\n";
    reset();
    while (window.isOpen()) {
        if (!handleEvents(window)) return true;
        draw(window);
        sf::sleep(sf::milliseconds(16));
    }
    return false;
}

void MinesweeperGame::reset() {
    board.assign(ROWS, std::vector<int>(COLS, 0));
    revealed.assign(ROWS, std::vector<bool>(COLS, false));
    flagged.assign(ROWS, std::vector<bool>(COLS, false));
    gameOver = false; won = false; firstClick = true;
}

void MinesweeperGame::placeMines(int sx, int sy) {
    std::uniform_int_distribution<int> dx(0, COLS - 1), dy(0, ROWS - 1);
    int placed = 0;
    while (placed < MINES) {
        int x = dx(rng), y = dy(rng);
        if (board[y][x] == -1) continue;
        if (std::abs(x - sx) <= 1 && std::abs(y - sy) <= 1) continue;
        board[y][x] = -1; placed++;
    }
    for (int y = 0; y < ROWS; ++y)
        for (int x = 0; x < COLS; ++x) {
            if (board[y][x] == -1) continue;
            int c = 0;
            for (int dy2 = -1; dy2 <= 1; ++dy2)
                for (int dx2 = -1; dx2 <= 1; ++dx2) {
                    int nx = x + dx2, ny = y + dy2;
                    if (nx >= 0 && nx < COLS && ny >= 0 && ny < ROWS && board[ny][nx] == -1) c++;
                }
            board[y][x] = c;
        }
}

void MinesweeperGame::reveal(int x, int y) {
    if (x < 0 || x >= COLS || y < 0 || y >= ROWS) return;
    if (revealed[y][x] || flagged[y][x]) return;
    revealed[y][x] = true;
    if (board[y][x] == -1) { gameOver = true; return; }
    if (board[y][x] == 0)
        for (int dy = -1; dy <= 1; ++dy)
            for (int dx = -1; dx <= 1; ++dx)
                reveal(x + dx, y + dy);

    int cnt = 0;
    for (int yy = 0; yy < ROWS; ++yy)
        for (int xx = 0; xx < COLS; ++xx)
            if (revealed[yy][xx]) cnt++;
    if (cnt == ROWS * COLS - MINES) won = true;
}

bool MinesweeperGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            if (key->code == sf::Keyboard::Key::Escape) return false;
            if (key->code == sf::Keyboard::Key::R) reset();
        }
        if (const auto* mb = event->getIf<sf::Event::MouseButtonPressed>()) {
            if (gameOver || won) continue;
            int x = (mb->position.x - offsetX) / CELL;
            int y = (mb->position.y - offsetY) / CELL;
            if (x < 0 || x >= COLS || y < 0 || y >= ROWS) continue;
            if (mb->button == sf::Mouse::Button::Left) {
                if (firstClick) { placeMines(x, y); firstClick = false; }
                if (!flagged[y][x]) reveal(x, y);
            } else if (mb->button == sf::Mouse::Button::Right) {
                if (!revealed[y][x]) flagged[y][x] = !flagged[y][x];
            }
        }
    }
    return true;
}

void MinesweeperGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape bg(sf::Vector2f(WIDTH, HEIGHT));
    bg.setPosition({(float)offsetX, (float)offsetY});
    bg.setFillColor(sf::Color(30, 30, 42));
    window.draw(bg);

    sf::Color numCols[] = {
        sf::Color::Black, sf::Color::Blue, sf::Color(0,150,0), sf::Color::Red,
        sf::Color(0,0,180), sf::Color(150,0,0), sf::Color(0,150,150), sf::Color::Black
    };

    for (int y = 0; y < ROWS; ++y)
        for (int x = 0; x < COLS; ++x) {
            sf::FloatRect r({(float)(offsetX + x * CELL + 1), (float)(offsetY + y * CELL + 1)},
                            {CELL - 2.f, CELL - 2.f});
            if (revealed[y][x]) {
                if (board[y][x] == -1) {
                    drawRect(window, r, sf::Color(200,50,50), 4.f);
                    sf::CircleShape c(9);
                    c.setPosition({r.position.x + CELL/2.f - 9, r.position.y + CELL/2.f - 9});
                    c.setFillColor(sf::Color::Black);
                    window.draw(c);
                } else {
                    drawRect(window, r, sf::Color(190,190,200), 4.f);
                    if (board[y][x] > 0) {
                        sf::Text t(font, std::to_string(board[y][x]), 20);
                        t.setFillColor(numCols[board[y][x]]);
                        t.setPosition({r.position.x + 10, r.position.y + 6});
                        window.draw(t);
                    }
                }
            } else {
                sf::Color c = flagged[y][x] ? sf::Color(210,210,60) : sf::Color(80,80,105);
                drawRect(window, r, c, 4.f);
                if (flagged[y][x]) {
                    sf::Text t(font, "F", 20);
                    t.setFillColor(sf::Color(255,50,50));
                    t.setPosition({r.position.x + 10, r.position.y + 5});
                    window.draw(t);
                }
            }
        }

    if (gameOver || won) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0, 0, 0, 150));
        window.draw(ov);
        sf::Text t(font, won ? "YOU WIN!" : "GAME OVER", 40);
        t.setFillColor(won ? sf::Color(0,255,0) : sf::Color::White);
        t.setPosition({offsetX + WIDTH/2.f - t.getLocalBounds().size.x/2.f, offsetY + HEIGHT/2.f - 30.f});
        window.draw(t);
        sf::Text t2(font, "Press R to restart", 22);
        t2.setFillColor(COLOR_WHITE);
        t2.setPosition({offsetX + WIDTH/2.f - t2.getLocalBounds().size.x/2.f, offsetY + HEIGHT/2.f + 30.f});
        window.draw(t2);
    }
    window.display();
}