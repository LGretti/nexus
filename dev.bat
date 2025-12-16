@echo off
echo ==========================================
echo      INICIANDO O SISTEMA NEXUS 🚀
echo ==========================================

:: 1. Abre uma nova janela, entra na pasta 'api' e roda o Go
start "Nexus Backend (Go)" cmd /k "cd api && go run cmd/api/main.go"

:: 2. Abre uma nova janela, entra na pasta 'frontend' e roda o NPM
start "Nexus Frontend (Next.js)" cmd /k "cd frontend && npm run dev"

echo.
echo Backend e Frontend iniciados em janelas separadas.
echo Pode minimizar esta aqui.
echo.