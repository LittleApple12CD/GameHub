package com.gamehub;

import com.badlogic.gdx.Game;
import com.badlogic.gdx.graphics.g2d.SpriteBatch;
import com.gamehub.screens.MainMenuScreen;

public class Main extends Game {
    public SpriteBatch batch;
    
    @Override
    public void create() {
        batch = new SpriteBatch();
        setScreen(new MainMenuScreen(this));
    }
    
    @Override
    public void render() {
        super.render();
    }
    
    @Override
    public void dispose() {
        batch.dispose();
        super.dispose();
    }
}