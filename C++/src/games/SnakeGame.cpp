#include "SnakeGame.hpp"
#include <random>
#include <iostream>

static std::mt19937 rng(std::random_device{}());

bool SnakeGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf"))
        std::cerr << "font missing\n";

    reset();
    sf::Clock clock;
    while (window.isOpen()) {
        float dt = clock.restart().asSeconds();
        if (!handleEvents(window)) return true;
        update(dt);
        draw(window);
    }
    return false;
}

void SnakeGame::reset() {
    snake.clear();
    snake.push_back({gridW / 2, gridH / 2});
    dir = {1, 0}; nextDir = {1, 0};
    score = 0; gameOver = false; moveTimer = 0.f;
    spawnFood();
}

void SnakeGame::spawnFood() {
    std::uniform_int_distribution<int> dx(0, gridW - 1), dy(0, gridH - 1);
    while (true) {
        sf::Vector2i p(dx(rng), dy(rng));
        bool onSnake = false;
        for (auto& s : snake) if (s == p) { onSnake = true; break; }
        if (!onSnake) { food = p; return; }
    }
}

bool SnakeGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            using K = sf::Keyboard::Key;
            if (key->code == K::Escape) return false;
            if (key->code == K::R) reset();
            if (key->code == K::Up && dir != sf::Vector2i(0, 1)) nextDir = {0, -1};
            else if (key->code == K::Down && dir != sf::Vector2i(0, -1)) nextDir = {0, 1};
            else if (key->code == K::Left && dir != sf::Vector2i(1, 0)) nextDir = {-1, 0};
            else if (key->code == K::Right && dir != sf::Vector2i(-1, 0)) nextDir = {1, 0};
        }
    }
    return true;
}

void SnakeGame::update(float dt) {
    if (gameOver) return;
    moveTimer += dt;
    if (moveTimer < moveDelay) return;
    moveTimer = 0.f;
    dir = nextDir;
    sf::Vector2i head = snake.front();
    sf::Vector2i nh = head + dir;

    if (nh.x < 0 || nh.x >= gridW || nh.y < 0 || nh.y >= gridH) { gameOver = true; return; }
    for (auto& s : snake) if (s == nh) { gameOver = true; return; }

    snake.push_front(nh);
    if (nh == food) { score += 10; spawnFood(); }
    else snake.pop_back();
}

void SnakeGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);

    sf::RectangleShape board(sf::Vector2f(WIDTH, HEIGHT));
    board.setPosition({(float)offsetX, (float)offsetY});
    board.setFillColor(sf::Color(20, 20, 30));
    window.draw(board);

    for (size_t i = 0; i < snake.size(); ++i) {
        sf::Color c = (i == 0) ? sf::Color(60, 220, 60) : sf::Color(40, 180, 40);
        sf::FloatRect r({(float)(offsetX + snake[i].x * CELL + 1),
                     (float)(offsetY + snake[i].y * CELL + 1)},
                    {(float)(CELL - 2), (float)(CELL - 2)});
        drawRect(window, r, c, 4.f);
    }

    sf::FloatRect fr({(float)(offsetX + food.x * CELL + 1),
                      (float)(offsetY + food.y * CELL + 1)},
                     {(float)(CELL - 2), (float)(CELL - 2)});
    drawRect(window, fr, sf::Color(255, 60, 60), 6.f);

    sf::Text st(font, "Score: " + std::to_string(score), 30);
    st.setFillColor(COLOR_WHITE);
    st.setPosition({(float)offsetX + 10, (float)offsetY + 10});
    window.draw(st);

    if (gameOver) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0, 0, 0, 180));
        window.draw(ov);
        sf::Text t(font, "GAME OVER - Press R", 30);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({offsetX + WIDTH / 2.f - t.getLocalBounds().size.x / 2.f,
                       offsetY + HEIGHT / 2.f});
        window.draw(t);
    }
    window.display();
}