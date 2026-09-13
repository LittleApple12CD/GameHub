#include "TetrisGame.hpp"
#include <random>
#include <iostream>

static std::mt19937 rng(std::random_device{}());

bool TetrisGame::run(sf::RenderWindow& window) {
    if (!font.openFromFile("assets/fonts/arial.ttf")) std::cerr << "font missing\n";
    shapes = {
        {'I', {{0,0,0,0},{1,1,1,1},{0,0,0,0},{0,0,0,0}}},
        {'O', {{1,1},{1,1}}},
        {'T', {{0,1,0},{1,1,1},{0,0,0}}},
        {'S', {{0,1,1},{1,1,0},{0,0,0}}},
        {'Z', {{1,1,0},{0,1,1},{0,0,0}}},
        {'L', {{1,0,0},{1,1,1},{0,0,0}}},
        {'J', {{0,0,1},{1,1,1},{0,0,0}}}
    };
    colors = {
        {'I', sf::Color(0,240,240)}, {'O', sf::Color(240,240,0)},
        {'T', sf::Color(180,0,240)}, {'S', sf::Color(0,240,0)},
        {'Z', sf::Color(240,0,0)},   {'L', sf::Color(240,160,0)},
        {'J', sf::Color(0,0,240)}
    };
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

void TetrisGame::reset() {
    board.assign(ROWS, std::vector<char>(COLS, 0));
    score = 0; gameOver = false; fallTimer = 0.f;
    spawnPiece();
}

void TetrisGame::spawnPiece() {
    std::vector<char> keys = {'I','O','T','S','Z','L','J'};
    std::uniform_int_distribution<int> d(0, (int)keys.size() - 1);
    currentPiece = keys[d(rng)];
    pieceShape = shapes[currentPiece];
    pieceX = COLS / 2 - (int)pieceShape[0].size() / 2;
    pieceY = 0;
    if (checkCollision(pieceShape, pieceX, pieceY)) gameOver = true;
}

bool TetrisGame::checkCollision(const std::vector<std::vector<int>>& s, int x, int y) {
    for (size_t r = 0; r < s.size(); ++r)
        for (size_t c = 0; c < s[r].size(); ++c)
            if (s[r][c]) {
                int bx = x + (int)c, by = y + (int)r;
                if (bx < 0 || bx >= COLS || by >= ROWS) return true;
                if (by >= 0 && board[by][bx]) return true;
            }
    return false;
}

void TetrisGame::rotatePiece() {
    int n = (int)pieceShape.size();
    std::vector<std::vector<int>> rot(n, std::vector<int>(n, 0));
    for (int y = 0; y < n; ++y)
        for (int x = 0; x < n; ++x)
            rot[x][n - 1 - y] = pieceShape[y][x];
    if (!checkCollision(rot, pieceX, pieceY)) pieceShape = rot;
}

void TetrisGame::lockPiece() {
    for (size_t r = 0; r < pieceShape.size(); ++r)
        for (size_t c = 0; c < pieceShape[r].size(); ++c)
            if (pieceShape[r][c]) {
                int by = pieceY + (int)r, bx = pieceX + (int)c;
                if (by >= 0) board[by][bx] = currentPiece;
            }
    clearLines();
    spawnPiece();
}

void TetrisGame::clearLines() {
    int lines = 0;
    for (int r = ROWS - 1; r >= 0; --r) {
        bool full = true;
        for (int c = 0; c < COLS; ++c) if (!board[r][c]) { full = false; break; }
        if (full) {
            board.erase(board.begin() + r);
            board.insert(board.begin(), std::vector<char>(COLS, 0));
            lines++; r++;
        }
    }
    int pts[] = {0, 100, 300, 500, 800};
    if (lines) score += pts[std::min(lines, 4)];
}

bool TetrisGame::handleEvents(sf::RenderWindow& window) {
    while (const auto event = window.pollEvent()) {
        if (event->is<sf::Event::Closed>()) { window.close(); return false; }
        if (const auto* key = event->getIf<sf::Event::KeyPressed>()) {
            using K = sf::Keyboard::Key;
            if (key->code == K::Escape) return false;
            if (key->code == K::R) reset();
            if (!gameOver) {
                if (key->code == K::Left && !checkCollision(pieceShape, pieceX - 1, pieceY)) pieceX--;
                else if (key->code == K::Right && !checkCollision(pieceShape, pieceX + 1, pieceY)) pieceX++;
                else if (key->code == K::Down && !checkCollision(pieceShape, pieceX, pieceY + 1)) pieceY++;
                else if (key->code == K::Space || key->code == K::Up) rotatePiece();
            }
        }
    }
    return true;
}

void TetrisGame::update(float dt) {
    if (gameOver) return;
    fallTimer += dt;
    if (fallTimer >= fallDelay) {
        fallTimer = 0.f;
        if (!checkCollision(pieceShape, pieceX, pieceY + 1)) pieceY++;
        else lockPiece();
    }
}

void TetrisGame::draw(sf::RenderWindow& window) {
    window.clear(COLOR_BG);
    sf::RectangleShape bg(sf::Vector2f(WIDTH, HEIGHT));
    bg.setPosition({(float)offsetX, (float)offsetY});
    bg.setFillColor(sf::Color(20, 20, 30));
    window.draw(bg);

    for (int y = 0; y < ROWS; ++y)
        for (int x = 0; x < COLS; ++x)
            if (board[y][x]) {
                sf::FloatRect r({(float)(offsetX + x * CELL + 1),
                                 (float)(offsetY + y * CELL + 1)},
                                {(float)(CELL - 2), (float)(CELL - 2)});
                drawRect(window, r, colors[board[y][x]], 4.f);
            }

    if (!gameOver)
        for (size_t r = 0; r < pieceShape.size(); ++r)
            for (size_t c = 0; c < pieceShape[r].size(); ++c)
                if (pieceShape[r][c]) {
                    sf::FloatRect s({(float)(offsetX + (pieceX + (int)c) * CELL + 1),
                                     (float)(offsetY + (pieceY + (int)r) * CELL + 1)},
                                    {(float)(CELL - 2), (float)(CELL - 2)});
                    drawRect(window, s, colors[currentPiece], 5.f);
                }

    sf::Text st(font, "Score: " + std::to_string(score), 28);
    st.setFillColor(COLOR_WHITE);
    st.setPosition({(float)(offsetX + WIDTH + 20), (float)offsetY + 20});
    window.draw(st);

    if (gameOver) {
        sf::RectangleShape ov(sf::Vector2f(WIDTH, HEIGHT));
        ov.setPosition({(float)offsetX, (float)offsetY});
        ov.setFillColor(sf::Color(0,0,0,180));
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