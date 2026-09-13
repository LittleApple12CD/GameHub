package com.gamehub.games;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.badlogic.gdx.math.MathUtils;
import com.gamehub.Main;
import com.gamehub.screens.GameScreen;

import java.util.ArrayList;

public class BreakoutGame implements GameInterface {
    private BreakoutScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new BreakoutScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Breakout";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class BreakoutScreen extends GameScreen {
        private static final int GAME_WIDTH = 600;
        private static final int GAME_HEIGHT = 600;
        private static final int PADDLE_WIDTH = 120;
        private static final int PADDLE_HEIGHT = 16;
        private static final int BALL_SIZE = 20;
        private static final float PADDLE_SPEED = 6f;
        private static final float BALL_SPEED = 4f;

        private int offsetX, offsetY;
        private float paddleX;
        private float ballX, ballY;
        private float ballDx, ballDy;
        private ArrayList<Brick> bricks;
        private int lives;
        private boolean waiting;

        static class Brick {
            float x, y, w, h;
            Color color;
            boolean alive;
            Brick(float x, float y, float w, float h, Color color) {
                this.x = x;
                this.y = y;
                this.w = w;
                this.h = h;
                this.color = color;
                this.alive = true;
            }
        }

        public BreakoutScreen(Main game, BreakoutGame parent) {
            super(game, parent, "Breakout");
            this.gameWidth = GAME_WIDTH;
            this.gameHeight = GAME_HEIGHT;
            reset();
        }

        public void reset() {
            paddleX = GAME_WIDTH / 2f - PADDLE_WIDTH / 2f;
            ballX = GAME_WIDTH / 2f - BALL_SIZE / 2f;
            ballY = 70;
            ballDx = MathUtils.randomBoolean() ? BALL_SPEED : -BALL_SPEED;
            ballDy = BALL_SPEED;
            lives = 3;
            score = 0;
            gameOver = false;
            waiting = true;
            createBricks();
        }

