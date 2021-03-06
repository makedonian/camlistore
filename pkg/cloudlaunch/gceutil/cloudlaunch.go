/*
Copyright 2015 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package gceutil provides utility functions to help with instances on
// Google Compute Engine.
package gceutil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type coreOSImage struct {
	SelfLink          string
	CreationTimestamp time.Time
	Name              string
}

type coreOSImageList struct {
	Items []coreOSImage
}

// CoreOSImageURL returns the URL of the latest stable CoreOS image for running on Google Compute Engine.
func CoreOSImageURL(cl *http.Client) (string, error) {
	resp, err := cl.Get("https://www.googleapis.com/compute/v1/projects/coreos-cloud/global/images")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	imageList := &coreOSImageList{}
	if err := json.NewDecoder(resp.Body).Decode(imageList); err != nil {
		return "", err
	}
	if imageList == nil || len(imageList.Items) == 0 {
		return "", errors.New("no images list in response")
	}

	imageURL := ""
	var max time.Time // latest stable image creation time
	for _, v := range imageList.Items {
		if !strings.HasPrefix(v.Name, "coreos-stable") {
			continue
		}
		if v.CreationTimestamp.After(max) {
			max = v.CreationTimestamp
			imageURL = v.SelfLink
		}
	}
	if imageURL == "" {
		return "", errors.New("no stable coreOS image found")
	}
	return imageURL, nil
}
