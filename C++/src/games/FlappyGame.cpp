#include "FlappyGame.hpp"
#include <random>
#include <iostream>

static std::mt19937 rng(std::random_device{}());

bool FlappyGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf")) std::cerr << "font missing\n";
    reset();
    while (window.isOpen()) {
        if (!handleEvents(window)) return true;
        update();
        draw(window);
        sf::sleep(sf::milliseconds(16));
    }
    return false;
}

void FlappyGame::reset() {
    birdY = HEIGHT / 2.f; birdVy = 0.f;
    pipes.clear(); score = 0; gameOver = false; pipeTimer = 0;
}

void FlappyGame::addPipe() {
    std::uniform_int_distribution<int> d(80, HEIGHT - 80 - pipeGap);
    pipes.push_back({(float)WIDTH, d(rng), false});
}

bool FlappyGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            using K = sf::Keyboard::Key;
            if (key->code == K::Escape) return false;
            if (key->code == K::R) reset();
            if (key->code == K::Space && !gameOver) birdVy = jump;
        }
        if (const auto* mb = event->getIf<sf::Event::MouseButtonPressed>()) {
            if (mb->button == sf::Mouse::Button::Left && !gameOver) birdVy = jump;
        }
    }
    return true;
}

void FlappyGame::update() {
    if (gameOver) return;
    birdVy += gravity;
    birdY += birdVy;

    if (birdY - birdRadius < 0) { birdY = birdRadius; birdVy = 0; }
    else if (birdY + birdRadius > HEIGHT) { gameOver = true; return; }

    if (++pipeTimer >= pipeDelay) { pipeTimer = 0; addPipe(); }

    for (auto it = pipes.begin(); it != pipes.end();) {
        it->x -= pipeSpeed;
        if (!it->passed && it->x + pipeW < birdX()) { it->passed = true; score++; }

        if (it->x < birdX() + birdRadius && it->x + pipeW > birdX() - birdRadius) {
            if (birdY - birdRadius < it->gapY || birdY + birdRadius > it->gapY + pipeGap) {
                gameOver = true; return;
            }
        }
        if (it->x + pipeW < 0) it = pipes.erase(it);
        else ++it;
    }
}

void FlappyGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape sky(sf::Vector2f(WIDTH, HEIGHT));
    sky.setPosition({(float)offsetX, (float)offsetY});
    sky.setFillColor(sf::Color(135, 206, 235));
    window.draw(sky);

    auto ox = [&](float x){ return x + offsetX; };
    auto oy = [&](float y){ return y + offsetY; };

    sf::CircleShape bird(birdRadius);
    bird.setPosition({ox(birdX()) - birdRadius, oy(birdY) - birdRadius});
    bird.setFillColor(sf::Color(255, 255, 50));
    window.draw(bird);

    sf::CircleShape eye(4);
    eye.setPosition({ox(birdX()) + 8 - 4, oy(birdY) - 5 - 4});
    eye.setFillColor(sf::Color::Black);
    window.draw(eye);

    for (auto& p : pipes) {
        drawRect(window, sf::FloatRect({ox(p.x), (float)offsetY}, {(float)pipeW, (float)p.gapY}),
                 sf::Color(40, 200, 40), 6.f);
        drawRect(window, sf::FloatRect({ox(p.x) - 6, oy((float)p.gapY) - 22}, {(float)pipeW + 12, 22}),
                 sf::Color(30, 160, 30), 6.f);
        float by = (float)(p.gapY + pipeGap);
        drawRect(window, sf::FloatRect({ox(p.x), oy(by)}, {(float)pipeW, HEIGHT - by}),
                 sf::Color(40, 200, 40), 6.f);
        drawRect(window, sf::FloatRect({ox(p.x) - 6, oy(by)}, {(float)pipeW + 12, 22}),
                 sf::Color(30, 160, 30), 6.f);
    }

    sf::Text st(font, std::to_string(score), 40);
    st.setFillColor(COLOR_WHITE);
    st.setPosition({ox(WIDTH / 2.f) - st.getLocalBounds().size.x / 2.f, oy(50.f)});
    window.draw(st);

    if (gameOver) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0, 0, 0, 150));
        window.draw(ov);
        sf::Text t(font, "GAME OVER", 40);
        t.setFillColor(COLOR_WHITE);
        t.setPosition({ox(WIDTH / 2.f) - t.getLocalBounds().size.x / 2.f, oy(HEIGHT / 2.f - 30.f)});
        window.draw(t);
        sf::Text t2(font, "Press R to restart", 26);
        t2.setFillColor(COLOR_WHITE);
        t2.setPosition({ox(WIDTH / 2.f) - t2.getLocalBounds().size.x / 2.f, oy(HEIGHT / 2.f + 20.f)});
        window.draw(t2);
    }
    window.display();
}