        private void createBricks() {
            bricks = new ArrayList<>();
            int rows = 5, cols = 8;
            float brickW = (GAME_WIDTH - 20) / cols - 4;
            float brickH = 22;
            Color[] colors = {
                new Color(0.9f, 0.2f, 0.2f, 1),
                new Color(0.9f, 0.59f, 0.2f, 1),
                new Color(0.9f, 0.9f, 0.2f, 1),
                new Color(0.2f, 0.9f, 0.2f, 1),
                new Color(0.2f, 0.59f, 0.9f, 1)
            };
            for (int row = 0; row < rows; row++) {
                for (int col = 0; col < cols; col++) {
                    float x = 10 + col * (brickW + 4);
                    float y = GAME_HEIGHT - 40 - (row + 1) * (brickH + 4);
                    bricks.add(new Brick(x, y, brickW, brickH, colors[row % colors.length]));
                }
            }
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isKeyJustPressed(Input.Keys.SPACE) && waiting && !gameOver) {
                waiting = false;
            }
        }

        @Override
        protected void update(float delta) {
            if (gameOver || waiting) return;

            if (Gdx.input.isKeyPressed(Input.Keys.LEFT) && paddleX > 0) {
                paddleX -= PADDLE_SPEED;
            }
            if (Gdx.input.isKeyPressed(Input.Keys.RIGHT) && paddleX + PADDLE_WIDTH < GAME_WIDTH) {
                paddleX += PADDLE_SPEED;
            }

            ballX += ballDx;
            ballY += ballDy;

            if (ballX <= 0 || ballX + BALL_SIZE >= GAME_WIDTH) ballDx = -ballDx;
            if (ballY + BALL_SIZE >= GAME_HEIGHT) ballDy = -ballDy;
            if (ballY <= 0) {
                lives--;
                if (lives <= 0) {
                    gameOver = true;
                } else {
                    waiting = true;
                    ballX = GAME_WIDTH / 2f - BALL_SIZE / 2f;
                    ballY = 70;
                    ballDx = MathUtils.randomBoolean() ? BALL_SPEED : -BALL_SPEED;
                    ballDy = BALL_SPEED;
                }
                return;
            }

            if (ballX + BALL_SIZE > paddleX && ballX < paddleX + PADDLE_WIDTH &&
                ballY < 40 + PADDLE_HEIGHT && ballY + BALL_SIZE > 40) {
                ballDy = Math.abs(ballDy);
                float hitPos = (ballX + BALL_SIZE / 2 - paddleX - PADDLE_WIDTH / 2) / (PADDLE_WIDTH / 2);
                ballDx = hitPos * BALL_SPEED * 0.9f;
                if (Math.abs(ballDx) < 2) ballDx = (ballDx >= 0 ? 3 : -3);
            }

            for (Brick brick : bricks) {
                if (!brick.alive) continue;
                if (ballX + BALL_SIZE > brick.x && ballX < brick.x + brick.w &&
                    ballY + BALL_SIZE > brick.y && ballY < brick.y + brick.h) {
                    brick.alive = false;
                    score += 10;

                    float overlapTop = ballY + BALL_SIZE - brick.y;
                    float overlapBottom = brick.y + brick.h - ballY;
                    float overlapLeft = ballX + BALL_SIZE - brick.x;
                    float overlapRight = brick.x + brick.w - ballX;

                    float minOverlap = Math.min(Math.min(overlapTop, overlapBottom),
                        Math.min(overlapLeft, overlapRight));

                    if (minOverlap == overlapTop || minOverlap == overlapBottom) {
                        ballDy = -ballDy;
                    } else {
                        ballDx = -ballDx;
                    }
                    break;
                }
            }

            boolean allDead = true;
            for (Brick b : bricks) if (b.alive) { allDead = false; break; }
            if (allDead) gameOver = true;
        }

        @Override
        protected void draw(float delta) {
            int winW = Gdx.graphics.getWidth();
            int winH = Gdx.graphics.getHeight();
            offsetX = (winW - GAME_WIDTH) / 2;
            offsetY = (winH - GAME_HEIGHT) / 2;

            // 背景
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.08f, 0.08f, 0.12f, 1);
            shapeRenderer.rect(offsetX, offsetY, GAME_WIDTH, GAME_HEIGHT);
            shapeRenderer.end();

            // 砖块
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            for (Brick brick : bricks) {
                if (!brick.alive) continue;
                shapeRenderer.setColor(brick.color);
                float drawY = offsetY + brick.y;
                shapeRenderer.setColor(1, 1, 1, 0.4f);
                drawRoundedRect(offsetX + brick.x - 1, drawY - 1, brick.w + 2, brick.h + 2, 8);
                shapeRenderer.setColor(brick.color);
                drawRoundedRect(offsetX + brick.x, drawY, brick.w, brick.h, 7);
            }
            shapeRenderer.end();

            // 挡板
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(1, 1, 1, 1);
            float paddleDrawY = offsetY + 24;
            drawRoundedRect(offsetX + paddleX, paddleDrawY, PADDLE_WIDTH, PADDLE_HEIGHT, 8);
            shapeRenderer.end();

            // 球
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(1, 1, 0.47f, 1);
            float ballDrawY = offsetY + ballY + BALL_SIZE / 2;
            shapeRenderer.circle(offsetX + ballX + BALL_SIZE / 2, ballDrawY, BALL_SIZE / 2);
            shapeRenderer.end();

            // 分数和生命
            game.batch.begin();
            font.draw(game.batch, "Score: " + score, offsetX + 10, offsetY + GAME_HEIGHT - 10);
            font.draw(game.batch, "Lives: " + lives, offsetX + GAME_WIDTH - 100, offsetY + GAME_HEIGHT - 10);
            game.batch.end();

            if (waiting && !gameOver) {
                game.batch.begin();
                font.draw(game.batch, "Press SPACE to start",
                    offsetX + GAME_WIDTH / 2 - 120,
                    offsetY + GAME_HEIGHT / 2 + 30);
                game.batch.end();
            }

            if (gameOver) {
                drawGameOverOverlay(lives <= 0 ? "GAME OVER" : "YOU WIN!", "Press R to Restart");
            }
        }
    }
}