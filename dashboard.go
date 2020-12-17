package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type ClaimPair struct {
	ReleaseName string
	Owner       string
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugw("dashboard", "msg", "handling dashboard/list")

	defer func(now time.Time) {
		log.Debugw("dashboard", "took", fmt.Sprintf("%s", time.Since(now)))
	}(time.Now())

	var out bytes.Buffer

	// helm ls
	args := fmt.Sprintf("--namespace %s ls", namespace)
	log.Debugw("dashboard", "cmd", "helm", "args", args)
	cmd := exec.Command("helm", strings.Split(args, " ")...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorw("dashboard", "error while running helm", err)
		return
	}

	t := template.New("dashboard.html")
	t, err = t.ParseFiles("tmpl/dashboard.html")
	if err != nil {
		panic(fmt.Sprintf("cannot ParseFiles with tmpl/dashboard: %s", err))
	}

	output := strings.Replace(out.String(), "\n", "<br>", -1)

	data := struct {
		HelmOutput string
		Claims     []ClaimPair
	}{
		output,
		nil,
	}

	for k, v := range claims.List() {
		data.Claims = append(data.Claims, ClaimPair{
			ReleaseName: k,
			Owner:       v.Owner,
		})
	}

	err = t.Execute(w, data)
	if err != nil {
		panic(fmt.Sprintf("cannot execute template: %s", err))
	}
}
