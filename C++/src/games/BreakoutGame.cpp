#include "BreakoutGame.hpp"
#include <random>
#include <iostream>
#include <algorithm>

static std::mt19937 rng(std::random_device{}());

bool BreakoutGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf"))
        std::cerr << "font missing\n";
    reset();
    while (window.isOpen()) {
        if (!handleEvents(window)) return true;
        update();
        draw(window);
        sf::sleep(sf::milliseconds(16));
    }
    return false;
}

void BreakoutGame::reset() {
    paddle = sf::FloatRect({WIDTH / 2.f - 60.f, HEIGHT - 40.f}, {120.f, 16.f});
    ballPos = {WIDTH / 2.f - 10.f, HEIGHT - 70.f};
    ballDir = {5.f, -6.f};
    ballSpeed = 6.f;
    bricks.clear();
    score = 0; lives = 3; gameOver = false; waiting = true;

    int rows = 5, cols = 8;
    float bw = (WIDTH - 20.f) / cols - 4.f;
    float bh = 22.f;
    sf::Color colors[] = {
        sf::Color(230,50,50), sf::Color(230,150,50), sf::Color(230,230,50),
        sf::Color(50,230,50), sf::Color(50,150,230)
    };
    for (int r = 0; r < rows; ++r)
        for (int c = 0; c < cols; ++c) {
            float x = 10.f + c * (bw + 4.f);
            float y = 40.f + r * (bh + 4.f);
            bricks.push_back({sf::FloatRect({x, y}, {bw, bh}), colors[r % 5], true});
        }
}

bool BreakoutGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            using K = sf::Keyboard::Key;
            if (key->code == K::Escape) return false;
            if (key->code == K::R) reset();
            if (key->code == K::Space && waiting && !gameOver) waiting = false;
        }
    }
    if (!gameOver && !waiting) {
        if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::Left) && paddle.position.x > 0)
            paddle.position.x -= paddleSpeed;
        if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::Right) &&
            paddle.position.x + paddle.size.x < WIDTH)
            paddle.position.x += paddleSpeed;
    }
    return true;
}

void BreakoutGame::update() {
    if (gameOver || waiting) return;

    ballPos += ballDir;
    sf::FloatRect ball(ballPos, {ballSize, ballSize});

    if (ball.position.x <= 0 || ball.position.x + ballSize >= WIDTH) ballDir.x = -ballDir.x;
    if (ball.position.y <= 0) ballDir.y = -ballDir.y;

    if (ball.position.y + ballSize >= HEIGHT) {
        lives--;
        if (lives <= 0) gameOver = true;
        else {
            waiting = true;
            ballPos = {WIDTH / 2.f - 10.f, HEIGHT - 70.f};
            std::uniform_int_distribution<int> d(0, 1);
            ballDir.x = d(rng) ? ballSpeed : -ballSpeed;
            ballDir.y = -ballSpeed;
        }
        return;
    }

    if (ball.findIntersection(paddle)) {
        ballDir.y = -std::abs(ballDir.y);
        float hit = (ball.position.x + ballSize / 2 - (paddle.position.x + paddle.size.x / 2))
                    / (paddle.size.x / 2);
        ballDir.x = hit * ballSpeed * 0.9f;
        if (std::abs(ballDir.x) < 1.5f) ballDir.x = (ballDir.x >= 0) ? 3.f : -3.f;
    }

    for (auto& b : bricks) {
        if (!b.alive) continue;
        if (ball.findIntersection(b.rect)) {
            b.alive = false;
            score += 10;
            float ot = (ball.position.y + ballSize) - b.rect.position.y;
            float ob = (b.rect.position.y + b.rect.size.y) - ball.position.y;
            float ol = (ball.position.x + ballSize) - b.rect.position.x;
            float orr = (b.rect.position.x + b.rect.size.x) - ball.position.x;
            float mo = std::min({ot, ob, ol, orr});
            if (mo == ot || mo == ob) ballDir.y = -ballDir.y;
            else ballDir.x = -ballDir.x;
            break;
        }
    }

    bool allDead = true;
    for (auto& b : bricks) if (b.alive) { allDead = false; break; }
    if (allDead) gameOver = true;
}

void BreakoutGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape board(sf::Vector2f(WIDTH, HEIGHT));
    board.setPosition({(float)offsetX, (float)offsetY});
    board.setFillColor(sf::Color(20, 20, 30));
    window.draw(board);

    auto off = [&](sf::FloatRect r) {
        return sf::FloatRect({r.position.x + offsetX, r.position.y + offsetY}, r.size);
    };

    for (auto& b : bricks) {
        if (!b.alive) continue;
        drawRect(window, off(b.rect), b.color, 8.f, 1.f, sf::Color::White);
    }
    drawRect(window, off(paddle), sf::Color::White, 10.f, 2.f, sf::Color(200,200,220));

    sf::CircleShape ballC(ballSize / 2.f);
    ballC.setPosition({(float)offsetX + ballPos.x, (float)offsetY + ballPos.y});
    ballC.setFillColor(sf::Color(255, 255, 120));
    window.draw(ballC);

    sf::Text st(font, "Score: " + std::to_string(score), 30);
    st.setFillColor(COLOR_WHITE);
    st.setPosition({(float)offsetX + 10, (float)offsetY + 10});
    window.draw(st);

    sf::Text lt(font, "Lives: " + std::to_string(lives), 30);
    lt.setFillColor(COLOR_WHITE);
    lt.setPosition({(float)offsetX + WIDTH - 120, (float)offsetY + 10});
    window.draw(lt);

    if (waiting && !gameOver) {
        sf::Text t(font, "Press SPACE", 48);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({offsetX + WIDTH / 2.f - t.getLocalBounds().size.x / 2.f,
                       offsetY + HEIGHT / 2.f});
        window.draw(t);
    } else if (gameOver) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0, 0, 0, 150));
        window.draw(ov);
        std::string msg = (lives <= 0) ? "GAME OVER" : "YOU WIN!";
        sf::Text t(font, msg, 48);
        t.setFillColor(lives <= 0 ? sf::Color::White : sf::Color(0, 255, 0));
        t.setPosition({offsetX + WIDTH / 2.f - t.getLocalBounds().size.x / 2.f,
                       offsetY + HEIGHT / 2.f - 30.f});
        window.draw(t);
        sf::Text t2(font, "Press R to restart", 30);
        t2.setFillColor(COLOR_WHITE);
        t2.setPosition({offsetX + WIDTH / 2.f - t2.getLocalBounds().size.x / 2.f,
                        offsetY + HEIGHT / 2.f + 30.f});
        window.draw(t2);
    }
    window.display();
}