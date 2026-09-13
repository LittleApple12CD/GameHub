#pragma once
#include "../Common.hpp"
#include <vector>
#include <set>

class TankGame : public Game {
public:
    bool run(sf::RenderWindow& window) override;
private:
    struct Tank { int x, y; sf::Vector2i dir; bool alive; float timer, moveTimer; };
    struct Bullet { float x, y; sf::Vector2i dir; bool isPlayer, alive; };

    void reset();
    std::set<std::pair<int,int>> generateWalls();
    bool canMove(int x, int y, bool isPlayer);
    void shoot(int x, int y, sf::Vector2i dir, bool isPlayer);
    void update(float dt);
    void draw(sf::RenderWindow& window);
    bool handleEvents(sf::RenderWindow& window);
    void drawTank(sf::RenderTarget& t, int x, int y, sf::Color body, sf::Color gun, sf::Vector2i dir);

    static constexpr int WIDTH = 600, HEIGHT = 600, CELL = 30;
    int GRID = WIDTH / CELL;
    int offsetX = (SCREEN_WIDTH - WIDTH) / 2;
    int offsetY = (SCREEN_HEIGHT - HEIGHT) / 2;

    std::set<std::pair<int,int>> walls;
    Tank player;
    std::vector<Tank> enemies;
    std::vector<Bullet> bullets;
    int score = 0;
    bool gameOver = false;
    float enemySpawnTimer = 0.f, enemySpawnDelay = 2.f;
    float playerMoveDelay = 0.2f, enemyMoveDelay = 0.25f;
    float bulletSpeed = 12.f;
    sf::Font font;
};