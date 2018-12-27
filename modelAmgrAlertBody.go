/* 2018-12-25 (cc) <paul4hough@gmail.com>
   Prometheus AlertManager Alerts Body
*/

package main

import (
	"time"
)

type AmgrAlert struct {

	Annotations map[string]string `json:"annotations,omitempty"`

	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt,omitempty"`

	GeneratorURL string `json:"generatorURL"`

	Labels map[string]string `json:"labels"`

	Status string `json:"status"`

	RemedOut string
}

type AmgrAlertBody struct {

	Alerts []AmgrAlert `json:"alerts"`

	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty"`

	CommonLabels map[string]string `json:"commonLabels,omitempty"`

	ExternalURL string `json:"externalURL"`

	GroupKey string `json:"groupKey"`

	GroupLabels map[string]string `json:"groupLabels,omitempty"`

	Receiver string `json:"receiver"`

	Status string `json:"status"`

	Version string `json:"version"`
}
