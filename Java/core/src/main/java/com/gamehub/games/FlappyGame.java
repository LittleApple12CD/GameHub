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
import java.util.Random;

public class FlappyGame implements GameInterface {
    private FlappyScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new FlappyScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Flappy Bird";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class FlappyScreen extends GameScreen {
        private static final int GAME_WIDTH = 400;
        private static final int GAME_HEIGHT = 600;
        private static final int BIRD_RADIUS = 15;
        private static final int HITBOX_RADIUS = 11;
        private static final int PIPE_WIDTH = 60;
        private static final int PIPE_GAP = 160;
        private static final float GRAVITY = 0.2f;
        private static final float JUMP_FORCE = -7f;
        private static final float PIPE_SPEED = 2f;
        private static final int PIPE_DELAY = 90;

        private float birdY;
        private float birdVY;
        private int birdX;
        private ArrayList<Pipe> pipes;
        private float pipeTimer;
        private int offsetX, offsetY;
        private Random random;

        static class Pipe {
            float x;
            float gapY;
            boolean passed;
            Pipe(float x, float gapY) {
                this.x = x;
                this.gapY = gapY;
                this.passed = false;
            }
        }

        public FlappyScreen(Main game, FlappyGame parent) {
            super(game, parent, "Flappy Bird");
            this.gameWidth = GAME_WIDTH;
            this.gameHeight = GAME_HEIGHT;
            random = new Random();
            birdX = 80;
            reset();
        }

        public void reset() {
            birdY = GAME_HEIGHT / 2f;
            birdVY = 0;
            pipes = new ArrayList<>();
            pipeTimer = 0;
            score = 0;
            gameOver = false;
        }

        private void addPipe() {
            float gapY = MathUtils.random(80, GAME_HEIGHT - PIPE_GAP - 80);
            pipes.add(new Pipe(GAME_WIDTH, gapY));
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isKeyJustPressed(Input.Keys.SPACE) && !gameOver) {
                birdVY = JUMP_FORCE;
            }
            if (Gdx.input.isButtonJustPressed(Input.Buttons.LEFT) && !gameOver) {
                birdVY = JUMP_FORCE;
            }
        }

        @Override
        protected void update(float delta) {
            if (gameOver) return;

            birdVY += GRAVITY;
            birdY += birdVY;

            if (birdY - BIRD_RADIUS < 0) {
                birdY = BIRD_RADIUS;
                birdVY = 0;
            }
            if (birdY + BIRD_RADIUS > GAME_HEIGHT) {
                gameOver = true;
                return;
            }

            pipeTimer += delta * 60;
            if (pipeTimer >= PIPE_DELAY) {
                pipeTimer = 0;
                addPipe();
            }

            for (Pipe pipe : pipes) {
                pipe.x -= PIPE_SPEED;

                if (!pipe.passed && pipe.x + PIPE_WIDTH < birdX) {
                    pipe.passed = true;
                    score++;
                }

                if (pipe.x < birdX + HITBOX_RADIUS && pipe.x + PIPE_WIDTH > birdX - HITBOX_RADIUS) {
                    if (birdY - HITBOX_RADIUS < pipe.gapY ||
                        birdY + HITBOX_RADIUS > pipe.gapY + PIPE_GAP) {
                        gameOver = true;
                        return;
                    }
                }
            }

            pipes.removeIf(pipe -> pipe.x + PIPE_WIDTH < 0);
        }

        @Override
        protected void draw(float delta) {
            int winW = Gdx.graphics.getWidth();
            int winH = Gdx.graphics.getHeight();
            offsetX = (winW - GAME_WIDTH) / 2;
            offsetY = (winH - GAME_HEIGHT) / 2;

            // 背景
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.53f, 0.81f, 0.92f, 1);
            shapeRenderer.rect(offsetX, offsetY, GAME_WIDTH, GAME_HEIGHT);
            shapeRenderer.end();

            // 管道
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.16f, 0.78f, 0.16f, 1);
            for (Pipe pipe : pipes) {
                float topPipeY = offsetY + GAME_HEIGHT - pipe.gapY;
                drawRoundedRect(offsetX + pipe.x, topPipeY, PIPE_WIDTH, pipe.gapY, 6);

                float bottomPipeH = GAME_HEIGHT - pipe.gapY - PIPE_GAP;
                drawRoundedRect(offsetX + pipe.x, offsetY, PIPE_WIDTH, bottomPipeH, 6);

                shapeRenderer.setColor(0.12f, 0.63f, 0.12f, 1);
                drawRoundedRect(offsetX + pipe.x - 6, topPipeY, PIPE_WIDTH + 12, 22, 6);
                drawRoundedRect(offsetX + pipe.x - 6, offsetY + bottomPipeH - 22, PIPE_WIDTH + 12, 22, 6);
                shapeRenderer.setColor(0.16f, 0.78f, 0.16f, 1);
            }
            shapeRenderer.end();

            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(1, 1, 0.2f, 1);
            float birdDrawY = offsetY + (gameHeight - birdY);
            shapeRenderer.circle(offsetX + birdX, birdDrawY, BIRD_RADIUS);
            shapeRenderer.setColor(0, 0, 0, 1);
            shapeRenderer.circle(offsetX + birdX + 8, birdDrawY + 5, 4);
            shapeRenderer.setColor(1, 1, 1, 1);
            shapeRenderer.circle(offsetX + birdX + 10, birdDrawY + 3, 2);
            shapeRenderer.end();

            game.batch.begin();
            font.draw(game.batch, String.valueOf(score),
                offsetX + GAME_WIDTH / 2 - 20,
                offsetY + GAME_HEIGHT - 40);
            game.batch.end();

            if (gameOver) {
                drawGameOverOverlay("GAME OVER", "Press R to Restart");
            }
        }
    }
}
