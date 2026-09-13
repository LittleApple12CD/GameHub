package com.gamehub.screens;

import com.badlogic.gdx.Gdx;
import com.badlogic.gdx.Input;
import com.badlogic.gdx.Screen;
import com.badlogic.gdx.graphics.Color;
import com.badlogic.gdx.graphics.GL20;
import com.badlogic.gdx.graphics.Pixmap;
import com.badlogic.gdx.graphics.Texture;
import com.badlogic.gdx.graphics.g2d.BitmapFont;
import com.badlogic.gdx.graphics.g2d.freetype.FreeTypeFontGenerator;
import com.badlogic.gdx.scenes.scene2d.Actor;
import com.badlogic.gdx.scenes.scene2d.Stage;
import com.badlogic.gdx.scenes.scene2d.ui.Label;
import com.badlogic.gdx.scenes.scene2d.ui.Table;
import com.badlogic.gdx.scenes.scene2d.ui.TextButton;
import com.badlogic.gdx.scenes.scene2d.utils.ChangeListener;
import com.badlogic.gdx.scenes.scene2d.utils.TextureRegionDrawable;
import com.badlogic.gdx.utils.viewport.ScreenViewport;
import com.gamehub.Main;
import com.gamehub.games.*;

public class MainMenuScreen implements Screen {
    private final Main game;
    private Stage stage;
    private BitmapFont titleFont;
    private BitmapFont buttonFont;

    private final String[] gameNames = {
        "Snake", "Tetris", "Minesweeper",
        "Flappy Bird", "Tank", "Breakout", "Tic Tac Toe"
    };

    private final Class<?>[] gameClasses = {
        SnakeGame.class, TetrisGame.class, MinesweeperGame.class,
        FlappyGame.class, TankGame.class, BreakoutGame.class,
        TicTacToeGame.class
    };

    public MainMenuScreen(final Main game) {
        this.game = game;
        this.stage = new Stage(new ScreenViewport());
        Gdx.input.setInputProcessor(stage);

        loadFonts();
        buildMenu();
    }

    private void loadFonts() {
        try {
            FreeTypeFontGenerator generator = new FreeTypeFontGenerator(
                Gdx.files.internal("assets/fonts/arial.ttf")
            );

            FreeTypeFontGenerator.FreeTypeFontParameter titleParam = new FreeTypeFontGenerator.FreeTypeFontParameter();
            titleParam.size = 48;
            titleParam.color = Color.WHITE;
            titleFont = generator.generateFont(titleParam);

            FreeTypeFontGenerator.FreeTypeFontParameter buttonParam = new FreeTypeFontGenerator.FreeTypeFontParameter();
            buttonParam.size = 38;
            buttonParam.color = Color.WHITE;
            buttonFont = generator.generateFont(buttonParam);

            generator.dispose();
        } catch (Exception e) {
            e.printStackTrace();
            titleFont = new BitmapFont();
            buttonFont = new BitmapFont();
        }
    }

    private Pixmap createRoundedRectPixmap(int width, int height, int radius, Color fill, Color border) {
        Pixmap pixmap = new Pixmap(width, height, Pixmap.Format.RGBA8888);

        pixmap.setColor(fill);
        pixmap.fillRectangle(radius, 0, width - 2 * radius, height);
        pixmap.fillRectangle(0, radius, width, height - 2 * radius);

        for (int y = 0; y <= radius; y++) {
            for (int x = 0; x <= radius; x++) {
                if (x * x + y * y <= radius * radius) {
                    pixmap.drawPixel(radius - x, radius - y);
                    pixmap.drawPixel(width - radius + x - 1, radius - y);
                    pixmap.drawPixel(radius - x, height - radius + y - 1);
                    pixmap.drawPixel(width - radius + x - 1, height - radius + y - 1);
                }
            }
        }

        pixmap.setColor(border);
        for (int x = radius; x < width - radius; x++) {
            pixmap.drawPixel(x, 0);
            pixmap.drawPixel(x, height - 1);
        }
        for (int y = radius; y < height - radius; y++) {
            pixmap.drawPixel(0, y);
            pixmap.drawPixel(width - 1, y);
        }
        for (int y = 0; y <= radius; y++) {
            for (int x = 0; x <= radius; x++) {
                int dist = x * x + y * y;
                if (dist >= (radius - 1) * (radius - 1) && dist <= radius * radius) {
                    pixmap.drawPixel(radius - x, radius - y);
                    pixmap.drawPixel(width - radius + x - 1, radius - y);
                    pixmap.drawPixel(radius - x, height - radius + y - 1);
                    pixmap.drawPixel(width - radius + x - 1, height - radius + y - 1);
                }
            }
        }

        return pixmap;
    }

