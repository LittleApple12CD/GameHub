#pragma once
#include "../Common.hpp"
#include <vector>

class FlappyGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    struct Pipe { float x; int gapY; bool passed; };
    void reset();
    void addPipe();
    void update();
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);
    float birdX() const { return 80.f; }

    static constexpr int WIDTH = 400, HEIGHT = 600;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;

    float birdY = 0.f, birdVy = 0.f, birdRadius = 15.f;
    float gravity = 0.6f, jump = -11.f;
    std::vector<Pipe> pipes;
    float pipeW = 60.f; int pipeGap = 160; float pipeSpeed = 4.f;
    int score = 0; bool gameOver = false;
    int pipeTimer = 0, pipeDelay = 90;
    sf::Font font;
};
