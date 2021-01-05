# tweety

`tweety` is a service controlling all Lotus canary nodes. For more information, check https://docs.google.com/document/d/1ohYQvRpwwojjcaX35MucPesxGWHhmOWvUM2rNLf9Ubk/edit

** Note that this repository is currently in heavy work-in-progress. **

## Responsibilities

1. Keeps state of which canaries are free and available for testing.

2. Performs actions (listed below) such as upgrade, install-snapshot, release, etc.

## Actions

- [x] @lotusbot upgrade <GIT_REF> <HELM_DEPLOYMENT_REF> <OWNER> - (thin wrapper around helm) - should do a clean install of given git commit on given helm deployment (as we have one per VM), and keep the existing datadir (i.e. if we upgrade Lotus, we continue to sync mainnet from where previous version left off)

- [x] @lotusbot install-snapshot <GIT_REF> <HELM_DEPLOYMENT_REF> <OWNER> - (thin wrapper around helm) - should do a clean install of given git commit on given helm deployment (as we have one per k8s node), and remove datadir and start sync from snapshot as in https://docs.filecoin.io/get-started/lotus/installation/#start-the-lotus-daemon-and-sync-the-chain

- [ ] free <HELM_DEPLOYMENT_REF>

- developers should be able to free a given helm installation when they are done using it.
- acquiring is implicit when a developer calls upgrade, install-snapshot
- developers should not be able to deploy to a busy deployment

## Github and Slack integration

TODO

Probably the simplest integration between Github and the Canary bot is with Github Actions.

## Logs

Logs should be emitted with filebeat to https://logz.io

## Metrics

- [ ] All Lotus nodes emit metrics that are visualised in Grafana
- [ ] On `install-snapshot` and `upgrade` the Canary bot should send an annotation, so that it is clear when a specific deployment was updated (or reset to sync from scratch).

## Access to developers

Full access via AWS IAM and kubeconfig, since this will be hosted on `mainnet-us-east-2-dev-eks` EKS cluster.

Developers have `kubectl` and `helm` locally in case they want to attach to a node and run lotus CLI commands for example.

## How to use

```
curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Upgrade","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2", "Owner": "anton"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.InstallSnapshot","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.CreateCanary","params":[{"ReleaseName": "lotus-0", "ImageTag": "v1.2.2"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.DeleteCanary","params":[{"ReleaseName": "lotus-0"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Acquire","params":[{"ReleaseName": "lotus-1", "Owner": "anton"}],"id":68}' localhost:1337

curl -X POST -H "Content-Type: application/json" --data \
'{"jsonrpc":"2.0","method":"TweetyService.Release","params":[{"ReleaseName": "lotus-1", "Owner": "anton"}],"id":68}' localhost:1337
```
