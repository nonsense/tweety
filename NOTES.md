TweetyService
=============

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Upgrade","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2", "Owner": "anton"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.InstallSnapshot","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.CreateCanary","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.DeleteCanary","params":[{"ReleaseName": "lotus-0"}],"id":68}' localhost:1337

Acquire and Relese
==================

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Acquire","params":[{"ReleaseName": "lotus-1", "Owner": "anton"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Release","params":[{"ReleaseName": "lotus-1", "Owner": "anton"}],"id":68}' localhost:1337

docker build -t nonsens3/lotus --target=lotus --build-arg=BUILDER_BASE=builder-git --build-arg=TAG=v1.2.3  -f Dockerfile.lotus .
docker build -t nonsens3/lotus:v1.2.2 --target=lotus --build-arg=BUILDER_BASE=builder-git --build-arg=TAG=v1.2.2  -f Dockerfile.lotus .

helm upgrade --namespace lotus --install lotus-0 ./lotus-fullnode-minimal
helm upgrade --namespace lotus --install lotus-0 ./lotus-fullnode-minimal --set image.tag=v1.2.3 --set daemonArgs=null
