version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: aniverse
      POSTGRES_USER: aniverse
      POSTGRES_PASSWORD: aniverse123
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: aniverse
      DB_PASSWORD: aniverse123
      DB_NAME: aniverse
      REDIS_HOST: redis
      REDIS_PORT: 6379
      JWT_SECRET: aniverse-secret-key-2024
    depends_on:
      - postgres
      - redis

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080
    depends_on:
      - backend

volumes:
  postgres_data: