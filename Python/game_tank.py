import pygame
import random

class TankGame:
    def __init__(self, screen):
        self.screen = screen
        self.width, self.height = 600, 600
        self.offset_x = (screen.get_width() - self.width) // 2
        self.offset_y = (screen.get_height() - self.height) // 2
        self.cell_size = 30
        self.grid_size = self.width // self.cell_size
        self.bullet_speed = 1
        self.clock = pygame.time.Clock()
        self.font = pygame.font.Font(None, 36)
        self.font_big = pygame.font.Font(None, 48)
        self.reset()
    
    def reset(self):
        self.walls = self.generate_walls()
        self.player = {'x': 1, 'y': 1, 'dir': (1, 0), 'alive': True, 'move_timer': 0}
        self.enemies = [{'x': self.grid_size - 2, 'y': self.grid_size - 2, 
                        'dir': (-1, 0), 'alive': True, 'timer': 0, 'move_timer': 0}]
        self.bullets = []
        self.score = 0
        self.game_over = False
        self.enemy_spawn_timer = 0
        self.enemy_spawn_delay = 120
        self.move_delay = 200
        self.enemy_move_delay = 250
    
    def generate_walls(self):
        walls = []
        for i in range(self.grid_size):
            walls.append((i, 0))
            walls.append((i, self.grid_size - 1))
            walls.append((0, i))
            walls.append((self.grid_size - 1, i))
        for _ in range(12):
            x = random.randint(2, self.grid_size - 3)
            y = random.randint(2, self.grid_size - 3)
            if (x, y) != (1, 1) and (x, y) != (self.grid_size - 2, self.grid_size - 2):
                walls.append((x, y))
        return list(set(walls))
    
    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                return False
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_ESCAPE:
                    return False
                if event.key == pygame.K_r:
                    self.reset()
                if event.key == pygame.K_SPACE and not self.game_over and self.player['alive']:
                    self.shoot_bullet(self.player['x'], self.player['y'], 
                                     self.player['dir'], True)
        
        if not self.game_over and self.player['alive']:
            keys = pygame.key.get_pressed()
            dx, dy = 0, 0
            if keys[pygame.K_w]:
                dy = -1
            elif keys[pygame.K_s]:
                dy = 1
            elif keys[pygame.K_a]:
                dx = -1
            elif keys[pygame.K_d]:
                dx = 1
            
            if dx or dy:
                self.player['dir'] = (dx, dy)
                self.player['move_timer'] += self.clock.get_time()
                if self.player['move_timer'] >= self.move_delay:
                    self.player['move_timer'] = 0
                    new_x = self.player['x'] + dx
                    new_y = self.player['y'] + dy
                    if self.can_move(new_x, new_y, self.player['dir'], True):
                        self.player['x'] = new_x
                        self.player['y'] = new_y
            else:
                self.player['move_timer'] = 0
        
        return True
    
    def can_move(self, x, y, direction, is_player):
        if (x, y) in self.walls:
            return False
        
        if is_player:
            for enemy in self.enemies:
                if enemy['alive'] and enemy['x'] == x and enemy['y'] == y:
                    return False
        else:
            if self.player['alive'] and self.player['x'] == x and self.player['y'] == y:
                return False
            for enemy in self.enemies:
                if enemy is not self and enemy['alive'] and enemy['x'] == x and enemy['y'] == y:
                    return False
        
        return 0 <= x < self.grid_size and 0 <= y < self.grid_size
    
    def shoot_bullet(self, x, y, direction, is_player):
        self.bullets.append({
            'x': x + direction[0], 'y': y + direction[1],
            'dir': direction,
            'is_player': is_player,
            'alive': True
        })
    
    def update(self):
        if self.game_over:
            return
        
        for bullet in self.bullets[:]:
            if not bullet['alive']:
                continue
        
            dx = bullet['dir'][0] * self.bullet_speed
            dy = bullet['dir'][1] * self.bullet_speed
        
            steps = int(max(abs(dx), abs(dy))) + 1
            step_x = dx / steps
            step_y = dy / steps
        
            hit = False
            for _ in range(steps):
                bullet['x'] += step_x
                bullet['y'] += step_y
            
                bx = int(round(bullet['x']))
                by = int(round(bullet['y']))
            
                if not (0 <= bx < self.grid_size and 0 <= by < self.grid_size):
                    bullet['alive'] = False
                    hit = True
                    break
            
                if (bx, by) in self.walls:
                    bullet['alive'] = False
                    hit = True
                    break
            
                if bullet['is_player']:
                    for enemy in self.enemies:
                        if enemy['alive'] and enemy['x'] == bx and enemy['y'] == by:
                            enemy['alive'] = False
                            bullet['alive'] = False
                            self.score += 10
                            hit = True
                            break
                    if hit:
                        break
                else:
                    if self.player['alive'] and self.player['x'] == bx and self.player['y'] == by:
                        self.player['alive'] = False
                        bullet['alive'] = False
                        self.game_over = True
                        hit = True
                        break
        
            if hit:
                continue
    
        self.bullets = [b for b in self.bullets if b['alive']]
        
        for enemy in self.enemies:
            if not enemy['alive']:
                continue
            
            enemy['timer'] += 1
            if enemy['timer'] >= 20:
                enemy['timer'] = 0
                if random.random() < 0.25:
                    enemy['dir'] = random.choice([(1,0), (-1,0), (0,1), (0,-1)])
                
                enemy['move_timer'] += 20
                if enemy['move_timer'] >= self.enemy_move_delay:
                    enemy['move_timer'] = 0
                    new_x = enemy['x'] + enemy['dir'][0]
                    new_y = enemy['y'] + enemy['dir'][1]
                    if self.can_move(new_x, new_y, enemy['dir'], False):
                        enemy['x'] = new_x
                        enemy['y'] = new_y
            
            if random.random() < 0.015:
                self.shoot_bullet(enemy['x'], enemy['y'], enemy['dir'], False)
        
        self.enemy_spawn_timer += 1
        if self.enemy_spawn_timer >= self.enemy_spawn_delay:
            self.enemy_spawn_timer = 0
            alive_count = sum(1 for e in self.enemies if e['alive'])
            if alive_count < 3:
                self.enemies.append({
                    'x': self.grid_size - 2, 'y': self.grid_size - 2,
                    'dir': (-1, 0), 'alive': True, 'timer': 0, 'move_timer': 0
                })
    
    def draw_tank(self, surface, x, y, color, gun_color, direction):
        cx = x * self.cell_size + 15
        cy = y * self.cell_size + 15
        rect = pygame.Rect(x * self.cell_size + 2, y * self.cell_size + 2, 26, 26)
        pygame.draw.rect(surface, color, rect, border_radius=8)
        dx, dy = direction
        if dx == 1:
            gun_rect = pygame.Rect(cx + 4, cy - 4, 14, 8)
        elif dx == -1:
            gun_rect = pygame.Rect(cx - 18, cy - 4, 14, 8)
        elif dy == -1:
            gun_rect = pygame.Rect(cx - 4, cy - 18, 8, 14)
        else:
            gun_rect = pygame.Rect(cx - 4, cy + 4, 8, 14)
        pygame.draw.rect(surface, gun_color, gun_rect, border_radius=4)
        pygame.draw.circle(surface, color, (cx, cy), 6)
    
    def draw(self):
        self.screen.fill((25, 25, 35))
        
        game_surface = pygame.Surface((self.width, self.height))
        game_surface.fill((30, 30, 42))
        
        for x, y in self.walls:
            rect = pygame.Rect(x * self.cell_size + 1, y * self.cell_size + 1,
                               self.cell_size - 2, self.cell_size - 2)
            pygame.draw.rect(game_surface, (100, 100, 130), rect, border_radius=8)
        
        if self.player['alive']:
            self.draw_tank(game_surface, self.player['x'], self.player['y'],
                          (50, 230, 50), (30, 200, 30), self.player['dir'])
        
        for enemy in self.enemies:
            if enemy['alive']:
                self.draw_tank(game_surface, enemy['x'], enemy['y'],
                              (230, 50, 50), (200, 30, 30), enemy['dir'])
        
        for bullet in self.bullets:
            x, y = bullet['x'], bullet['y']
            color = (255, 255, 80) if bullet['is_player'] else (255, 150, 50)
            rect = pygame.Rect(x * self.cell_size + 10, y * self.cell_size + 10, 10, 10)
            pygame.draw.rect(game_surface, color, rect, border_radius=5)
        
        self.screen.blit(game_surface, (self.offset_x, self.offset_y))
        
        score_text = self.font.render(f"Score: {self.score}", True, (255, 255, 255))
        self.screen.blit(score_text, (self.offset_x + 10, self.offset_y + 10))
        
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
        
        pygame.display.flip()
    
    def run(self):
        running = True
        while running:
            if not self.handle_events():
                break
            self.update()
            self.draw()
            self.clock.tick(60)