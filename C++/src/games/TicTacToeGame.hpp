#pragma once
#include "../Common.hpp"

class TicTacToeGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    void reset();
    char checkWinner();
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);

    static constexpr int WIDTH = 500, HEIGHT = 500;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;
    int cell = WIDTH / 3;

    char board[3][3];
    char currentPlayer = 'X';
    char winner = 0;
    bool gameOver = false;
    int moveCount = 0;
    sf::Font font;
};