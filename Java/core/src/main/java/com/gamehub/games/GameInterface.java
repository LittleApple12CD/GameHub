package com.gamehub.games;

import com.badlogic.gdx.Screen;
import com.gamehub.Main;

public interface GameInterface {
    Screen createScreen(Main game);
    String getName();
    void reset();
}