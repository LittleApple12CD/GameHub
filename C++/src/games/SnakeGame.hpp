#pragma once
#include "../Common.hpp"
#include <vector>
#include <deque>

class SnakeGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    void reset();
    void spawnFood();
    void update(float dt);
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);

    static constexpr int WIDTH = 600, HEIGHT = 600;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;
    static constexpr int CELL = 20;
    int gridW = WIDTH / CELL, gridH = HEIGHT / CELL;

    std::deque<sf::Vector2i> snake;
    sf::Vector2i dir{1, 0}, nextDir{1, 0};
    sf::Vector2i food;
    int score = 0;
    bool gameOver = false;
    float moveTimer = 0.f, moveDelay = 0.15f;

    sf::Font font;
};