#include "TankGame.hpp"
#include <random>
#include <iostream>
#include <algorithm>

static std::mt19937 rng(std::random_device{}());

bool TankGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf")) std::cerr << "font missing\n";
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

std::set<std::pair<int,int>> TankGame::generateWalls() {
    std::set<std::pair<int,int>> w;
    for (int i = 0; i < GRID; ++i) {
        w.insert({i, 0}); w.insert({i, GRID - 1});
        w.insert({0, i}); w.insert({GRID - 1, i});
    }
    std::uniform_int_distribution<int> d(2, GRID - 3);
    for (int k = 0; k < 12; ++k) {
        int x = d(rng), y = d(rng);
        if ((x == 1 && y == 1) || (x == GRID - 2 && y == GRID - 2)) continue;
        w.insert({x, y});
    }
    return w;
}

void TankGame::reset() {
    walls = generateWalls();
    player = {1, 1, {1, 0}, true, 0.f, 0.f};
    enemies.clear();
    enemies.push_back({GRID - 2, GRID - 2, {-1, 0}, true, 0.f, 0.f});
    bullets.clear();
    score = 0; gameOver = false; enemySpawnTimer = 0.f;
}

bool TankGame::canMove(int x, int y, bool isPlayer) {
    if (x < 0 || x >= GRID || y < 0 || y >= GRID) return false;
    if (walls.count({x, y})) return false;
    if (isPlayer) {
        for (auto& e : enemies)
            if (e.alive && e.x == x && e.y == y) return false;
    } else {
        if (player.alive && player.x == x && player.y == y) return false;
        for (auto& e : enemies)
            if (e.alive && e.x == x && e.y == y && !(e.x == player.x && e.y == player.y)) {}
    }
    return true;
}

void TankGame::shoot(int x, int y, sf::Vector2i dir, bool isPlayer) {
    bullets.push_back({(float)(x + dir.x), (float)(y + dir.y), dir, isPlayer, true});
}

bool TankGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            using K = sf::Keyboard::Key;
            if (key->code == K::Escape) return false;
            if (key->code == K::R) reset();
            if (key->code == K::Space && !gameOver && player.alive)
                shoot(player.x, player.y, player.dir, true);
        }
    }
    return true;
}

void TankGame::update(float dt) {
    if (gameOver) return;

    if (player.alive) {
        sf::Vector2i nd = player.dir;
        bool moved = false;
        if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::W)) { nd = {0, -1}; moved = true; }
        else if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::S)) { nd = {0, 1}; moved = true; }
        else if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::A)) { nd = {-1, 0}; moved = true; }
        else if (sf::Keyboard::isKeyPressed(sf::Keyboard::Key::D)) { nd = {1, 0}; moved = true; }

        if (moved) {
            player.dir = nd;
            player.moveTimer += dt;
            if (player.moveTimer >= playerMoveDelay) {
                player.moveTimer = 0.f;
                int nx = player.x + nd.x, ny = player.y + nd.y;
                if (canMove(nx, ny, true)) { player.x = nx; player.y = ny; }
            }
        } else player.moveTimer = 0.f;
    }

    // 子弹
    for (auto& b : bullets) {
        if (!b.alive) continue;
        b.x += b.dir.x * bulletSpeed * dt * 6.f;
        b.y += b.dir.y * bulletSpeed * dt * 6.f;
        int bx = (int)std::round(b.x), by = (int)std::round(b.y);
        if (bx < 0 || bx >= GRID || by < 0 || by >= GRID || walls.count({bx, by})) {
            b.alive = false; continue;
        }
        if (b.isPlayer) {
            for (auto& e : enemies)
                if (e.alive && e.x == bx && e.y == by) {
                    e.alive = false; b.alive = false; score += 10; break;
                }
        } else {
            if (player.alive && player.x == bx && player.y == by) {
                player.alive = false; b.alive = false; gameOver = true;
            }
        }
    }
    bullets.erase(std::remove_if(bullets.begin(), bullets.end(),
                  [](const Bullet& b){ return !b.alive; }), bullets.end());

    // 敌人
    std::uniform_real_distribution<float> rd(0.f, 1.f);
    for (auto& e : enemies) {
        if (!e.alive) continue;
        e.timer += dt;
        if (e.timer >= 0.33f) {
            e.timer = 0.f;
            if (rd(rng) < 0.25f) {
                sf::Vector2i dirs[] = {{1,0},{-1,0},{0,1},{0,-1}};
                e.dir = dirs[(int)(rd(rng) * 4) % 4];
            }
            e.moveTimer += 0.33f;
            if (e.moveTimer >= enemyMoveDelay) {
                e.moveTimer = 0.f;
                int nx = e.x + e.dir.x, ny = e.y + e.dir.y;
                if (canMove(nx, ny, false)) { e.x = nx; e.y = ny; }
            }
        }
        if (rd(rng) < 0.015f) shoot(e.x, e.y, e.dir, false);
    }

    enemySpawnTimer += dt;
    if (enemySpawnTimer >= enemySpawnDelay) {
        enemySpawnTimer = 0.f;
        int alive = 0;
        for (auto& e : enemies) if (e.alive) alive++;
        if (alive < 3)
            enemies.push_back({GRID - 2, GRID - 2, {-1, 0}, true, 0.f, 0.f});
    }
}

