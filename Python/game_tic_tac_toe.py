import pygame

class TicTacToeGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 500, 500
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.cell_size = self.width // 3
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 80)
        self.font_big = pygame.font.Font(None, 48)
        self.font_small = pygame.font.Font(None, 28)
        self.reset()
    
    def reset(self):
        self.board = [[''] * 3 for _ in range(3)]
        self.current_player = 'X'
        self.winner = None
        self.game_over = False
        self.move_count = 0
    
    def check_winner(self):
        for row in range(3):
            if self.board[row][0] and self.board[row][0] == self.board[row][1] == self.board[row][2]:
                return self.board[row][0]
        
        for col in range(3):
            if self.board[0][col] and self.board[0][col] == self.board[1][col] == self.board[2][col]:
                return self.board[0][col]
        
        if self.board[0][0] and self.board[0][0] == self.board[1][1] == self.board[2][2]:
            return self.board[0][0]
        
        if self.board[0][2] and self.board[0][2] == self.board[1][1] == self.board[2][0]:
            return self.board[0][2]
        
        return None
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return False
                if event.key == pygame.K_r:
                    self.reset()
            if event.type == pygame.MOUSEBUTTONDOWN and not self.game_over:
                x = (event.pos[0] - self.offset_x) // self.cell_size
                y = (event.pos[1] - self.offset_y) // self.cell_size
                if 0 <= x < 3 and 0 <= y < 3 and self.board[y][x] == '':
                    self.board[y][x] = self.current_player
                    self.move_count += 1
                    self.winner = self.check_winner()
                    if self.winner:
                        self.game_over = True
                    elif self.move_count == 9:
                        self.game_over = True
                    else:
                        self.current_player = 'O' if self.current_player == 'X' else 'X'
        return True
    
    def update(self):
        pass
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((30, 30, 42))
        
        for row in range(3):
            for col in range(3):
                rect = pygame.Rect(col * self.cell_size + 2, row * self.cell_size + 2,
                                   self.cell_size - 4, self.cell_size - 4)
                pygame.draw.rect(game_surface, (50, 50, 65), rect, border_radius=8)
                
                if self.board[row][col]:
                    color = (60, 200, 255) if self.board[row][col] == 'X' else (255, 200, 60)
                    text = self.font.render(self.board[row][col], True, color)
                    text_rect = text.get_rect(center=(col * self.cell_size + self.cell_size // 2,
                                                       row * self.cell_size + self.cell_size // 2))
                    game_surface.blit(text, text_rect)
        
        for i in range(1, 3):
            pygame.draw.line(game_surface, (80, 80, 100), 
                            (i * self.cell_size, 5), (i * self.cell_size, self.height - 5), 3)
            pygame.draw.line(game_surface, (80, 80, 100),
                            (5, i * self.cell_size), (self.width - 5, i * self.cell_size), 3)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        status_y = self.offset_y + self.height + 20
        if self.game_over:
            if self.winner:
                text = self.font_big.render(f"Player {self.winner} Wins!", True, (0, 255, 0))
            else:
                text = self.font_big.render("Draw!", True, (255, 255, 100))
            text_rect = text.get_rect(center=(self.screen.get_width() // 2, status_y))
            self.screen.blit(text, text_rect)
            text2 = self.font_small.render("Press R to restart", True, (200, 200, 200))
            text2_rect = text2.get_rect(center=(self.screen.get_width() // 2, status_y + 40))
            self.screen.blit(text2, text2_rect)
        else:
            text = self.font_big.render(f"Player {self.current_player}'s Turn", True, (255, 255, 255))
            text_rect = text.get_rect(center=(self.screen.get_width() // 2, status_y))
            self.screen.blit(text, text_rect)
        
        pygame.display.flip()
    
    def run(self):
        running = True
        while running:
            if not self.handle_events():
                break
            self.update()
            self.draw()
            self.clock.tick(60)