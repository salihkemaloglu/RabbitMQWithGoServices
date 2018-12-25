#!/bin/bash
BLUE='\033[1;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${BLUE}Open sender folder${NC}"
cd sender

echo -e "${GREEN}Build sender  project...${NC}"
docker build -t sendergo .

echo -e "${BLUE}Go back main folder${NC}"
cd ..

echo -e "${BLUE}Open consumer folder${NC}"
cd consumer

echo -e "${GREEN}Build consumer testing project...${NC}"
docker build -t consumergo .

echo -e "${BLUE}Go back main folder${NC}"
cd ..

echo -e "${GREEN}Services are starting...${NC}"
docker-compose up