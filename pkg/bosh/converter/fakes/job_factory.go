// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/converter"
	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedjob/v1alpha1"
)

type FakeJobFactory struct {
	BPMConfigsJobStub        func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)
	bPMConfigsJobMutex       sync.RWMutex
	bPMConfigsJobArgsForCall []struct {
		arg1 manifest.Manifest
	}
	bPMConfigsJobReturns struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	bPMConfigsJobReturnsOnCall map[int]struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	InstanceGroupManifestJobStub        func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)
	instanceGroupManifestJobMutex       sync.RWMutex
	instanceGroupManifestJobArgsForCall []struct {
		arg1 manifest.Manifest
	}
	instanceGroupManifestJobReturns struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	instanceGroupManifestJobReturnsOnCall map[int]struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	VariableInterpolationJobStub        func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)
	variableInterpolationJobMutex       sync.RWMutex
	variableInterpolationJobArgsForCall []struct {
		arg1 manifest.Manifest
	}
	variableInterpolationJobReturns struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	variableInterpolationJobReturnsOnCall map[int]struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeJobFactory) BPMConfigsJob(arg1 manifest.Manifest) (*v1alpha1.ExtendedJob, error) {
	fake.bPMConfigsJobMutex.Lock()
	ret, specificReturn := fake.bPMConfigsJobReturnsOnCall[len(fake.bPMConfigsJobArgsForCall)]
	fake.bPMConfigsJobArgsForCall = append(fake.bPMConfigsJobArgsForCall, struct {
		arg1 manifest.Manifest
	}{arg1})
	fake.recordInvocation("BPMConfigsJob", []interface{}{arg1})
	fake.bPMConfigsJobMutex.Unlock()
	if fake.BPMConfigsJobStub != nil {
		return fake.BPMConfigsJobStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.bPMConfigsJobReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeJobFactory) BPMConfigsJobCallCount() int {
	fake.bPMConfigsJobMutex.RLock()
	defer fake.bPMConfigsJobMutex.RUnlock()
	return len(fake.bPMConfigsJobArgsForCall)
}

func (fake *FakeJobFactory) BPMConfigsJobCalls(stub func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)) {
	fake.bPMConfigsJobMutex.Lock()
	defer fake.bPMConfigsJobMutex.Unlock()
	fake.BPMConfigsJobStub = stub
}

func (fake *FakeJobFactory) BPMConfigsJobArgsForCall(i int) manifest.Manifest {
	fake.bPMConfigsJobMutex.RLock()
	defer fake.bPMConfigsJobMutex.RUnlock()
	argsForCall := fake.bPMConfigsJobArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeJobFactory) BPMConfigsJobReturns(result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.bPMConfigsJobMutex.Lock()
	defer fake.bPMConfigsJobMutex.Unlock()
	fake.BPMConfigsJobStub = nil
	fake.bPMConfigsJobReturns = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) BPMConfigsJobReturnsOnCall(i int, result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.bPMConfigsJobMutex.Lock()
	defer fake.bPMConfigsJobMutex.Unlock()
	fake.BPMConfigsJobStub = nil
	if fake.bPMConfigsJobReturnsOnCall == nil {
		fake.bPMConfigsJobReturnsOnCall = make(map[int]struct {
			result1 *v1alpha1.ExtendedJob
			result2 error
		})
	}
	fake.bPMConfigsJobReturnsOnCall[i] = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) InstanceGroupManifestJob(arg1 manifest.Manifest) (*v1alpha1.ExtendedJob, error) {
	fake.instanceGroupManifestJobMutex.Lock()
	ret, specificReturn := fake.instanceGroupManifestJobReturnsOnCall[len(fake.instanceGroupManifestJobArgsForCall)]
	fake.instanceGroupManifestJobArgsForCall = append(fake.instanceGroupManifestJobArgsForCall, struct {
		arg1 manifest.Manifest
	}{arg1})
	fake.recordInvocation("InstanceGroupManifestJob", []interface{}{arg1})
	fake.instanceGroupManifestJobMutex.Unlock()
	if fake.InstanceGroupManifestJobStub != nil {
		return fake.InstanceGroupManifestJobStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.instanceGroupManifestJobReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeJobFactory) InstanceGroupManifestJobCallCount() int {
	fake.instanceGroupManifestJobMutex.RLock()
	defer fake.instanceGroupManifestJobMutex.RUnlock()
	return len(fake.instanceGroupManifestJobArgsForCall)
}

func (fake *FakeJobFactory) InstanceGroupManifestJobCalls(stub func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)) {
	fake.instanceGroupManifestJobMutex.Lock()
	defer fake.instanceGroupManifestJobMutex.Unlock()
	fake.InstanceGroupManifestJobStub = stub
}

