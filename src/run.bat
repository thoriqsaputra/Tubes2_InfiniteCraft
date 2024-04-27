@echo off

rem Start the frontend in one command prompt window
start cmd /k npm run dev

rem Navigate to the backend folder
cd Backend

rem Start the backend in another command prompt window
start cmd /k go run back.go BFS.go IDS.go
