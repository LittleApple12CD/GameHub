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
import java.util.HashMap;
import java.util.Map;

public class TankGame implements GameInterface {
    private TankScreen screen;

    @Override
    public Screen createScreen(Main game) {
        screen = new TankScreen(game, this);
        return screen;
    }

    @Override
    public String getName() {
        return "Tank";
    }

    @Override
    public void reset() {
        if (screen != null) screen.reset();
    }

    public static class TankScreen extends GameScreen {
        private static final int GAME_SIZE = 600;
        private static final int CELL_SIZE = 30;
        private static final float BULLET_SPEED = 1.0f;
        private static final float MOVE_DELAY = 0.2f;

        private int gridSize;
        private Map<String, Boolean> walls;
        private PlayerTank player;
        private ArrayList<EnemyTank> enemies;
        private ArrayList<Bullet> bullets;
        private int enemySpawnTimer;
        private int enemySpawnDelay;
        private int offsetX, offsetY;

        static class PlayerTank {
            int x, y;
            int[] dir;
            boolean alive;
            float moveTimer;
            PlayerTank(int x, int y) {
                this.x = x;
                this.y = y;
                this.dir = new int[]{1, 0};
                this.alive = true;
                this.moveTimer = 0;
            }
        }

        static class EnemyTank {
            int x, y;
            int[] dir;
            boolean alive;
            float timer;
            float moveTimer;
            EnemyTank(int x, int y) {
                this.x = x;
                this.y = y;
                this.dir = new int[]{-1, 0};
                this.alive = true;
                this.timer = 0;
                this.moveTimer = 0;
            }
        }

        static class Bullet {
            float x, y;
            int[] dir;
            boolean isPlayer;
            boolean alive;
            Bullet(float x, float y, int[] dir, boolean isPlayer) {
                this.x = x;
                this.y = y;
                this.dir = dir.clone();
                this.isPlayer = isPlayer;
                this.alive = true;
            }
        }

        public TankScreen(Main game, TankGame parent) {
            super(game, parent, "Tank");
            this.gameWidth = GAME_SIZE;
            this.gameHeight = GAME_SIZE;
            gridSize = GAME_SIZE / CELL_SIZE;
            enemySpawnDelay = 120;
            reset();
        }

        public void reset() {
            walls = generateWalls();
            player = new PlayerTank(1, 1);
            enemies = new ArrayList<>();
            enemies.add(new EnemyTank(gridSize - 2, gridSize - 2));
            bullets = new ArrayList<>();
            enemySpawnTimer = 0;
            score = 0;
            gameOver = false;
        }

        private Map<String, Boolean> generateWalls() {
            Map<String, Boolean> walls = new HashMap<>();
            for (int i = 0; i < gridSize; i++) {
                walls.put(i + "," + 0, true);
                walls.put(i + "," + (gridSize - 1), true);
                walls.put(0 + "," + i, true);
                walls.put((gridSize - 1) + "," + i, true);
            }
            for (int i = 0; i < 12; i++) {
                int x = MathUtils.random(2, gridSize - 3);
                int y = MathUtils.random(2, gridSize - 3);
                if (!(x == 1 && y == 1) && !(x == gridSize - 2 && y == gridSize - 2)) {
                    walls.put(x + "," + y, true);
                }
            }
            return walls;
        }

        private boolean canMove(int x, int y, int[] dir, boolean isPlayer) {
            if (walls.containsKey(x + "," + y)) return false;
            if (isPlayer) {
                for (EnemyTank e : enemies) {
                    if (e.alive && e.x == x && e.y == y) return false;
                }
            } else {
                if (player.alive && player.x == x && player.y == y) return false;
                for (EnemyTank e : enemies) {
                    if (e.alive && e.x == x && e.y == y) return false;
                }
            }
            return x >= 0 && x < gridSize && y >= 0 && y < gridSize;
        }

        private void shootBullet(float x, float y, int[] dir, boolean isPlayer) {
            bullets.add(new Bullet(x + dir[0], y + dir[1], dir, isPlayer));
        }

        @Override
        protected void handleInput() {
            super.handleInput();
            if (Gdx.input.isKeyJustPressed(Input.Keys.SPACE) && !gameOver && player.alive) {
                shootBullet(player.x, player.y, player.dir, true);
            }
            if (Gdx.input.isKeyJustPressed(Input.Keys.SPACE) && gameOver) {
                reset();
            }
        }

