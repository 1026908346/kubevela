/*
Copyright 2021 The KubeVela Authors.

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

package resourcekeeper

import (
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha1"
)

type rtConfig struct {
	useRoot bool
	skipRT  bool
}

// MetaOnlyOption record only meta part in resourcetracker, which disables the configuration-drift-prevention
type MetaOnlyOption struct{}

// ApplyToDispatchConfig apply change to dispatch config
func (option MetaOnlyOption) ApplyToDispatchConfig(cfg *dispatchConfig) { cfg.metaOnly = true }

// SkipRTOption skip the rt recording during dispatch/delete resources, which means the resource will not be controlled
// by application resourcetracker
type SkipRTOption struct{}

// ApplyToDispatchConfig apply change to dispatch config
func (option SkipRTOption) ApplyToDispatchConfig(cfg *dispatchConfig) { cfg.skipRT = true }

// ApplyToDeleteConfig apply change to delete config
func (option SkipRTOption) ApplyToDeleteConfig(cfg *deleteConfig) { cfg.skipRT = true }

// UseRootOption let the recording and management of target resource belongs to the RootRT instead of VersionedRT. This
// will let the resource be alive as long as the application is still alive.
type UseRootOption struct{}

// ApplyToDispatchConfig apply change to dispatch config
func (option UseRootOption) ApplyToDispatchConfig(cfg *dispatchConfig) { cfg.useRoot = true }

// ApplyToDeleteConfig apply change to delete config
func (option UseRootOption) ApplyToDeleteConfig(cfg *deleteConfig) { cfg.useRoot = true }

// PassiveGCOption disable the active gc for outdated versions. Old versioned resourcetracker will not be recycled
// except all of their managed resources have already been deleted or controlled by later resourcetrackers.
type PassiveGCOption struct{}

// ApplyToGCConfig apply change to gc config
func (option PassiveGCOption) ApplyToGCConfig(cfg *gcConfig) { cfg.passive = true }

// DisableMarkStageGCOption disable the mark stage in gc process (no rt will be marked to be deleted)
// this option should be switched on when application workflow is suspending/terminating since workflow is not
// finished so outdated versions should be kept
type DisableMarkStageGCOption struct{}

// ApplyToGCConfig apply change to gc config
func (option DisableMarkStageGCOption) ApplyToGCConfig(cfg *gcConfig) { cfg.disableMark = true }

// DisableGCComponentRevisionOption disable the component revision gc process
// this option should be switched on when application workflow is suspending/terminating
type DisableGCComponentRevisionOption struct{}

// ApplyToGCConfig apply change to gc config
func (option DisableGCComponentRevisionOption) ApplyToGCConfig(cfg *gcConfig) {
	cfg.disableComponentRevisionGC = true
}

// GarbageCollectStrategyOption apply garbage collect strategy to resourcetracker recording
type GarbageCollectStrategyOption v1alpha1.GarbageCollectStrategy

func (option GarbageCollectStrategyOption) applyToRTConfig(cfg *rtConfig) {
	switch v1alpha1.GarbageCollectStrategy(option) {
	case v1alpha1.GarbageCollectStrategyOnAppUpdate:
		cfg.skipRT = false
		cfg.useRoot = false
	case v1alpha1.GarbageCollectStrategyOnAppDelete:
		cfg.skipRT = false
		cfg.useRoot = true
	case v1alpha1.GarbageCollectStrategyNever:
		cfg.skipRT = true
	}
}

// ApplyToDispatchConfig apply change to dispatch config
func (option GarbageCollectStrategyOption) ApplyToDispatchConfig(cfg *dispatchConfig) {
	option.applyToRTConfig(&cfg.rtConfig)
}

// ApplyToDeleteConfig apply change to delete config
func (option GarbageCollectStrategyOption) ApplyToDeleteConfig(cfg *deleteConfig) {
	option.applyToRTConfig(&cfg.rtConfig)
}