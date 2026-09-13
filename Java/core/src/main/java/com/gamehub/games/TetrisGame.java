package com.gamehub.games;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.gamehub.Main;
import com.gamehub.screens.GameScreen;

import java.util.HashMap;
import java.util.Map;
import java.util.Random;

public class TetrisGame implements GameInterface {
    private TetrisScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new TetrisScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Tetris";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class TetrisScreen extends GameScreen {
        private static final int COLS = 10;
        private static final int ROWS = 20;
        private static final int CELL_SIZE = 30;
        private static final int GAME_WIDTH = COLS * CELL_SIZE;
        private static final int GAME_HEIGHT = ROWS * CELL_SIZE;
        private static final int INFO_WIDTH = 140;

        private String[][] board;
        private int[][][] shapes;
        private String[] shapeNames;
        private Map<String, Color> shapeColors;
        private String currentPiece;
        private int[][] currentShape;
        private int pieceX, pieceY;
        private float fallTimer;
        private float fallDelay;
        private Random random;
        private int offsetX, offsetY;

        public TetrisScreen(Main game, TetrisGame parent) {
            super(game, parent, "Tetris");
            this.gameWidth = GAME_WIDTH;
            this.gameHeight = GAME_HEIGHT;
            initShapes();
            reset();
        }

        private void initShapes() {
            shapes = new int[][][]{
                {{0,0,0,0},{1,1,1,1},{0,0,0,0},{0,0,0,0}},
                {{1,1},{1,1}},
                {{0,1,0},{1,1,1},{0,0,0}},
                {{0,1,1},{1,1,0},{0,0,0}},
                {{1,1,0},{0,1,1},{0,0,0}},
                {{1,0,0},{1,1,1},{0,0,0}},
                {{0,0,1},{1,1,1},{0,0,0}}
            };
            shapeNames = new String[]{"I","O","T","S","Z","L","J"};
            shapeColors = new HashMap<>();
            shapeColors.put("I", Color.CYAN);
            shapeColors.put("O", Color.YELLOW);
            shapeColors.put("T", Color.PURPLE);
            shapeColors.put("S", Color.GREEN);
            shapeColors.put("Z", Color.RED);
            shapeColors.put("L", Color.ORANGE);
            shapeColors.put("J", Color.BLUE);
            random = new Random();
            fallDelay = 0.5f;
        }

        public void reset() {
            board = new String[ROWS][COLS];
            for (int r = 0; r < ROWS; r++) {
                for (int c = 0; c < COLS; c++) {
                    board[r][c] = "";
                }
            }
            score = 0;
            gameOver = false;
            fallTimer = 0;
            spawnPiece();
        }

        private void spawnPiece() {
            if (gameOver) return;
            int idx = random.nextInt(shapes.length);
            currentPiece = shapeNames[idx];
            currentShape = shapes[idx];
            pieceX = COLS / 2 - currentShape[0].length / 2;
            pieceY = 0;
            if (checkCollision(currentShape, pieceX, pieceY)) {
                gameOver = true;
            }
        }

        private boolean checkCollision(int[][] shape, int x, int y) {
            for (int row = 0; row < shape.length; row++) {
                for (int col = 0; col < shape[row].length; col++) {
                    if (shape[row][col] == 1) {
                        int bx = x + col;
                        int by = y + row;
                        if (bx < 0 || bx >= COLS || by >= ROWS) return true;
                        if (by >= 0 && !board[by][bx].isEmpty()) return true;
                    }
                }
            }
            return false;
        }

        private void rotatePiece() {
            int[][] rotated = new int[currentShape[0].length][currentShape.length];
            for (int row = 0; row < currentShape.length; row++) {
                for (int col = 0; col < currentShape[row].length; col++) {
                    rotated[col][currentShape.length - 1 - row] = currentShape[row][col];
                }
            }
            if (!checkCollision(rotated, pieceX, pieceY)) {
                currentShape = rotated;
            }
        }

