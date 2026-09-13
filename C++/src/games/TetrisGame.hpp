#pragma once
#include "../Common.hpp"
#include <map>
#include <vector>

class TetrisGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    void reset();
    void spawnPiece();
    bool checkCollision(const std::vector<std::vector<int>>& s, int x, int y);
    void rotatePiece();
    void lockPiece();
    void clearLines();
    void update(float dt);
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);

    static constexpr int COLS = 10, ROWS = 20, CELL = 30;
    int WIDTH = COLS * CELL, HEIGHT = ROWS * CELL;
    int offsetX = (SCREEN_WIDTH - WIDTH - 140) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;

    std::map<char, std::vector<std::vector<int>>> shapes;
    std::map<char, sf::Color> colors;
    std::vector<std::vector<char>> board;
    char currentPiece;
    std::vector<std::vector<int>> pieceShape;
    int pieceX = 0, pieceY = 0;
    int score = 0;
    bool gameOver = false;
    float fallTimer = 0.f, fallDelay = 0.5f;
    sf::Font font;
};