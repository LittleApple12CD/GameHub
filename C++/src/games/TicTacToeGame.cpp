#include "TicTacToeGame.hpp"
#include <iostream>

bool TicTacToeGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf")) std::cerr << "font missing\n";
    reset();
    while (window.isOpen()) {
        if (!handleEvents(window)) return true;
        draw(window);
        sf::sleep(sf::milliseconds(16));
    }
    return false;
}

void TicTacToeGame::reset() {
    for (int y = 0; y < 3; ++y)
        for (int x = 0; x < 3; ++x) board[y][x] = 0;
    currentPlayer = 'X';
    winner = 0; gameOver = false; moveCount = 0;
}

char TicTacToeGame::checkWinner() {
    for (int r = 0; r < 3; ++r)
        if (board[r][0] && board[r][0] == board[r][1] && board[r][1] == board[r][2])
            return board[r][0];
    for (int c = 0; c < 3; ++c)
        if (board[0][c] && board[0][c] == board[1][c] && board[1][c] == board[2][c])
            return board[0][c];
    if (board[0][0] && board[0][0] == board[1][1] && board[1][1] == board[2][2]) return board[0][0];
    if (board[0][2] && board[0][2] == board[1][1] && board[1][1] == board[2][0]) return board[0][2];
    return 0;
}

bool TicTacToeGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            if (key->code == sf::Keyboard::Key::Escape) return false;
            if (key->code == sf::Keyboard::Key::R) reset();
        }
        if (const auto* mb = event->getIf<sf::Event::MouseButtonPressed>()) {
            if (mb->button != sf::Mouse::Button::Left || gameOver) continue;
            int x = (mb->position.x - offsetX) / cell;
            int y = (mb->position.y - offsetY) / cell;
            if (x < 0 || x >= 3 || y < 0 || y >= 3) continue;
            if (board[y][x]) continue;
            board[y][x] = currentPlayer;
            moveCount++;
            winner = checkWinner();
            if (winner) gameOver = true;
            else if (moveCount == 9) gameOver = true;
            else currentPlayer = (currentPlayer == 'X') ? 'O' : 'X';
        }
    }
    return true;
}

void TicTacToeGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape bg(sf::Vector2f(WIDTH, HEIGHT));
    bg.setPosition({(float)offsetX, (float)offsetY});
    bg.setFillColor(sf::Color(30, 30, 42));
    window.draw(bg);

    for (int r = 0; r < 3; ++r)
        for (int c = 0; c < 3; ++c) {
            sf::FloatRect rr({(float)(offsetX + c * cell + 2), (float)(offsetY + r * cell + 2)},
                             {cell - 4.f, cell - 4.f});
            drawRect(window, rr, sf::Color(50, 50, 65), 8.f);
            if (board[r][c]) {
                sf::Color col = (board[r][c] == 'X') ? sf::Color(60,200,255) : sf::Color(255,200,60);
                sf::Text t(font, std::string(1, board[r][c]), 80);
                t.setFillColor(col);
                t.setPosition({offsetX + c * cell + cell / 2.f - t.getLocalBounds().size.x / 2.f,
                               offsetY + r * cell + cell / 2.f - t.getLocalBounds().size.y / 2.f - 10});
                window.draw(t);
            }
        }

    for (int i = 1; i < 3; ++i) {
        sf::RectangleShape v(sf::Vector2f(3, HEIGHT - 10));
        v.setPosition({(float)(offsetX + i * cell - 1.5f), (float)(offsetY + 5)});
        v.setFillColor(sf::Color(80, 80, 100));
        window.draw(v);
        sf::RectangleShape h(sf::Vector2f(WIDTH - 10, 3));
        h.setPosition({(float)(offsetX + 5), (float)(offsetY + i * cell - 1.5f)});
        h.setFillColor(sf::Color(80, 80, 100));
        window.draw(h);
    }

    float statusY = offsetY + HEIGHT + 15.f;
    if (gameOver) {
        std::string msg;
        sf::Color col;
        if (winner) { msg = std::string("Player ") + winner + " Wins!"; col = sf::Color(0,255,0); }
        else { msg = "Draw!"; col = sf::Color(255,255,100); }
        sf::Text t(font, msg, 44);
        t.setFillColor(col);
        t.setPosition({SCREEN_WIDTH / 2.f - t.getLocalBounds().size.x / 2.f, statusY});
        window.draw(t);
        sf::Text t2(font, "Press R to restart", 24);
        t2.setFillColor(sf::Color(200,200,200));
        t2.setPosition({SCREEN_WIDTH / 2.f - t2.getLocalBounds().size.x / 2.f, statusY + 45.f});
        window.draw(t2);
    } else {
        std::string msg = std::string("Player ") + currentPlayer + "'s Turn";
        sf::Text t(font, msg, 44);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({SCREEN_WIDTH / 2.f - t.getLocalBounds().size.x / 2.f, statusY});
        window.draw(t);
    }
    window.display();
}