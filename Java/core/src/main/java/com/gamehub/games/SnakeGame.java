package com.gamehub.games;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.badlogic.gdx.math.MathUtils;
import com.gamehub.Main;
import com.gamehub.screens.GameScreen;

import java.util.LinkedList;

public class SnakeGame implements GameInterface {
    private SnakeScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new SnakeScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Snake";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class SnakeScreen extends GameScreen {
        private static final int CELL_SIZE = 30;
        private static final int GAME_WIDTH = 600;
        private static final int GAME_HEIGHT = 600;

        private LinkedList<int[]> snake;
        private int[] direction;
        private int[] nextDirection;
        private int[] food;
        private float moveTimer;
        private final float MOVE_DELAY = 0.15f;
        private int gridWidth, gridHeight;

        public SnakeScreen(Main game, SnakeGame parent) {
            super(game, parent, "Snake");
            this.gameWidth = GAME_WIDTH;
            this.gameHeight = GAME_HEIGHT;
            gridWidth = GAME_WIDTH / CELL_SIZE;
            gridHeight = GAME_HEIGHT / CELL_SIZE;
            reset();
        }

        public void reset() {
            snake = new LinkedList<>();
            snake.add(new int[]{gridWidth / 2, gridHeight / 2});
            direction = new int[]{1, 0};
            nextDirection = new int[]{1, 0};
            spawnFood();
            score = 0;
            gameOver = false;
            moveTimer = 0;
        }

        private void spawnFood() {
            do {
                food = new int[]{
                    MathUtils.random(gridWidth - 1),
                    MathUtils.random(gridHeight - 1)
                };
            } while (isOnSnake(food[0], food[1]));
        }

        private boolean isOnSnake(int x, int y) {
            for (int[] seg : snake) {
                if (seg[0] == x && seg[1] == y) return true;
            }
            return false;
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isKeyJustPressed(Input.Keys.UP) && direction[1] != 1) {
                nextDirection = new int[]{0, -1};
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.DOWN) && direction[1] != -1) {
                nextDirection = new int[]{0, 1};
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.LEFT) && direction[0] != 1) {
                nextDirection = new int[]{-1, 0};
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.RIGHT) && direction[0] != -1) {
                nextDirection = new int[]{1, 0};
            }
        }

        @Override
        protected void update(float delta) {
            moveTimer += delta;
            if (moveTimer >= MOVE_DELAY) {
                moveTimer = 0;
                direction = nextDirection;
                int[] head = snake.getFirst();
                int newX = head[0] + direction[0];
                int newY = head[1] + direction[1];

                if (newX < 0 || newX >= gridWidth || newY < 0 || newY >= gridHeight) {
                    gameOver = true;
                    return;
                }
                for (int[] seg : snake) {
                    if (seg[0] == newX && seg[1] == newY) {
                        gameOver = true;
                        return;
                    }
                }
                snake.addFirst(new int[]{newX, newY});
                if (newX == food[0] && newY == food[1]) {
                    score += 10;
                    spawnFood();
                } else {
                    snake.removeLast();
                }
            }
        }

        @Override
        protected void draw(float delta) {
            offsetX = (Gdx.graphics.getWidth() - GAME_WIDTH) / 2;
            offsetY = (Gdx.graphics.getHeight() - GAME_HEIGHT) / 2;

            // 背景
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.08f, 0.08f, 0.12f, 1);
            shapeRenderer.rect(offsetX, offsetY, GAME_WIDTH, GAME_HEIGHT);
            shapeRenderer.end();

            // 蛇
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            for (int i = 0; i < snake.size(); i++) {
                int[] seg = snake.get(i);
                if (i == 0) {
                    shapeRenderer.setColor(0.24f, 0.86f, 0.24f, 1);
                } else {
                    shapeRenderer.setColor(0.16f, 0.71f, 0.16f, 1);
                }
                // 翻转 Y：gameHeight - (seg[1] + 1) * CELL_SIZE
                float drawY = offsetY + (gameHeight - (seg[1] + 1) * CELL_SIZE);
                drawRoundedRect(
                    offsetX + seg[0] * CELL_SIZE + 1,
                    drawY + 1,
                    CELL_SIZE - 2, CELL_SIZE - 2, 6
                );
            }

            // 食物
            shapeRenderer.setColor(1, 0.24f, 0.24f, 1);
            float drawY = offsetY + (gameHeight - (food[1] + 1) * CELL_SIZE);
            drawRoundedRect(
                offsetX + food[0] * CELL_SIZE + 3,
                drawY + 3,
                CELL_SIZE - 6, CELL_SIZE - 6, 8
            );
            shapeRenderer.end();

            // 分数
            game.batch.begin();
            font.draw(game.batch, "Score: " + score, offsetX + 10, offsetY + GAME_HEIGHT - 10);
            game.batch.end();

            if (gameOver) {
                drawGameOverOverlay("GAME OVER", "Press R to Restart");
            }
        }
    }
}