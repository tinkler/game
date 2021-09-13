PWD=pwd
cd ../cmd/server
CGO_ENABLED=0 go build -tags netgo -a -v
cd ../../
docker build -f build/alpine/Dockerfile -t game-server .
rm cmd/server/server
cd ${PWD}