// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package autodiscovery

import (
	"fmt"
	"sync"

	"github.com/DataDog/datadog-agent/pkg/collector/check"
)

// TemplateCache is a data structure to store configuration templates
type TemplateCache struct {
	id2digests      map[string][]string     // map an AD identifier to the all the configs that have it
	digest2ids      map[string][]string     // map a config to the list of AD identifiers it has
	digest2template map[string]check.Config // map a digest to the corresponding config object
	m               sync.RWMutex
}

// NewTemplateCache creates a new cache
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		id2digests:      make(map[string][]string, 0),
		digest2ids:      make(map[string][]string, 0),
		digest2template: make(map[string]check.Config, 0),
	}
}

// Set stores or updates a template in the cache
func (cache *TemplateCache) Set(tpl check.Config) error {
	// return an error if configuration has no AD identifiers
	if len(tpl.ADIdentifiers) == 0 {
		return fmt.Errorf("template has no AD identifiers, unable to store it in the cache")
	}

	cache.m.Lock()
	defer cache.m.Unlock()

	// compute the template digest once
	d := tpl.Digest()

	// do nothing if the template is already in cache
	if _, found := cache.digest2ids[d]; found {
		return nil
	}

	// store the template
	cache.digest2template[d] = tpl
	cache.digest2ids[d] = tpl.ADIdentifiers
	for _, id := range tpl.ADIdentifiers {
		cache.id2digests[id] = append(cache.id2digests[id], d)
	}

	return nil
}

// Get retrieves a template from the cache
func (cache *TemplateCache) Get(adID string) ([]check.Config, error) {
	cache.m.RLock()
	defer cache.m.RUnlock()

	if digests, found := cache.id2digests[adID]; found {
		templates := make([]check.Config, len(digests))
		for _, digest := range digests {
			templates = append(templates, cache.digest2template[digest])
		}
		return templates, nil
	}

	return nil, fmt.Errorf("Autodiscovery id not found in cache")
}

// Del removes a template from the cache
func (cache *TemplateCache) Del(tpl check.Config) error {
	// compute the digest once
	d := tpl.Digest()

	cache.m.Lock()
	defer cache.m.Unlock()

	// returns an error in case the template isn't there
	if _, found := cache.digest2ids[d]; !found {
		return fmt.Errorf("template not found in cache")
	}

	// remove the template
	delete(cache.digest2ids, d)
	delete(cache.digest2template, d)

	// iterate through the AD identifiers for this config
	for _, id := range tpl.ADIdentifiers {
		digests := cache.id2digests[id]
		// remove the template from id2templates
		for i, digest := range digests {
			if digest == d {
				cache.id2digests[id] = append(digests[:i], digests[i+1:]...)
				break
			}
		}
	}

	return nil
}