        private void lockPiece() {
            for (int row = 0; row < currentShape.length; row++) {
                for (int col = 0; col < currentShape[row].length; col++) {
                    if (currentShape[row][col] == 1) {
                        int by = pieceY + row;
                        int bx = pieceX + col;
                        if (by >= 0) board[by][bx] = currentPiece;
                    }
                }
            }
            clearLines();
            spawnPiece();
        }

        private void clearLines() {
            int lines = 0;
            for (int row = ROWS - 1; row >= 0; row--) {
                boolean full = true;
                for (int col = 0; col < COLS; col++) {
                    if (board[row][col].isEmpty()) {
                        full = false;
                        break;
                    }
                }
                if (full) {
                    for (int r = row; r > 0; r--) {
                        board[r] = board[r-1].clone();
                    }
                    board[0] = new String[COLS];
                    for (int c = 0; c < COLS; c++) board[0][c] = "";
                    lines++;
                    row++;
                }
            }
            if (lines > 0) {
                score += lines * 100;
            }
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isKeyJustPressed(Input.Keys.LEFT) && !gameOver) {
                if (!checkCollision(currentShape, pieceX - 1, pieceY)) pieceX--;
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.RIGHT) && !gameOver) {
                if (!checkCollision(currentShape, pieceX + 1, pieceY)) pieceX++;
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.DOWN) && !gameOver) {
                if (!checkCollision(currentShape, pieceX, pieceY + 1)) pieceY++;
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.UP) && !gameOver) {
                rotatePiece();
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.SPACE) && !gameOver) {
                while (!checkCollision(currentShape, pieceX, pieceY + 1)) pieceY++;
                lockPiece();
            }
        }

        @Override
        protected void update(float delta) {
            if (gameOver) return;
            fallTimer += delta;
            if (fallTimer >= fallDelay) {
                fallTimer = 0;
                if (!checkCollision(currentShape, pieceX, pieceY + 1)) {
                    pieceY++;
                } else {
                    lockPiece();
                }
            }
        }

        @Override
        protected void draw(float delta) {
            int winW = Gdx.graphics.getWidth();
            int winH = Gdx.graphics.getHeight();
            offsetX = (winW - GAME_WIDTH - INFO_WIDTH) / 2;
            offsetY = (winH - GAME_HEIGHT) / 2;

            // 背景
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.08f, 0.08f, 0.12f, 1);
            shapeRenderer.rect(offsetX, offsetY, GAME_WIDTH, GAME_HEIGHT);
            shapeRenderer.end();

            // 已放置方块
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            for (int row = 0; row < ROWS; row++) {
                for (int col = 0; col < COLS; col++) {
                    if (!board[row][col].isEmpty()) {
                        Color color = shapeColors.get(board[row][col]);
                        shapeRenderer.setColor(color);
                        float drawY = offsetY + (gameHeight - (row + 1) * CELL_SIZE);
                        drawRoundedRect(
                            offsetX + col * CELL_SIZE + 1,
                            drawY + 1,
                            CELL_SIZE - 2, CELL_SIZE - 2, 4
                        );
                    }
                }
            }

            // 当前方块
            if (!gameOver) {
                for (int row = 0; row < currentShape.length; row++) {
                    for (int col = 0; col < currentShape[row].length; col++) {
                        if (currentShape[row][col] == 1) {
                            Color color = shapeColors.get(currentPiece);
                            shapeRenderer.setColor(color);
                            int gameRow = pieceY + row;
                            float drawY = offsetY + (gameHeight - (gameRow + 1) * CELL_SIZE);
                            drawRoundedRect(
                                offsetX + (pieceX + col) * CELL_SIZE + 1,
                                drawY + 1,
                                CELL_SIZE - 2, CELL_SIZE - 2, 4
                            );
                        }
                    }
                }
            }
            shapeRenderer.end();

            // 分数
            game.batch.begin();
            font.draw(game.batch, "Score: " + score, offsetX + GAME_WIDTH + 20, offsetY + GAME_HEIGHT - 20);
            game.batch.end();

            if (gameOver) {
                drawGameOverOverlay("GAME OVER", "Press R to Restart");
            }
        }
    }
}