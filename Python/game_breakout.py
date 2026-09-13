import pygame
import random

class BreakoutGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 600, 600
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 36)
        self.font_big = pygame.font.Font(None, 48)
        self.reset()
    
    def reset(self):
        self.paddle = pygame.Rect(self.width // 2 - 60, self.height - 40, 120, 16)
        self.paddle_speed = 8
        self.ball = pygame.Rect(self.width // 2 - 10, self.height - 70, 20, 20)
        self.ball_dx = 5
        self.ball_dy = -6
        self.ball_speed = 6
        self.bricks = []
        self.score = 0
        self.lives = 3
        self.game_over = False
        self.waiting = True
        
        rows, cols = 5, 8
        brick_width = (self.width - 20) // cols - 4
        brick_height = 22
        colors = [(230, 50, 50), (230, 150, 50), (230, 230, 50), 
                  (50, 230, 50), (50, 150, 230)]
        for row in range(rows):
            for col in range(cols):
                x = 10 + col * (brick_width + 4)
                y = 40 + row * (brick_height + 4)
                self.bricks.append({
                    'rect': pygame.Rect(x, y, brick_width, brick_height),
                    'color': colors[row % len(colors)],
                    'alive': True
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
                if event.key == pygame.K_SPACE and self.waiting and not self.game_over:
                    self.waiting = False
        
        if not self.game_over and not self.waiting:
            keys = pygame.key.get_pressed()
            if keys[pygame.K_LEFT] and self.paddle.left > 0:
                self.paddle.x -= self.paddle_speed
            if keys[pygame.K_RIGHT] and self.paddle.right < self.width:
                self.paddle.x += self.paddle_speed
        return True
    
    def update(self):
        if self.game_over or self.waiting:
            return
        
        self.ball.x += self.ball_dx
        self.ball.y += self.ball_dy
        
        if self.ball.left <= 0 or self.ball.right >= self.width:
            self.ball_dx = -self.ball_dx
        if self.ball.top <= 0:
            self.ball_dy = -self.ball_dy
        
        if self.ball.bottom >= self.height:
            self.lives -= 1
            if self.lives <= 0:
                self.game_over = True
            else:
                self.waiting = True
                self.ball.x = self.width // 2 - 10
                self.ball.y = self.height - 70
                self.ball_dx = self.ball_speed * (1 if random.choice([True, False]) else -1)
                self.ball_dy = -self.ball_speed
            return
        
        if self.ball.colliderect(self.paddle):
            self.ball_dy = -abs(self.ball_dy)
            hit_pos = (self.ball.centerx - self.paddle.centerx) / (self.paddle.width / 2)
            self.ball_dx = hit_pos * self.ball_speed * 0.9
            if abs(self.ball_dx) < 1.5:
                self.ball_dx = 3 if self.ball_dx >= 0 else -3
        
        for brick in self.bricks:
            if not brick['alive']:
                continue
            if self.ball.colliderect(brick['rect']):
                brick['alive'] = False
                self.score += 10
                
                overlap_top = self.ball.bottom - brick['rect'].top
                overlap_bottom = brick['rect'].bottom - self.ball.top
                overlap_left = self.ball.right - brick['rect'].left
                overlap_right = brick['rect'].right - self.ball.left
                min_overlap = min(overlap_top, overlap_bottom, overlap_left, overlap_right)
                
                if min_overlap == overlap_top or min_overlap == overlap_bottom:
                    self.ball_dy = -self.ball_dy
                else:
                    self.ball_dx = -self.ball_dx
                break
        
        if all(not brick['alive'] for brick in self.bricks):
            self.game_over = True
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((20, 20, 30))
        
        for brick in self.bricks:
            if brick['alive']:
                pygame.draw.rect(game_surface, brick['color'], brick['rect'], border_radius=8)
                pygame.draw.rect(game_surface, (255, 255, 255), brick['rect'], 1, border_radius=8)
        
        pygame.draw.rect(game_surface, (255, 255, 255), self.paddle, border_radius=10)
        pygame.draw.rect(game_surface, (200, 200, 220), self.paddle, 2, border_radius=10)
        
        pygame.draw.ellipse(game_surface, (255, 255, 120), self.ball)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        score_text = self.font.render(f"Score: {self.score}", True, (255, 255, 255))
        self.screen.blit(score_text, (self.offset_x + 10, self.offset_y + 10))
        lives_text = self.font.render(f"Lives: {self.lives}", True, (255, 255, 255))
        self.screen.blit(lives_text, (self.offset_x + self.width - 120, self.offset_y + 10))
        
        if self.waiting and not self.game_over:
            text = self.font_big.render("Press SPACE", True, (255, 255, 255))
            text_rect = text.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 + 30))
            self.screen.blit(text, text_rect)
            text2 = self.font.render("to start", True, (255, 255, 255))
            text2_rect = text2.get_rect(center=(self.offset_x + self.width // 2, self.offset_y + self.height // 2 + 70))
            self.screen.blit(text2, text2_rect)
        elif self.game_over:
            overlay = pygame.Surface((self.width, self.height))
            overlay.set_alpha(150)
            overlay.fill((0, 0, 0))
            self.screen.blit(overlay, (self.offset_x, self.offset_y))
            if self.lives <= 0:
                text = self.font_big.render("GAME OVER", True, (255, 255, 255))
            else:
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