void TankGame::drawTank(sf::RenderTarget& t, int x, int y, sf::Color body, sf::Color gun, sf::Vector2i dir) {
    float cx = x * CELL + 15.f, cy = y * CELL + 15.f;
    sf::FloatRect bodyRect({(float)(x * CELL + 2), (float)(y * CELL + 2)}, {26.f, 26.f});
    drawRect(t, bodyRect, body, 4.f);

    sf::RectangleShape g;
    if (dir.x == 1) g.setSize({14, 8});
    else if (dir.x == -1) g.setSize({14, 8});
    else if (dir.y == -1) g.setSize({8, 14});
    else g.setSize({8, 14});
    g.setFillColor(gun);
    g.setPosition({cx - g.getSize().x/2.f + dir.x * 6.f, cy - g.getSize().y/2.f + dir.y * 6.f});
    t.draw(g);

    sf::CircleShape c(6);
    c.setPosition({cx - 6, cy - 6});
    c.setFillColor(body);
    t.draw(c);
}

void TankGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape bg(sf::Vector2f(WIDTH, HEIGHT));
    bg.setPosition({(float)offsetX, (float)offsetY});
    bg.setFillColor(sf::Color(30, 30, 42));
    window.draw(bg);

    sf::View oldView = window.getView();
    sf::View v(sf::FloatRect({(float)-offsetX, (float)-offsetY}, {(float)SCREEN_WIDTH, (float)SCREEN_HEIGHT}));
    window.setView(v);

    for (auto& w : walls) {
        sf::FloatRect r({(float)(w.first * CELL + 1), (float)(w.second * CELL + 1)},
                        {(float)(CELL - 2), (float)(CELL - 2)});
        drawRect(window, r, sf::Color(100, 100, 130), 4.f);
    }

    if (player.alive)
        drawTank(window, player.x, player.y, sf::Color(50, 230, 50), sf::Color(30, 200, 30), player.dir);
    for (auto& e : enemies)
        if (e.alive)
            drawTank(window, e.x, e.y, sf::Color(230, 50, 50), sf::Color(200, 30, 30), e.dir);

    for (auto& b : bullets) {
        sf::RectangleShape r(sf::Vector2f(10, 10));
        r.setPosition({b.x * CELL + 10, b.y * CELL + 10});
        r.setFillColor(b.isPlayer ? sf::Color(255, 255, 80) : sf::Color(255, 150, 50));
        window.draw(r);
    }

    window.setView(oldView);

    sf::Text st(font, "Score: " + std::to_string(score), 30);
    st.setFillColor(COLOR_WHITE);
    st.setPosition({(float)offsetX + 10, (float)offsetY + 10});
    window.draw(st);

    if (gameOver) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0,0,0,150));
        window.draw(ov);
        sf::Text t(font, "GAME OVER", 44);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({offsetX + WIDTH/2.f - t.getLocalBounds().size.x/2.f, offsetY + HEIGHT/2.f - 30.f});
        window.draw(t);
        sf::Text t2(font, "Press R to restart", 26);
        t2.setFillColor(COLOR_WHITE);
        t2.setPosition({offsetX + WIDTH/2.f - t2.getLocalBounds().size.x/2.f, offsetY + HEIGHT/2.f + 30.f});
        window.draw(t2);
    }
    window.display();
}