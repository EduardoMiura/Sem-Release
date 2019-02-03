package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/poc-git/helpers"
	"github.com/poc-git/semrelease"
	"go.uber.org/zap"
)

type handlerFuncError func(w http.ResponseWriter, r *http.Request) error

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("V1.0"))
}

func addConfigHandler(svc *semrelease.Service) handlerFuncError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		infile, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error parsing uploaded file: "+err.Error(), http.StatusBadRequest)
			return err
		}

		// THIS IS VERY INSECURE! DO NOT DO THIS!
		outfile, err := os.Create("./config/" + header.Filename)
		if err != nil {
			http.Error(w, "Error saving file: "+err.Error(), http.StatusBadRequest)
			return err
		}

		_, err = io.Copy(outfile, infile)
		if err != nil {
			http.Error(w, "Error saving file: "+err.Error(), http.StatusBadRequest)
			return err
		}
		if err := svc.CheckHealth(ctx, time.Second); err != nil { // TODO: timeout on property
			return err
		}
		return responseWriter(ctx, w, http.StatusOK, nil)
	}
}
func getAllRepositoryHandler(svc *semrelease.Service) handlerFuncError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		rep, _ := svc.GetRepositories()
		return responseWriter(ctx, w, http.StatusOK, rep)

	}
}

func getAllConfigHandler(svc *semrelease.Service) handlerFuncError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		if err := svc.CheckHealth(ctx, time.Second); err != nil { // TODO: timeout on property
			return err
		}
		return responseWriter(ctx, w, http.StatusOK, nil)
	}
}

func healthCheckHandler(svc *semrelease.Service) handlerFuncError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		if err := svc.CheckHealth(ctx, time.Second); err != nil { // TODO: timeout on property
			return err
		}
		return responseWriter(ctx, w, http.StatusOK, nil)
	}
}

func webhookHandler(svc *semrelease.Service) handlerFuncError {
	return func(w http.ResponseWriter, r *http.Request) error {
		//sd := svc.ReturnClient()
		body, _ := ioutil.ReadAll(r.Body)
		event := semrelease.Event{}
		//m := map[string]interface{}{}
		json.Unmarshal(body, &event)
		repositoryname := event.Repository.Name
		login := event.Repository.Owner.Login

		fmt.Println(event.PullRequest, " --------------------------------------dssssssssssssssssssssssssssssssssssssssssssss--------------------------------------")

		fmt.Println(repositoryname, login)

		if event.PullRequest != nil && event.PullRequest.Merged {
			svc.CreateRelease(login, repositoryname)
		}

		// for k, v := range m {
		// 	if k == "action" {
		// 		if v == "closed" {
		// 			count++
		// 		}
		// 	}
		// 	if k == "pull_request" {
		// 		msg := v.(map[string]interface{})
		// 		for key, value := range msg {
		// 			if key == "merged" {
		// 				merged := value.(bool)
		// 				if merged {
		// 					count++
		// 				}
		// 			}
		// 		}
		// 	}
		// 	if count > 1 {
		// 		///chamar o release
		// 		createRelease(login, repositoryname)
		// 		break
		// 	}
		// }
		return nil

	}
}

func errorWrapper(fn handlerFuncError, logger *zap.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			helpers.ErrorWriter(r.Context(), w, err, logger)
		}
	}
}

func responseWriter(ctx context.Context, w http.ResponseWriter, code int, content interface{}) error {
	if content == nil {
		w.WriteHeader(code)
		return nil
	}

	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	contentType := "application/json"
	// if content.Version() != "" {
	// 	contentType = fmt.Sprintf("application/%s+json", content.Version())
	// }

	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(code)
	w.Write(jsonContent)

	return nil
}
