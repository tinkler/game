PWD=pwd
cd ..
docker build -f build/Dockerfile -t game-server .
cd ${PWD}