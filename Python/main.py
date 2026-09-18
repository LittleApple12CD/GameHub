import pygame
import sys
from game_snake import SnakeGame
from game_tetris import TetrisGame
from game_minesweeper import MinesweeperGame
from game_flappy import FlappyGame
from game_tank import TankGame
from game_breakout import BreakoutGame
from game_tic_tac_toe import TicTacToeGame

pygame.init()
SCREEN_WIDTH = 850
SCREEN_HEIGHT = 650

try:
    icon = pygame.image.load("assets/icons/icon.png")
    pygame.display.set_icon(icon)
except FileNotFoundError:
    icon = pygame.Surface((32, 32))
    icon.fill((100, 200, 255))
    pygame.display.set_icon(icon)

screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
pygame.display.set_caption("Game Hub")
clock = pygame.time.Clock()
font = pygame.font.Font(None, 52)
font_small = pygame.font.Font(None, 38)

COLOR_BG = (25, 25, 35)
COLOR_WHITE = (255, 255, 255)
COLOR_HOVER = (80, 200, 255)
COLOR_BUTTON = (55, 55, 75)
COLOR_BORDER = (120, 120, 160)
COLOR_QUIT = (200, 50, 50)
COLOR_QUIT_HOVER = (255, 80, 80)

games = [
    ("SnakeGame", SnakeGame),
    ("TetrisGame", TetrisGame),
    ("Minesweeper", MinesweeperGame),
    ("Flappy Bird", FlappyGame),
    ("TankBattle", TankGame),
    ("Breakout", BreakoutGame),
    ("Tic Tac Toe", TicTacToeGame)
]

def draw_menu(selected, quit_selected):
    screen.fill(COLOR_BG)
    
    title = font.render("GAME HUB", True, COLOR_WHITE)
    title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, 55))
    screen.blit(title, title_rect)
    
    for i, (name, _) in enumerate(games):
        y = 130 + i * 52
        rect = pygame.Rect(SCREEN_WIDTH // 2 - 140, y, 280, 46)
        color = COLOR_HOVER if i == selected else COLOR_BUTTON
        pygame.draw.rect(screen, color, rect, border_radius=12)
        pygame.draw.rect(screen, COLOR_BORDER, rect, 2, border_radius=12)
        text = font_small.render(name, True, COLOR_WHITE)
        text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, y + 23))
        screen.blit(text, text_rect)
    
    quit_rect = pygame.Rect(SCREEN_WIDTH // 2 - 60, 130 + len(games) * 52 + 20, 120, 40)
    quit_color = COLOR_QUIT_HOVER if quit_selected else COLOR_QUIT
    pygame.draw.rect(screen, quit_color, quit_rect, border_radius=10)
    pygame.draw.rect(screen, COLOR_BORDER, quit_rect, 2, border_radius=10)
    quit_text = font_small.render("EXIT", True, COLOR_WHITE)
    quit_text_rect = quit_text.get_rect(center=(SCREEN_WIDTH // 2, 130 + len(games) * 52 + 40))
    screen.blit(quit_text, quit_text_rect)
    
    pygame.display.flip()

def main():
    selected = 0
    quit_selected = False
    running = True
    
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                sys.exit()
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    pygame.quit()
                    sys.exit()
                elif event.key == pygame.K_UP:
                    if quit_selected:
                        quit_selected = False
                        selected = len(games) - 1
                    else:
                        selected = (selected - 1) % len(games)
                elif event.key == pygame.K_DOWN:
                    if selected == len(games) - 1:
                        quit_selected = True
                    else:
                        selected = (selected + 1) % len(games)
                        quit_selected = False
                elif event.key == pygame.K_RETURN:
                    if quit_selected:
                        pygame.quit()
                        sys.exit()
                    else:
                        game_class = games[selected][1]
                        game = game_class(screen)
                        game.run()
            if event.type == pygame.MOUSEMOTION:
                quit_rect = pygame.Rect(SCREEN_WIDTH // 2 - 60, 130 + len(games) * 52 + 20, 120, 40)
                if quit_rect.collidepoint(event.pos):
                    quit_selected = True
                else:
                    quit_selected = False
                    for i, (name, _) in enumerate(games):
                        y = 130 + i * 52
                        rect = pygame.Rect(SCREEN_WIDTH // 2 - 140, y, 280, 46)
                        if rect.collidepoint(event.pos):
                            selected = i
                            break
            if event.type == pygame.MOUSEBUTTONDOWN:
                quit_rect = pygame.Rect(SCREEN_WIDTH // 2 - 60, 130 + len(games) * 52 + 20, 120, 40)
                if quit_rect.collidepoint(event.pos):
                    pygame.quit()
                    sys.exit()
                for i, (name, _) in enumerate(games):
                    y = 130 + i * 52
                    rect = pygame.Rect(SCREEN_WIDTH // 2 - 140, y, 280, 46)
                    if rect.collidepoint(event.pos):
                        game_class = games[i][1]
                        game = game_class(screen)
                        game.run()
        
        draw_menu(selected, quit_selected)
        clock.tick(60)

if __name__ == "__main__":
    main()
