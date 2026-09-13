package com.gamehub.screens;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.GL20;
import com.badlogic.gdx.graphics.g2d.BitmapFont;
import com.badlogic.gdx.graphics.g2d.GlyphLayout;
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator;
import com.badlogic.gdx.graphics.glutils.ShapeRenderer;
import com.gamehub.Main;
import com.gamehub.games.GameInterface;

public abstract class GameScreen implements Screen {
    protected final Main game;
    protected final GameInterface gameInstance;
    protected final String gameName;
    protected ShapeRenderer shapeRenderer;
    protected BitmapFont font;
    protected BitmapFont smallFont;
    protected boolean running = true;
    protected boolean gameOver = false;
    protected int score = 0;
    protected int gameWidth;
    protected int gameHeight;
    protected int offsetX, offsetY;

    public GameScreen(Main game, GameInterface gameInstance, String gameName) {
        this.game = game;
        this.gameInstance = gameInstance;
        this.gameName = gameName;
        this.shapeRenderer = new ShapeRenderer();
        loadFonts();
    }

    private void loadFonts() {
        try {
            FreeTypeFontGenerator generator = new FreeTypeFontGenerator(
                Gdx.files.internal("assets/fonts/arial.ttf")
            );

            FreeTypeFontGenerator.FreeTypeFontParameter param = new FreeTypeFontGenerator.FreeTypeFontParameter();
            param.size = 36;
            param.color = Color.WHITE;
            font = generator.generateFont(param);

            param.size = 24;
            smallFont = generator.generateFont(param);

            generator.dispose();
        } catch (Exception e) {
            font = new BitmapFont();
            smallFont = new BitmapFont();
        }
    }

    @Override
    public void render(float delta) {
        Gdx.gl.glClearColor(0.1f, 0.1f, 0.14f, 1);
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT);

        handleInput();

        if (!gameOver && running) {
            update(delta);
        }

        draw(delta);
    }

    protected void handleInput() {
        if (Gdx.input.isKeyJustPressed(Input.Keys.ESCAPE)) {
            goToMenu();
        }
        if (Gdx.input.isKeyJustPressed(Input.Keys.R)) {
            resetGame();
        }
    }

    protected void goToMenu() {
        Gdx.input.setInputProcessor(null);
        gameInstance.reset();
        game.setScreen(new MainMenuScreen(game));
    }

    protected void resetGame() {
        gameInstance.reset();
        gameOver = false;
        score = 0;
    }

    protected abstract void update(float delta);
    protected abstract void draw(float delta);

    protected void drawRoundedRect(float x, float y, float width, float height, float radius) {
        if (radius <= 0) {
            shapeRenderer.rect(x, y, width, height);
            return;
        }
        shapeRenderer.rect(x + radius, y, width - 2 * radius, height);
        shapeRenderer.rect(x, y + radius, width, height - 2 * radius);
        shapeRenderer.circle(x + radius, y + radius, radius);
        shapeRenderer.circle(x + width - radius, y + radius, radius);
        shapeRenderer.circle(x + radius, y + height - radius, radius);
        shapeRenderer.circle(x + width - radius, y + height - radius, radius);
    }

    protected void drawGameOverOverlay(String title, String sub) {
        int x = (Gdx.graphics.getWidth() - gameWidth) / 2;
        int y = (Gdx.graphics.getHeight() - gameHeight) / 2;

        shapeRenderer.begin(ShapeRenderer.ShapeType.Filled);
        shapeRenderer.setColor(0, 0, 0, 0.7f);
        shapeRenderer.rect(x, y, gameWidth, gameHeight);
        shapeRenderer.end();

        GlyphLayout layout = new GlyphLayout();

        game.batch.begin();
        font.setColor(Color.WHITE);
        layout.setText(font, title);
        font.draw(game.batch, title,
            x + gameWidth / 2 - layout.width / 2,
            y + gameHeight / 2 + 20);

        smallFont.setColor(Color.WHITE);
        layout.setText(smallFont, sub);
        smallFont.draw(game.batch, sub,
            x + gameWidth / 2 - layout.width / 2,
            y + gameHeight / 2 - 20);
        game.batch.end();
    }

    @Override
    public void dispose() {
        shapeRenderer.dispose();
        if (font != null) font.dispose();
        if (smallFont != null) smallFont.dispose();
    }

    @Override
    public void show() {}
    @Override
    public void resize(int width, int height) {}
    @Override
    public void pause() {}
    @Override
    public void resume() {}
    @Override
    public void hide() {}
}