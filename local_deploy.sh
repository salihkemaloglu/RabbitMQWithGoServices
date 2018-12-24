#!/bin/bash
BLUE='\033[1;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${BLUE}Open client folder${NC}"
cd client

echo -e "${GREEN}Build client  project...${NC}"
docker build -t clientgo .

echo -e "${BLUE}Go back main folder${NC}"
cd ..

echo -e "${BLUE}Open service folder${NC}"
cd service

echo -e "${GREEN}Build service testing project...${NC}"
docker build -t servicego .

echo -e "${BLUE}Go back main folder${NC}"
cd ..

echo -e "${GREEN}Services are starting...${NC}"
docker-compose up