#pragma once
#include "Common.hpp"
#include <vector>
#include <functional>
#include <memory>

class Menu {
public:
    Menu();
    void run(sf::RenderWindow& window);

private:
    void drawMenu(sf::RenderWindow& window, int selected, bool quitSelected);
    sf::Font font;
    std::vector<std::string> names;
    std::vector<std::function<std::unique_ptr<Game>()>> factories;
    int selected = 0;
    bool quitSelected = false;
};