#include "Common.hpp"
#include "Menu.hpp"

int main() {
    sf::RenderWindow window(sf::VideoMode({SCREEN_WIDTH, SCREEN_HEIGHT}), "Game Hub");
    window.setFramerateLimit(60);

    sf::Image icon;
    if (icon.loadFromFile("assets/icons/icon.png"))
        window.setIcon(icon.getSize(), icon.getPixelsPtr());

    Menu menu;
    menu.run(window);
    return 0;
}