    private void buildMenu() {
        Table table = new Table();
        table.setFillParent(true);
        table.center();

        Label.LabelStyle titleStyle = new Label.LabelStyle();
        titleStyle.font = titleFont;
        titleStyle.fontColor = Color.WHITE;
        Label title = new Label("GAME HUB", titleStyle);
        table.add(title).padBottom(60).row();

        TextButton.TextButtonStyle gameStyle = new TextButton.TextButtonStyle();
        gameStyle.font = buttonFont;
        gameStyle.fontColor = Color.WHITE;
        gameStyle.overFontColor = new Color(0.3f, 0.8f, 1f, 1f);

        Pixmap bgPixmap = createRoundedRectPixmap(280, 46, 12,
            new Color(55f/255f, 55f/255f, 75f/255f, 1),
            new Color(120f/255f, 120f/255f, 160f/255f, 1)
        );
        gameStyle.up = new TextureRegionDrawable(new Texture(bgPixmap));
        bgPixmap.dispose();

        Pixmap hoverPixmap = createRoundedRectPixmap(280, 46, 12,
            new Color(60f/255f, 80f/255f, 120f/255f, 1),
            new Color(150f/255f, 220f/255f, 255f/255f, 1)
        );
        gameStyle.over = new TextureRegionDrawable(new Texture(hoverPixmap));
        hoverPixmap.dispose();

        for (int i = 0; i < gameNames.length; i++) {
            final int index = i;
            TextButton button = new TextButton(gameNames[i], gameStyle);
            button.addListener(new ChangeListener() {
                @Override
                public void changed(ChangeEvent event, Actor actor) {
                    launchGame(index);
                }
            });
            table.add(button).width(280).height(46).pad(3).row();
        }

        TextButton.TextButtonStyle quitStyle = new TextButton.TextButtonStyle();
        quitStyle.font = buttonFont;
        quitStyle.fontColor = Color.WHITE;
        quitStyle.overFontColor = Color.WHITE;
        
        Pixmap quitBg = createRoundedRectPixmap(120, 40, 10,
            new Color(200f/255f, 50f/255f, 50f/255f, 1),
            new Color(150f/255f, 40f/255f, 40f/255f, 1)
        );
        quitStyle.up = new TextureRegionDrawable(new Texture(quitBg));
        quitBg.dispose();

        Pixmap quitHover = createRoundedRectPixmap(120, 40, 10,
            new Color(255f/255f, 80f/255f, 80f/255f, 1),
            new Color(220f/255f, 60f/255f, 60f/255f, 1)
        );
        quitStyle.over = new TextureRegionDrawable(new Texture(quitHover));
        quitHover.dispose();

        TextButton quitButton = new TextButton("EXIT", quitStyle);
        quitButton.addListener(new ChangeListener() {
            @Override
            public void changed(ChangeEvent event, Actor actor) {
                Gdx.app.exit();
            }
        });
        table.add(quitButton).width(120).height(40).pad(20).row();

        stage.addActor(table);
    }

    private void launchGame(int index) {
        try {
            Gdx.input.setInputProcessor(null);
            GameInterface gameInstance = (GameInterface) gameClasses[index].getDeclaredConstructor().newInstance();
            game.setScreen(gameInstance.createScreen(game));
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void render(float delta) {
        Gdx.gl.glClearColor(0.1f, 0.1f, 0.14f, 1);
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT);

        if (Gdx.input.isKeyJustPressed(Input.Keys.ESCAPE)) {
            Gdx.app.exit();
        }

        stage.act(delta);
        stage.draw();
    }

    @Override
    public void resize(int width, int height) {
        stage.getViewport().update(width, height, true);
    }

    @Override
    public void dispose() {
        stage.dispose();
        if (titleFont != null) titleFont.dispose();
        if (buttonFont != null) buttonFont.dispose();
    }

    @Override
    public void show() {}
    @Override
    public void pause() {}
    @Override
    public void resume() {}
    @Override
    public void hide() {}
}
