import pygame
import random

class SnakeGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 600, 600
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.cell_size = 20
        self.grid_width = self.width // self.cell_size
        self.grid_height = self.height // self.cell_size
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 36)
        self.reset()
    
    def reset(self):
        self.snake = [(self.grid_width // 2, self.grid_height // 2)]
        self.direction = (1, 0)
        self.next_direction = (1, 0)
        self.food = self.spawn_food()
        self.score = 0
        self.game_over = False
        self.move_timer = 0
        self.move_delay = 150
    
    def spawn_food(self):
        while True:
            pos = (random.randint(0, self.grid_width - 1), random.randint(0, self.grid_height - 1))
            if pos not in self.snake:
                return pos
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_UP and self.direction != (0, 1):
                    self.next_direction = (0, -1)
                elif event.key == pygame.K_DOWN and self.direction != (0, -1):
                    self.next_direction = (0, 1)
                elif event.key == pygame.K_LEFT and self.direction != (1, 0):
                    self.next_direction = (-1, 0)
                elif event.key == pygame.K_RIGHT and self.direction != (-1, 0):
                    self.next_direction = (1, 0)
                elif event.key == pygame.K_ESCAPE:
                    return False
                elif event.key == pygame.K_r:
                    self.reset()
        return True
    
    def update(self):
        if self.game_over:
            return
        
        self.move_timer += self.clock.get_time()
        if self.move_timer >= self.move_delay:
            self.move_timer = 0
            self.direction = self.next_direction
            head = self.snake[0]
            new_head = (head[0] + self.direction[0], head[1] + self.direction[1])
            
            if (new_head[0] < 0 or new_head[0] >= self.grid_width or
                new_head[1] < 0 or new_head[1] >= self.grid_height):
                self.game_over = True
                return
            
            if new_head in self.snake:
                self.game_over = True
                return
            
            self.snake.insert(0, new_head)
            
            if new_head == self.food:
                self.score += 10
                self.food = self.spawn_food()
            else:
                self.snake.pop()
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((20, 20, 30))
        
        for i, segment in enumerate(self.snake):
            color = (60, 220, 60) if i == 0 else (40, 180, 40)
            rect = pygame.Rect(segment[0] * self.cell_size + 1, segment[1] * self.cell_size + 1,
                               self.cell_size - 2, self.cell_size - 2)
            pygame.draw.rect(game_surface, color, rect, border_radius=4)
        
        food_rect = pygame.Rect(self.food[0] * self.cell_size + 1, self.food[1] * self.cell_size + 1,
                                self.cell_size - 2, self.cell_size - 2)
        pygame.draw.rect(game_surface, (255, 60, 60), food_rect, border_radius=6)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        score_text = self.font.render(f"Score: {self.score}", True, (255, 255, 255))
        self.screen.blit(score_text, (self.offset_x + 10, self.offset_y + 10))
        
        if self.game_over:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(180)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            text = self.font.render("GAME OVER - Press R", True, (255, 255, 255))
            text_rect = text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2))
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