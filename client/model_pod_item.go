/*
 * Pipeline API
 *
 * Pipeline v0.3.0 swagger
 *
 * API version: 0.3.0
 * Contact: info@banzaicloud.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

type PodItem struct {
	Name            string                 `json:"name,omitempty"`
	Namespace       string                 `json:"namespace,omitempty"`
	CreatedAt       string                 `json:"createdAt,omitempty"`
	Labels          PodItemLabels          `json:"labels,omitempty"`
	RestartPolicy   string                 `json:"restartPolicy,omitempty"`
	Conditions      []PodCondition         `json:"conditions,omitempty"`
	ResourceSummary PodItemResourceSummary `json:"resourceSummary,omitempty"`
}