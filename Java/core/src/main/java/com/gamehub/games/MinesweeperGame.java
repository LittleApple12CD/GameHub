package com.gamehub.games;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.gamehub.Main;
import com.gamehub.screens.GameScreen;

import java.util.Random;

public class MinesweeperGame implements GameInterface {
    private MinesweeperScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new MinesweeperScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Minesweeper";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class MinesweeperScreen extends GameScreen {
        private static final int ROWS = 16;
        private static final int COLS = 16;
        private static final int MINES = 40;
        private static final int CELL_SIZE = 30;

        private int[][] board;
        private boolean[][] revealed;
        private boolean[][] flagged;
        private boolean firstClick;
        private boolean won;
        private int offsetX, offsetY;
        private Random random;

        public MinesweeperScreen(Main game, MinesweeperGame parent) {
            super(game, parent, "Minesweeper");
            this.gameWidth = COLS * CELL_SIZE;
            this.gameHeight = ROWS * CELL_SIZE;
            random = new Random();
            reset();
        }

        public void reset() {
            board = new int[ROWS][COLS];
            revealed = new boolean[ROWS][COLS];
            flagged = new boolean[ROWS][COLS];
            firstClick = true;
            gameOver = false;
            won = false;
            score = 0;
            for (int r = 0; r < ROWS; r++) {
                for (int c = 0; c < COLS; c++) {
                    board[r][c] = 0;
                }
            }
        }

        private void placeMines(int safeX, int safeY) {
            int placed = 0;
            while (placed < MINES) {
                int x = random.nextInt(COLS);
                int y = random.nextInt(ROWS);
                if (board[y][x] == -1) continue;
                if (Math.abs(x - safeX) <= 1 && Math.abs(y - safeY) <= 1) continue;
                board[y][x] = -1;
                placed++;
            }
            for (int y = 0; y < ROWS; y++) {
                for (int x = 0; x < COLS; x++) {
                    if (board[y][x] == -1) continue;
                    int count = 0;
                    for (int dy = -1; dy <= 1; dy++) {
                        for (int dx = -1; dx <= 1; dx++) {
                            int nx = x + dx, ny = y + dy;
                            if (nx >= 0 && nx < COLS && ny >= 0 && ny < ROWS && board[ny][nx] == -1) count++;
                        }
                    }
                    board[y][x] = count;
                }
            }
        }

        private void reveal(int x, int y) {
            if (x < 0 || x >= COLS || y < 0 || y >= ROWS) return;
            if (revealed[y][x] || flagged[y][x]) return;
            revealed[y][x] = true;
            if (board[y][x] == -1) {
                gameOver = true;
                return;
            }
            if (board[y][x] == 0) {
                for (int dy = -1; dy <= 1; dy++) {
                    for (int dx = -1; dx <= 1; dx++) {
                        reveal(x + dx, y + dy);
                    }
                }
            }
            checkWin();
        }

