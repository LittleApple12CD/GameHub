import pygame
import random

class MinesweeperGame:
    def __init__(self, screen):
        self.screen = screen
        self.rows, self.cols = 16, 16
        self.mines = 40
        self.cell_size = 30
        self.width = self.cols * self.cell_size
        self.height = self.rows * self.cell_size
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 22)
        self.font_big = pygame.font.Font(None, 40)
        self.reset()
    
    def reset(self):
        self.board = [[0] * self.cols for _ in range(self.rows)]
        self.revealed = [[False] * self.cols for _ in range(self.rows)]
        self.flagged = [[False] * self.cols for _ in range(self.rows)]
        self.game_over = False
        self.won = False
        self.first_click = True
        self.mine_count = self.mines
    
    def place_mines(self, safe_x, safe_y):
        placed = 0
        while placed < self.mines:
            x = random.randint(0, self.cols - 1)
            y = random.randint(0, self.rows - 1)
            if self.board[y][x] == -1:
                continue
            if abs(x - safe_x) <= 1 and abs(y - safe_y) <= 1:
                continue
            self.board[y][x] = -1
            placed += 1
        
        for y in range(self.rows):
            for x in range(self.cols):
                if self.board[y][x] == -1:
                    continue
                count = 0
                for dy in [-1, 0, 1]:
                    for dx in [-1, 0, 1]:
                        nx, ny = x + dx, y + dy
                        if 0 <= nx < self.cols and 0 <= ny < self.rows:
                            if self.board[ny][nx] == -1:
                                count += 1
                self.board[y][x] = count
    
    def reveal(self, x, y):
        if not (0 <= x < self.cols and 0 <= y < self.rows):
            return
        if self.revealed[y][x] or self.flagged[y][x]:
            return
        
        self.revealed[y][x] = True
        
        if self.board[y][x] == -1:
            self.game_over = True
            return
        
        if self.board[y][x] == 0:
            for dy in [-1, 0, 1]:
                for dx in [-1, 0, 1]:
                    self.reveal(x + dx, y + dy)
        
        revealed_count = sum(sum(row) for row in self.revealed)
        if revealed_count == self.rows * self.cols - self.mines:
            self.won = True
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return False
                if event.key == pygame.K_r:
                    self.reset()
            if event.type == pygame.MOUSEBUTTONDOWN and not self.game_over and not self.won:
                x = (event.pos[0] - self.offset_x) // self.cell_size
                y = (event.pos[1] - self.offset_y) // self.cell_size
                if not (0 <= x < self.cols and 0 <= y < self.rows):
                    continue
                
                if event.button == 1:
                    if self.first_click:
                        self.place_mines(x, y)
                        self.first_click = False
                    if not self.flagged[y][x]:
                        self.reveal(x, y)
                elif event.button == 3:
                    if not self.revealed[y][x]:
                        self.flagged[y][x] = not self.flagged[y][x]
        return True
    
    def update(self):
        pass
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((30, 30, 42))
        
        for y in range(self.rows):
            for x in range(self.cols):
                rect = pygame.Rect(x * self.cell_size + 1, y * self.cell_size + 1,
                                   self.cell_size - 2, self.cell_size - 2)
                
                if self.revealed[y][x]:
                    if self.board[y][x] == -1:
                        pygame.draw.rect(game_surface, (200, 50, 50), rect, border_radius=4)
                        pygame.draw.circle(game_surface, (0, 0, 0),
                                           (x * self.cell_size + 15, y * self.cell_size + 15), 9)
                    else:
                        pygame.draw.rect(game_surface, (190, 190, 200), rect, border_radius=4)
                        if self.board[y][x] > 0:
                            colors = [(0,0,0), (0,0,255), (0,150,0), (255,0,0),
                                      (0,0,180), (150,0,0), (0,150,150), (0,0,0)]
                            text = self.font.render(str(self.board[y][x]), True, colors[self.board[y][x]])
                            game_surface.blit(text, (x * self.cell_size + 10, y * self.cell_size + 6))
                else:
                    color = (80, 80, 105) if not self.flagged[y][x] else (210, 210, 60)
                    pygame.draw.rect(game_surface, color, rect, border_radius=4)
                    if self.flagged[y][x]:
                        text = self.font.render("F", True, (255, 50, 50))
                        game_surface.blit(text, (x * self.cell_size + 10, y * self.cell_size + 5))
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        if self.game_over:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(150)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            text = self.font_big.render("GAME OVER", True, (255, 255, 255))
            text_rect = text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 - 20))
            self.screen.blit(text, text_rect)
            text2 = self.font.render("Press R to restart", True, (255, 255, 255))
            text2_rect = text2.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 + 30))
            self.screen.blit(text2, text2_rect)
        elif self.won:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(150)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            text = self.font_big.render("YOU WIN!", True, (0, 255, 0))
            text_rect = text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 - 20))
            self.screen.blit(text, text_rect)
            text2 = self.font.render("Press R to restart", True, (255, 255, 255))
            text2_rect = text2.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 + 30))
            self.screen.blit(text2, text2_rect)
        
        pygame.display.flip()
    
    def run(self):
        running = True
        while running:
            if not self.handle_events():
                break
            self.update()
            self.draw()
            self.clock.tick(60)