        @Override
        protected void update(float delta) {
            if (gameOver) return;

            if (player.alive) {
                int dx = 0, dy = 0;
                // W 向上（-1），S 向下（+1）
                if (Gdx.input.isKeyPressed(Input.Keys.W)) dy = -1;
                else if (Gdx.input.isKeyPressed(Input.Keys.S)) dy = 1;
                else if (Gdx.input.isKeyPressed(Input.Keys.A)) dx = -1;
                else if (Gdx.input.isKeyPressed(Input.Keys.D)) dx = 1;

                if (dx != 0 || dy != 0) {
                    player.dir = new int[]{dx, dy};
                    player.moveTimer += delta;
                    if (player.moveTimer >= MOVE_DELAY) {
                        player.moveTimer = 0;
                        int newX = player.x + dx;
                        int newY = player.y + dy;
                        if (canMove(newX, newY, player.dir, true)) {
                            player.x = newX;
                            player.y = newY;
                        }
                    }
                } else {
                    player.moveTimer = 0;
                }
            }

            for (int i = bullets.size() - 1; i >= 0; i--) {
                Bullet b = bullets.get(i);
                if (!b.alive) continue;
                b.x += b.dir[0] * BULLET_SPEED;
                b.y += b.dir[1] * BULLET_SPEED;

                int bx = (int)b.x;
                int by = (int)b.y;

                if (bx < 0 || bx >= gridSize || by < 0 || by >= gridSize ||
                    walls.containsKey(bx + "," + by)) {
                    b.alive = false;
                    continue;
                }

                if (b.isPlayer) {
                    for (EnemyTank e : enemies) {
                        if (e.alive && e.x == bx && e.y == by) {
                            e.alive = false;
                            b.alive = false;
                            score += 10;
                            break;
                        }
                    }
                } else {
                    if (player.alive && player.x == bx && player.y == by) {
                        player.alive = false;
                        b.alive = false;
                        gameOver = true;
                    }
                }
            }
            bullets.removeIf(b -> !b.alive);

            for (EnemyTank e : enemies) {
                if (!e.alive) continue;
                e.timer += delta;
                if (e.timer >= 0.33f) {
                    e.timer = 0;
                    if (MathUtils.random() < 0.25f) {
                        int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
                        e.dir = dirs[MathUtils.random(3)];
                    }
                    e.moveTimer += delta;
                    if (e.moveTimer >= MOVE_DELAY) {
                        e.moveTimer = 0;
                        int newX = e.x + e.dir[0];
                        int newY = e.y + e.dir[1];
                        if (canMove(newX, newY, e.dir, false)) {
                            e.x = newX;
                            e.y = newY;
                        }
                    }
                }
                if (MathUtils.random() < 0.015f) {
                    shootBullet(e.x, e.y, e.dir, false);
                }
            }

            enemySpawnTimer++;
            if (enemySpawnTimer >= enemySpawnDelay) {
                enemySpawnTimer = 0;
                int aliveCount = 0;
                for (EnemyTank e : enemies) if (e.alive) aliveCount++;
                if (aliveCount < 3) {
                    enemies.add(new EnemyTank(gridSize - 2, gridSize - 2));
                }
            }
        }

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
            shapeRenderer.end();

            // 墙壁
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            shapeRenderer.setColor(0.39f, 0.39f, 0.51f, 1);
            for (String key : walls.keySet()) {
                String[] parts = key.split(",");
                int x = Integer.parseInt(parts[0]);
                int y = Integer.parseInt(parts[1]);
                float drawY = offsetY + (gameHeight - (y + 1) * CELL_SIZE);
                drawRoundedRect(
                    offsetX + x * CELL_SIZE + 1,
                    drawY + 1,
                    CELL_SIZE - 2, CELL_SIZE - 2, 4
                );
            }
            shapeRenderer.end();

            // 玩家和敌人坦克（统一包在 begin/end 里）
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            if (player.alive) {
                drawTank(player.x, player.y, Color.GREEN, player.dir);
            }
            for (EnemyTank e : enemies) {
                if (e.alive) drawTank(e.x, e.y, Color.RED, e.dir);
            }
            shapeRenderer.end();

            // 子弹
            shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
            for (Bullet b : bullets) {
                if (!b.alive) continue;
                if (b.isPlayer) {
                    shapeRenderer.setColor(1, 1, 0.31f, 1);
                } else {
                    shapeRenderer.setColor(1, 0.59f, 0.2f, 1);
                }
                float drawY = offsetY + (gameHeight - ((int)b.y + 1) * CELL_SIZE);
                drawRoundedRect(
                    offsetX + (int)b.x * CELL_SIZE + 10,
                    drawY + 10,
                    10, 10, 4
                );
            }
            shapeRenderer.end();

            // 分数
            game.batch.begin();
            font.draw(game.batch, "Score: " + score, offsetX + 10, offsetY + GAME_SIZE - 10);
            game.batch.end();

            if (gameOver) {
                drawGameOverOverlay("GAME OVER", "Press R or SPACE to Restart");
            }
        }

        private void drawTank(int x, int y, Color color, int[] dir) {
            int cx = offsetX + x * CELL_SIZE + CELL_SIZE / 2;
            float drawY = offsetY + (gameHeight - (y + 1) * CELL_SIZE);
            float drawCY = offsetY + (gameHeight - y * CELL_SIZE) - CELL_SIZE / 2;

            // 注意：这里不要 begin/end，由 draw() 统一管理
            shapeRenderer.setColor(color);
            // 主体
            drawRoundedRect(
                offsetX + x * CELL_SIZE + 2,
                drawY + 2,
                CELL_SIZE - 4, CELL_SIZE - 4, 4
            );

            // 炮管
            shapeRenderer.setColor(color.cpy().mul(0.8f));
            if (dir[0] == 1) {
                shapeRenderer.rect(cx + 4, drawCY - 4, 14, 8);
            } else if (dir[0] == -1) {
                shapeRenderer.rect(cx - 18, drawCY - 4, 14, 8);
            } else if (dir[1] == -1) {
                shapeRenderer.rect(cx - 4, drawCY + 4, 8, 14);
            } else {
                shapeRenderer.rect(cx - 4, drawCY - 18, 8, 14);
            }
        }
    }
}