func (fake *FakeJobFactory) InstanceGroupManifestJobArgsForCall(i int) manifest.Manifest {
	fake.instanceGroupManifestJobMutex.RLock()
	defer fake.instanceGroupManifestJobMutex.RUnlock()
	argsForCall := fake.instanceGroupManifestJobArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeJobFactory) InstanceGroupManifestJobReturns(result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.instanceGroupManifestJobMutex.Lock()
	defer fake.instanceGroupManifestJobMutex.Unlock()
	fake.InstanceGroupManifestJobStub = nil
	fake.instanceGroupManifestJobReturns = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) InstanceGroupManifestJobReturnsOnCall(i int, result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.instanceGroupManifestJobMutex.Lock()
	defer fake.instanceGroupManifestJobMutex.Unlock()
	fake.InstanceGroupManifestJobStub = nil
	if fake.instanceGroupManifestJobReturnsOnCall == nil {
		fake.instanceGroupManifestJobReturnsOnCall = make(map[int]struct {
			result1 *v1alpha1.ExtendedJob
			result2 error
		})
	}
	fake.instanceGroupManifestJobReturnsOnCall[i] = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) VariableInterpolationJob(arg1 manifest.Manifest) (*v1alpha1.ExtendedJob, error) {
	fake.variableInterpolationJobMutex.Lock()
	ret, specificReturn := fake.variableInterpolationJobReturnsOnCall[len(fake.variableInterpolationJobArgsForCall)]
	fake.variableInterpolationJobArgsForCall = append(fake.variableInterpolationJobArgsForCall, struct {
		arg1 manifest.Manifest
	}{arg1})
	fake.recordInvocation("VariableInterpolationJob", []interface{}{arg1})
	fake.variableInterpolationJobMutex.Unlock()
	if fake.VariableInterpolationJobStub != nil {
		return fake.VariableInterpolationJobStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.variableInterpolationJobReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeJobFactory) VariableInterpolationJobCallCount() int {
	fake.variableInterpolationJobMutex.RLock()
	defer fake.variableInterpolationJobMutex.RUnlock()
	return len(fake.variableInterpolationJobArgsForCall)
}

func (fake *FakeJobFactory) VariableInterpolationJobCalls(stub func(manifest.Manifest) (*v1alpha1.ExtendedJob, error)) {
	fake.variableInterpolationJobMutex.Lock()
	defer fake.variableInterpolationJobMutex.Unlock()
	fake.VariableInterpolationJobStub = stub
}

func (fake *FakeJobFactory) VariableInterpolationJobArgsForCall(i int) manifest.Manifest {
	fake.variableInterpolationJobMutex.RLock()
	defer fake.variableInterpolationJobMutex.RUnlock()
	argsForCall := fake.variableInterpolationJobArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeJobFactory) VariableInterpolationJobReturns(result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.variableInterpolationJobMutex.Lock()
	defer fake.variableInterpolationJobMutex.Unlock()
	fake.VariableInterpolationJobStub = nil
	fake.variableInterpolationJobReturns = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) VariableInterpolationJobReturnsOnCall(i int, result1 *v1alpha1.ExtendedJob, result2 error) {
	fake.variableInterpolationJobMutex.Lock()
	defer fake.variableInterpolationJobMutex.Unlock()
	fake.VariableInterpolationJobStub = nil
	if fake.variableInterpolationJobReturnsOnCall == nil {
		fake.variableInterpolationJobReturnsOnCall = make(map[int]struct {
			result1 *v1alpha1.ExtendedJob
			result2 error
		})
	}
	fake.variableInterpolationJobReturnsOnCall[i] = struct {
		result1 *v1alpha1.ExtendedJob
		result2 error
	}{result1, result2}
}

func (fake *FakeJobFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bPMConfigsJobMutex.RLock()
	defer fake.bPMConfigsJobMutex.RUnlock()
	fake.instanceGroupManifestJobMutex.RLock()
	defer fake.instanceGroupManifestJobMutex.RUnlock()
	fake.variableInterpolationJobMutex.RLock()
	defer fake.variableInterpolationJobMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeJobFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ converter.JobFactory = new(FakeJobFactory)
