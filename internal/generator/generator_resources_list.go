// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package generator

import "time"

type ResourcesListGenerator struct {
	GeneratedTimestamp time.Time
	Resources          []ResourceConfig
	Packages           []string
}

func (p ResourcesListGenerator) String() string {
	return renderTemplate(resourcesListTemplate, p)
}
