package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("tweety")

// Kubernetes namespace where Canaries are to be deployed
// and that Tweety would control
var namespace = "lotus"

type TweetyService struct{}

func NewTweetyService() *TweetyService {
	return &TweetyService{}
}

type UpgradeResponse struct {
	Result string
}

type UpgradeRequest struct {
	ImageTag    string
	ReleaseName string
	Owner       string
}

func (t *TweetyService) Upgrade(r *http.Request, req *UpgradeRequest, result *UpgradeResponse) error {
	log.Debugw("handle.upgrade", "msg", "handling upgrade", "image_tag", req.ImageTag, "release_name", req.ReleaseName, "owner", req.Owner)

	defer func(now time.Time) {
		log.Debugw("handle.upgrade", "took", fmt.Sprintf("%s", time.Since(now)))
	}(time.Now())

	if err := validateRequest(req); err != nil {
		log.Warnw("handle.upgrade", "request-not-valid", err)
		return err
	}

	// helm upgrade, but keep datadir, assuming node is already sycned up
	args := fmt.Sprintf("upgrade --namespace %s --install %s ./lotus-fullnode-minimal --set image.tag=%s --set daemonArgs=null", namespace, req.ReleaseName, req.ImageTag)
	log.Debugw("handle.upgrade", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorw("handle.upgrade", "error while running helm", err)
		return err
	}

	log.Debugw("handle.upgrade", "helm ran successfully", strings.Replace(out.String(), "\n", "; ", -1))
	return nil
}

func (t *TweetyService) InstallSnapshot(r *http.Request, req *UpgradeRequest, result *UpgradeResponse) error {
	log.Debugw("handle.install-snapshot", "msg", "handling install-snapshot", "image_tag", req.ImageTag, "release_name", req.ReleaseName, "owner", req.Owner)

	defer func(now time.Time) {
		log.Debugw("handle.install-snapshot", "took", fmt.Sprintf("%s", time.Since(now)))
	}(time.Now())

	if err := validateRequest(req); err != nil {
		log.Warnw("handle.install-snapshot", "request-not-valid", err)
		return err
	}

	var out bytes.Buffer

	// helm uninstall
	args := fmt.Sprintf("--namespace %s uninstall %s", namespace, req.ReleaseName)
	log.Debugw("handle.install-snapshot", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorw("handle.install-snapshot", "error while running helm", err)
		return err
	}

	// kubectl delete pvc
	args = fmt.Sprintf("--namespace %s delete pvc -l release=%s", namespace, req.ReleaseName)
	log.Debugw("handle.install-snapshot", "cmd", "kubectl", "args", args)
	cmd = exec.Command("kubectl", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Errorw("handle.install-snapshot", "error while running helm", err)
		return err
	}

	// helm install-snapshot, and sync from a snapshot from s3
	args = fmt.Sprintf("--namespace %s upgrade --install %s ./lotus-fullnode-minimal --set image.tag=%s", namespace, req.ReleaseName, req.ImageTag)
	log.Debugw("handle.install-snapshot", "cmd", "helm", "args", args)
	cmd = exec.Command("helm", strings.Split(args, " ")...)

	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Errorw("handle.install-snapshot", "error while running helm", err)
		return err
	}

	log.Debugw("handle.install-snapshot", "helm ran successfully", strings.Replace(out.String(), "\n", "; ", -1))
	return nil
}

func (t *TweetyService) CreateCanary(r *http.Request, req *UpgradeRequest, result *UpgradeResponse) error {
	log.Debugw("handle.create-canary", "msg", "handling create-canary", "image_tag", req.ImageTag, "release_name", req.ReleaseName, "owner", req.Owner)

	defer func(now time.Time) {
		log.Debugw("handle.create-canary", "took", fmt.Sprintf("%s", time.Since(now)))
	}(time.Now())

	if err := validateRequest(req); err != nil {
		log.Warnw("handle.create-canary", "request-not-valid", err)
		return err
	}

	if exists(req.ReleaseName) {
		return fmt.Errorf("canary %s already exists", req.ReleaseName)
	}

	var out bytes.Buffer

	// helm install-snapshot, and sync from a snapshot from s3
	args := fmt.Sprintf("--namespace %s upgrade --install %s ./lotus-fullnode-minimal --set image.tag=%s", namespace, req.ReleaseName, req.ImageTag)
	log.Debugw("handle.create-canary", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorw("handle.create-canary", "error while running helm", err)
		return err
	}

	log.Debugw("handle.create-canary", "helm ran successfully", strings.Replace(out.String(), "\n", "; ", -1))
	return nil
}

func (t *TweetyService) DeleteCanary(r *http.Request, req *UpgradeRequest, result *UpgradeResponse) error {
	log.Debugw("handle.delete-canary", "msg", "handling delete-canary", "release_name", req.ReleaseName, "owner", req.Owner)

	defer func(now time.Time) {
		log.Debugw("handle.delete-canary", "took", fmt.Sprintf("%s", time.Since(now)))
	}(time.Now())

	if err := validateRequest(req); err != nil {
		log.Warnw("handle.delete-canary", "request-not-valid", err)
		return err
	}

	if !exists(req.ReleaseName) {
		return fmt.Errorf("canary %s doesn't exists", req.ReleaseName)
	}

	var out bytes.Buffer

	// helm uninstall
	args := fmt.Sprintf("--namespace %s uninstall %s", namespace, req.ReleaseName)
	log.Debugw("handle.delete-canary", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorw("handle.delete-canary", "error while running helm", err)
		return err
	}

	// kubectl delete pvc
	args = fmt.Sprintf("--namespace %s delete pvc -l release=%s", namespace, req.ReleaseName)
	log.Debugw("handle.delete-canary", "cmd", "kubectl", "args", args)
	cmd = exec.Command("helm", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Errorw("handle.delete-canary", "error while running helm", err)
		return err
	}

	log.Debugw("handle.delete-canary", "helm ran successfully", strings.Replace(out.String(), "\n", "; ", -1))
	return nil
}

func exists(releaseName string) bool {
	var out bytes.Buffer

	args := fmt.Sprintf("--namespace %s status %s", namespace, releaseName)
	log.Debugw("handle.exists", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

func validateRequest(req *UpgradeRequest) error {
	if req.ReleaseName == "" {
		return errors.New("missing ReleaseName param")
	}

	if req.Owner == "" {
		return errors.New("missing Owner param")
	}

	currentOwner, ok := claims.Has(req.ReleaseName)

	if currentOwner != req.Owner {
		if !ok {
			return errors.New("canary is not claimed by anyone; you must first claim canary")
		}

		return fmt.Errorf("canary is not claimed by owner; current owner is %s", currentOwner)
	}

	return nil
}