        private void checkWin() {
            int revealedCount = 0;
            for (int y = 0; y < ROWS; y++) {
                for (int x = 0; x < COLS; x++) {
                    if (revealed[y][x]) revealedCount++;
                }
            }
            if (revealedCount == ROWS * COLS - MINES) {
                won = true;
                gameOver = true;
            }
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isButtonJustPressed(Input.Buttons.LEFT) && !gameOver) {
                int x = (int)((Gdx.input.getX() - offsetX) / CELL_SIZE);
                int y = (int)((Gdx.input.getY() - offsetY) / CELL_SIZE);
                if (x >= 0 && x < COLS && y >= 0 && y < ROWS) {
                    if (firstClick) {
                        placeMines(x, y);
                        firstClick = false;
                    }
                    if (!flagged[y][x]) reveal(x, y);
                }
            }
            if (Gdx.input.isButtonJustPressed(Input.Buttons.RIGHT) && !gameOver) {
                int x = (int)((Gdx.input.getX() - offsetX) / CELL_SIZE);
                int y = (int)((Gdx.input.getY() - offsetY) / CELL_SIZE);
                if (x >= 0 && x < COLS && y >= 0 && y < ROWS) {
                    if (!revealed[y][x]) flagged[y][x] = !flagged[y][x];
                }
            }
        }

        @Override
        protected void update(float delta) {}

        @Override
        protected void draw(float delta) {
            int winW = Gdx.graphics.getWidth();
            int winH = Gdx.graphics.getHeight();
            offsetX = (winW - COLS * CELL_SIZE) / 2;
            offsetY = (winH - ROWS * CELL_SIZE) / 2;

            // 背景 + 格子
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.08f, 0.08f, 0.12f, 1);
            shapeRenderer.rect(offsetX, offsetY, COLS * CELL_SIZE, ROWS * CELL_SIZE);

            for (int y = 0; y < ROWS; y++) {
                for (int x = 0; x < COLS; x++) {
                    float drawY = offsetY + (gameHeight - (y + 1) * CELL_SIZE);
                    float drawX = offsetX + x * CELL_SIZE + 1;
                    float size = CELL_SIZE - 2;

                    if (revealed[y][x]) {
                        if (board[y][x] == -1) {
                            shapeRenderer.setColor(0.8f, 0.2f, 0.2f, 1);
                            drawRoundedRect(drawX, drawY + 1, size, size, 4);
                            shapeRenderer.setColor(0, 0, 0, 1);
                            shapeRenderer.rectLine(drawX, drawY + 1, drawX + size, drawY + size + 1, 2);
                            shapeRenderer.rectLine(drawX + size, drawY + 1, drawX, drawY + size + 1, 2);
                        } else {
                            shapeRenderer.setColor(0.75f, 0.75f, 0.78f, 1);
                            drawRoundedRect(drawX, drawY + 1, size, size, 4);
                        }
                    } else {
                        if (flagged[y][x]) {
                            shapeRenderer.setColor(0.82f, 0.82f, 0.24f, 1);
                        } else {
                            shapeRenderer.setColor(0.31f, 0.31f, 0.41f, 1);
                        }
                        drawRoundedRect(drawX, drawY + 1, size, size, 4);
                    }
                }
            }
            shapeRenderer.end();

            // 文字（数字和 F）
            game.batch.begin();

            // 周围雷数
            for (int y = 0; y < ROWS; y++) {
                for (int x = 0; x < COLS; x++) {
                    if (revealed[y][x] && board[y][x] > 0) {
                        float drawY = offsetY + (gameHeight - (y + 1) * CELL_SIZE);
                        Color numColor;
                        switch (board[y][x]) {
                            case 1: numColor = new Color(0.2f, 0.4f, 1f, 1); break;
                            case 2: numColor = new Color(0.2f, 0.7f, 0.2f, 1); break;
                            case 3: numColor = new Color(1f, 0.3f, 0.3f, 1); break;
                            case 4: numColor = new Color(0.2f, 0.2f, 0.7f, 1); break;
                            case 5: numColor = new Color(0.7f, 0.2f, 0.2f, 1); break;
                            case 6: numColor = new Color(0.2f, 0.7f, 0.7f, 1); break;
                            case 7: numColor = new Color(0.1f, 0.1f, 0.1f, 1); break;
                            case 8: numColor = new Color(0.5f, 0.5f, 0.5f, 1); break;
                            default: numColor = Color.WHITE;
                        }
                        smallFont.setColor(numColor);
                        smallFont.draw(game.batch, String.valueOf(board[y][x]),
                            offsetX + x * CELL_SIZE + CELL_SIZE / 2 - 6,
                            drawY + CELL_SIZE / 2 + 8
                        );
                    }
                }
            }

            // 旗帜 F
            for (int y = 0; y < ROWS; y++) {
                for (int x = 0; x < COLS; x++) {
                    if (!revealed[y][x] && flagged[y][x]) {
                        float drawY = offsetY + (gameHeight - (y + 1) * CELL_SIZE);
                        smallFont.setColor(0, 0, 0, 1);
                        smallFont.draw(game.batch, "F",
                            offsetX + x * CELL_SIZE + CELL_SIZE / 2 - 6,
                            drawY + CELL_SIZE / 2 + 8
                        );
                    }
                }
            }

            game.batch.end();

            if (gameOver) {
                drawGameOverOverlay(won ? "YOU WIN!" : "GAME OVER", "Press R to Restart");
            }
        }
    }
}