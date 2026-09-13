#pragma once
#include "../Common.hpp"
#include <vector>

class BreakoutGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    struct Brick { sf::FloatRect rect; sf::Color color; bool alive; };
    void reset();
    void update();
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);

    static constexpr int WIDTH = 600, HEIGHT = 600;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;

    sf::FloatRect paddle;
    sf::Vector2f ballPos, ballDir;
    float ballSize = 20.f, ballSpeed = 6.f, paddleSpeed = 8.f;
    std::vector<Brick> bricks;
    int score = 0, lives = 3;
    bool gameOver = false, waiting = true;
    sf::Font font;
};