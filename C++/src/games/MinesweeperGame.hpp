#pragma once
#include "../Common.hpp"
#include <vector>

class MinesweeperGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    void reset();
    void placeMines(int sx, int sy);
    void reveal(int x, int y);
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);

    static constexpr int ROWS = 16, COLS = 16, MINES = 40, CELL = 30;
    int WIDTH = COLS * CELL, HEIGHT = ROWS * CELL;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;

    std::vector<std::vector<int>> board;
    std::vector<std::vector<bool>> revealed, flagged;
    bool gameOver = false, won = false, firstClick = true;
    sf::Font font;
};