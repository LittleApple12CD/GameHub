import pygame
import random

class FlappyGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 400, 600
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 40)
        self.font_small = pygame.font.Font(None, 26)
        self.reset()
    
    def reset(self):
        self.bird_y = self.height // 2
        self.bird_vy = 0
        self.bird_radius = 15
        self.gravity = 0.6
        self.jump_strength = -11
        self.pipes = []
        self.pipe_width = 60
        self.pipe_gap = 160
        self.pipe_speed = 3
        self.score = 0
        self.game_over = False
        self.pipe_timer = 0
        self.pipe_delay = 90
    
    def add_pipe(self):
        gap_y = random.randint(80, self.height - 80 - self.pipe_gap)
        self.pipes.append({
            'x': self.width,
            'gap_y': gap_y,
            'passed': False
        })
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return False
                if event.key == pygame.K_r:
                    self.reset()
                if event.key == pygame.K_SPACE and not self.game_over:
                    self.bird_vy = self.jump_strength
            if event.type == pygame.MOUSEBUTTONDOWN and not self.game_over:
                self.bird_vy = self.jump_strength
        return True
    
    def update(self):
        if self.game_over:
            return
        
        self.bird_vy += self.gravity
        self.bird_y += self.bird_vy
        
        if self.bird_y - self.bird_radius < 0:
            self.bird_y = self.bird_radius
            self.bird_vy = 0
        elif self.bird_y + self.bird_radius > self.height:
            self.game_over = True
            return
        
        self.pipe_timer += 1
        if self.pipe_timer >= self.pipe_delay:
            self.pipe_timer = 0
            self.add_pipe()
        
        for pipe in self.pipes[:]:
            pipe['x'] -= self.pipe_speed
            
            if not pipe['passed'] and pipe['x'] + self.pipe_width < self.bird_x:
                pipe['passed'] = True
                self.score += 1
            
            if (pipe['x'] < self.bird_x + self.bird_radius and 
                pipe['x'] + self.pipe_width > self.bird_x - self.bird_radius):
                if (self.bird_y - self.bird_radius < pipe['gap_y'] or 
                    self.bird_y + self.bird_radius > pipe['gap_y'] + self.pipe_gap):
                    self.game_over = True
                    return
            
            if pipe['x'] + self.pipe_width < 0:
                self.pipes.remove(pipe)
    
    @property
    def bird_x(self):
        return 80
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((135, 206, 235))
        
        pygame.draw.circle(game_surface, (255, 255, 50), 
                          (self.bird_x, int(self.bird_y)), self.bird_radius)
        pygame.draw.circle(game_surface, (0, 0, 0),
                          (self.bird_x + 8, int(self.bird_y) - 5), 4)
        pygame.draw.circle(game_surface, (255, 255, 255),
                          (self.bird_x + 10, int(self.bird_y) - 7), 2)
        
        for pipe in self.pipes:
            top_rect = pygame.Rect(pipe['x'], 0, self.pipe_width, pipe['gap_y'])
            pygame.draw.rect(game_surface, (40, 200, 40), top_rect, border_radius=6)
            pygame.draw.rect(game_surface, (30, 160, 30), 
                           (pipe['x'] - 6, pipe['gap_y'] - 22, self.pipe_width + 12, 22), border_radius=6)
            
            bottom_y = pipe['gap_y'] + self.pipe_gap
            bottom_rect = pygame.Rect(pipe['x'], bottom_y, self.pipe_width, self.height - bottom_y)
            pygame.draw.rect(game_surface, (40, 200, 40), bottom_rect, border_radius=6)
            pygame.draw.rect(game_surface, (30, 160, 30),
                           (pipe['x'] - 6, bottom_y, self.pipe_width + 12, 22), border_radius=6)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        score_text = self.font.render(str(self.score), True, (255, 255, 255))
        score_rect = score_text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + 50))
        self.screen.blit(score_text, score_rect)
        
        if self.game_over:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(150)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            text = self.font.render("GAME OVER", True, (255, 255, 255))
            text_rect = text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 - 30))
            self.screen.blit(text, text_rect)
            text2 = self.font_small.render("Press R to restart", True, (255, 255, 255))
            text2_rect = text2.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 + 20))
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