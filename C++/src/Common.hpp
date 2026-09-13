#pragma once
#include <SFML/Graphics.hpp>
#include <optional>
#include <string>
#include <cmath>
#include <algorithm>

// 颜色常量
inline const sf::Color COLOR_BG(25, 25, 35);
inline const sf::Color COLOR_WHITE(255, 255, 255);
inline const sf::Color COLOR_HOVER(80, 200, 255);
inline const sf::Color COLOR_BUTTON(55, 55, 75);
inline const sf::Color COLOR_BORDER(120, 120, 160);
inline const sf::Color COLOR_QUIT(200, 50, 50);
inline const sf::Color COLOR_QUIT_HOVER(255, 80, 80);

inline constexpr unsigned SCREEN_WIDTH = 850;
inline constexpr unsigned SCREEN_HEIGHT = 650;

class Game {
public:
    virtual ~Game() = default;
    virtual bool run(sf::RenderWindow& window) = 0;
};

// ============================================================
// 圆角矩形绘制
// ============================================================
inline void drawRect(sf::RenderTarget& target,
                     const sf::FloatRect& rect,
                     sf::Color color,
                     float radius = 0.f,
                     float outline = 0.f,
                     sf::Color outlineColor = sf::Color::Transparent) {

    if (radius <= 0.f) {
        sf::RectangleShape shape(sf::Vector2f(rect.size.x, rect.size.y));
        shape.setPosition(rect.position);
        shape.setFillColor(color);
        if (outline > 0.f) {
            shape.setOutlineThickness(outline);
            shape.setOutlineColor(outlineColor);
        }
        target.draw(shape);
        return;
    }
    radius = std::min(radius, std::min(rect.size.x, rect.size.y) / 2.f);

    const int segments = 8;
    const int pointsPerCorner = segments + 1;

    sf::ConvexShape shape;
    shape.setPointCount(pointsPerCorner * 4);
    shape.setFillColor(color);
    if (outline > 0.f) {
        shape.setOutlineThickness(outline);
        shape.setOutlineColor(outlineColor);
    }

    const float x = rect.position.x;
    const float y = rect.position.y;
    const float w = rect.size.x;
    const float h = rect.size.y;
    constexpr float PI = 3.14159265358979323846f;

    int idx = 0;

    for (int i = 0; i <= segments; ++i) {
        float a = PI + (PI / 2.f) * static_cast<float>(i) / segments;
        shape.setPoint(idx++, {x + radius + radius * std::cos(a),
                               y + radius + radius * std::sin(a)});
    }
    for (int i = 0; i <= segments; ++i) {
        float a = 1.5f * PI + (PI / 2.f) * static_cast<float>(i) / segments;
        shape.setPoint(idx++, {x + w - radius + radius * std::cos(a),
                               y + radius + radius * std::sin(a)});
    }
    for (int i = 0; i <= segments; ++i) {
        float a = 0.f + (PI / 2.f) * static_cast<float>(i) / segments;
        shape.setPoint(idx++, {x + w - radius + radius * std::cos(a),
                               y + h - radius + radius * std::sin(a)});
    }
    for (int i = 0; i <= segments; ++i) {
        float a = PI / 2.f + (PI / 2.f) * static_cast<float>(i) / segments;
        shape.setPoint(idx++, {x + radius + radius * std::cos(a),
                               y + h - radius + radius * std::sin(a)});
    }

    target.draw(shape);
}