import pygame
import random

class TetrisGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 300, 600
        self.offset_x = (screen.get_width() - self.width - 140) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.cols, self.rows = 10, 20
        self.cell_size = 30
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 32)
        self.font_big = pygame.font.Font(None, 48)
        
        self.shapes = {
            'I': [[0,0,0,0],[1,1,1,1],[0,0,0,0],[0,0,0,0]],
            'O': [[1,1],[1,1]],
            'T': [[0,1,0],[1,1,1],[0,0,0]],
            'S': [[0,1,1],[1,1,0],[0,0,0]],
            'Z': [[1,1,0],[0,1,1],[0,0,0]],
            'L': [[1,0,0],[1,1,1],[0,0,0]],
            'J': [[0,0,1],[1,1,1],[0,0,0]]
        }
        self.colors = {
            'I': (0, 240, 240),
            'O': (240, 240, 0),
            'T': (180, 0, 240),
            'S': (0, 240, 0),
            'Z': (240, 0, 0),
            'L': (240, 160, 0),
            'J': (0, 0, 240)
        }
        self.reset()
    
    def reset(self):
        self.board = [[0] * self.cols for _ in range(self.rows)]
        self.score = 0
        self.game_over = False
        self.fall_timer = 0
        self.fall_delay = 500
        self.spawn_piece()
    
    def spawn_piece(self):
        self.current_piece = random.choice(list(self.shapes.keys()))
        self.piece_shape = [row[:] for row in self.shapes[self.current_piece]]
        self.piece_x = self.cols // 2 - len(self.piece_shape[0]) // 2
        self.piece_y = 0
        
        if self.check_collision(self.piece_shape, self.piece_x, self.piece_y):
            self.game_over = True
    
    def check_collision(self, shape, x, y):
        for row_idx, row in enumerate(shape):
            for col_idx, cell in enumerate(row):
                if cell:
                    board_x = x + col_idx
                    board_y = y + row_idx
                    if (board_x < 0 or board_x >= self.cols or 
                        board_y >= self.rows or
                        (board_y >= 0 and self.board[board_y][board_x])):
                        return True
        return False
    
    def rotate_piece(self):
        shape = self.piece_shape
        rotated = [[shape[y][x] for y in range(len(shape))] for x in range(len(shape[0])-1, -1, -1)]
        if not self.check_collision(rotated, self.piece_x, self.piece_y):
            self.piece_shape = rotated
    
    def lock_piece(self):
        for row_idx, row in enumerate(self.piece_shape):
            for col_idx, cell in enumerate(row):
                if cell:
                    board_y = self.piece_y + row_idx
                    board_x = self.piece_x + col_idx
                    if board_y >= 0:
                        self.board[board_y][board_x] = self.current_piece
        self.clear_lines()
        self.spawn_piece()
    
    def clear_lines(self):
        lines_cleared = 0
        for row in range(self.rows - 1, -1, -1):
            if all(self.board[row]):
                del self.board[row]
                self.board.insert(0, [0] * self.cols)
                lines_cleared += 1
        if lines_cleared:
            self.score += [0, 100, 300, 500, 800][lines_cleared]
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return False
                if event.key == pygame.K_r:
                    self.reset()
                if not self.game_over:
                    if event.key == pygame.K_LEFT:
                        if not self.check_collision(self.piece_shape, self.piece_x - 1, self.piece_y):
                            self.piece_x -= 1
                    elif event.key == pygame.K_RIGHT:
                        if not self.check_collision(self.piece_shape, self.piece_x + 1, self.piece_y):
                            self.piece_x += 1
                    elif event.key == pygame.K_DOWN:
                        if not self.check_collision(self.piece_shape, self.piece_x, self.piece_y + 1):
                            self.piece_y += 1
                    elif event.key == pygame.K_SPACE:
                        self.rotate_piece()
                    elif event.key == pygame.K_UP:
                        self.rotate_piece()
        return True
    
    def update(self):
        if self.game_over:
            return
        self.fall_timer += self.clock.get_time()
        if self.fall_timer >= self.fall_delay:
            self.fall_timer = 0
            if not self.check_collision(self.piece_shape, self.piece_x, self.piece_y + 1):
                self.piece_y += 1
            else:
                self.lock_piece()
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((20, 20, 30))
        
        for y, row in enumerate(self.board):
            for x, cell in enumerate(row):
                if cell:
                    color = self.colors[cell]
                    rect = pygame.Rect(x * self.cell_size + 1, y * self.cell_size + 1,
                                       self.cell_size - 2, self.cell_size - 2)
                    pygame.draw.rect(game_surface, color, rect, border_radius=4)
        
        if not self.game_over:
            for y, row in enumerate(self.piece_shape):
                for x, cell in enumerate(row):
                    if cell:
                        color = self.colors[self.current_piece]
                        rect = pygame.Rect((self.piece_x + x) * self.cell_size + 1,
                                           (self.piece_y + y) * self.cell_size + 1,
                                           self.cell_size - 2, self.cell_size - 2)
                        pygame.draw.rect(game_surface, color, rect, border_radius=4)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        score_text = self.font.render(f"Score: {self.score}", True, (255, 255, 255))
        self.screen.blit(score_text, (self.offset_x + self.width + 20, self.offset_y + 20))
        
        if self.game_over:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(180)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            text = self.font_big.render("GAME OVER", True, (255, 255, 255))
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