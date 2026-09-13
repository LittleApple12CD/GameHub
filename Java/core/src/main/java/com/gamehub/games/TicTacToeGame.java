package com.gamehub.games;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.gamehub.Main;
import com.gamehub.screens.GameScreen;

public class TicTacToeGame implements GameInterface {
    private TicTacToeScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new TicTacToeScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Tic Tac Toe";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class TicTacToeScreen extends GameScreen {
        private static final int GAME_SIZE = 500;
        private static final int CELL_SIZE = GAME_SIZE / 3;

        private String[][] board;
        private String current;
        private String winner;
        private int moveCount;
        private int offsetX, offsetY;

        public TicTacToeScreen(Main game, TicTacToeGame parent) {
            super(game, parent, "Tic Tac Toe");
            this.gameWidth = GAME_SIZE;
            this.gameHeight = GAME_SIZE;
            reset();
        }

        public void reset() {
            board = new String[3][3];
            for (int r = 0; r < 3; r++) {
                for (int c = 0; c < 3; c++) {
                    board[r][c] = "";
                }
            }
            current = "X";
            winner = "";
            moveCount = 0;
            gameOver = false;
        }

        private String checkWinner() {
            for (int i = 0; i < 3; i++) {
                if (!board[i][0].isEmpty() && board[i][0].equals(board[i][1]) &&
                    board[i][1].equals(board[i][2])) return board[i][0];
                if (!board[0][i].isEmpty() && board[0][i].equals(board[1][i]) &&
                    board[1][i].equals(board[2][i])) return board[0][i];
            }
            if (!board[0][0].isEmpty() && board[0][0].equals(board[1][1]) &&
                board[1][1].equals(board[2][2])) return board[0][0];
            if (!board[0][2].isEmpty() && board[0][2].equals(board[1][1]) &&
                board[1][1].equals(board[2][0])) return board[0][2];
            return "";
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isButtonJustPressed(Input.Buttons.LEFT) && !gameOver) {
                int x = (int)((Gdx.input.getX() - offsetX) / CELL_SIZE);
                int y = (int)((Gdx.input.getY() - offsetY) / CELL_SIZE);
                if (x >= 0 && x < 3 && y >= 0 && y < 3 && board[y][x].isEmpty()) {
                    board[y][x] = current;
                    moveCount++;
                    winner = checkWinner();
                    if (!winner.isEmpty()) {
                        gameOver = true;
                    } else if (moveCount == 9) {
                        gameOver = true;
                    } else {
                        current = current.equals("X") ? "O" : "X";
                    }
                }
            }
        }

        @Override
        protected void update(float delta) {}

        @Override
        protected void draw(float delta) {
            int winW = Gdx.graphics.getWidth();
            int winH = Gdx.graphics.getHeight();
            offsetX = (winW - GAME_SIZE) / 2;
            offsetY = (winH - GAME_SIZE) / 2;

            // 背景
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.08f, 0.08f, 0.12f, 1);
            shapeRenderer.rect(offsetX, offsetY, GAME_SIZE, GAME_SIZE);

            // 格子
            shapeRenderer.setColor(0.2f, 0.2f, 0.25f, 1);
            for (int r = 0; r < 3; r++) {
                for (int c = 0; c < 3; c++) {
                    float drawY = offsetY + (gameHeight - (r + 1) * CELL_SIZE);
                    drawRoundedRect(
                        offsetX + c * CELL_SIZE + 2,
                        drawY + 2,
                        CELL_SIZE - 4, CELL_SIZE - 4, 8
                    );
                }
            }

            // 线条
            shapeRenderer.setColor(0.31f, 0.31f, 0.39f, 1);
            for (int i = 1; i < 3; i++) {
                float lineY1 = offsetY + i * CELL_SIZE;
                float lineY2 = offsetY + GAME_SIZE - i * CELL_SIZE;
                shapeRenderer.rectLine(offsetX + i * CELL_SIZE, offsetY + 5,
                    offsetX + i * CELL_SIZE, offsetY + GAME_SIZE - 5, 3);
                shapeRenderer.rectLine(offsetX + 5, lineY1,
                    offsetX + GAME_SIZE - 5, lineY1, 3);
            }
            shapeRenderer.end();

            // 棋子
            game.batch.begin();
            for (int r = 0; r < 3; r++) {
                for (int c = 0; c < 3; c++) {
                    if (!board[r][c].isEmpty()) {
                        Color color = board[r][c].equals("X") ?
                            new Color(0.24f, 0.78f, 1f, 1) :
                            new Color(1f, 0.78f, 0.24f, 1);
                        font.setColor(color);
                        float drawY = offsetY + (gameHeight - (r + 1) * CELL_SIZE);
                        font.draw(game.batch, board[r][c],
                            offsetX + c * CELL_SIZE + CELL_SIZE / 2 - 25,
                            drawY + CELL_SIZE / 2 + 20
                        );
                    }
                }
            }
            game.batch.end();

            // 状态栏
            int statusY = offsetY + GAME_SIZE + 30;
            game.batch.begin();
            if (gameOver) {
                String text = winner.isEmpty() ? "Draw!" : "Player " + winner + " Wins!";
                Color color = winner.isEmpty() ? Color.YELLOW : Color.GREEN;
                font.setColor(color);
                font.draw(game.batch, text, offsetX + GAME_SIZE / 2 - 100, statusY);
                font.setColor(Color.WHITE);
                smallFont.draw(game.batch, "Press R to Restart", offsetX + GAME_SIZE / 2 - 70, statusY + 40);
            } else {
                font.setColor(Color.WHITE);
                font.draw(game.batch, "Player " + current + "'s Turn", offsetX + GAME_SIZE / 2 - 90, statusY);
            }
            game.batch.end();
        }
    }
}