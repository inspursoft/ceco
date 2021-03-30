// +build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostAndPath) DeepCopyInto(out *HostAndPath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostAndPath.
func (in *HostAndPath) DeepCopy() *HostAndPath {
	if in == nil {
		return nil
	}
	out := new(HostAndPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsCo) DeepCopyInto(out *NatsCo) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsCo.
func (in *NatsCo) DeepCopy() *NatsCo {
	if in == nil {
		return nil
	}
	out := new(NatsCo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NatsCo) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsCoList) DeepCopyInto(out *NatsCoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NatsCo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsCoList.
func (in *NatsCoList) DeepCopy() *NatsCoList {
	if in == nil {
		return nil
	}
	out := new(NatsCoList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NatsCoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsCoSpec) DeepCopyInto(out *NatsCoSpec) {
	*out = *in
	if in.NatsServers != nil {
		in, out := &in.NatsServers, &out.NatsServers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Source = in.Source
	if in.Destinations != nil {
		in, out := &in.Destinations, &out.Destinations
		*out = make([]HostAndPath, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsCoSpec.
func (in *NatsCoSpec) DeepCopy() *NatsCoSpec {
	if in == nil {
		return nil
	}
	out := new(NatsCoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NatsCoStatus) DeepCopyInto(out *NatsCoStatus) {
	*out = *in
	if in.Destination != nil {
		in, out := &in.Destination, &out.Destination
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NatsCoStatus.
func (in *NatsCoStatus) DeepCopy() *NatsCoStatus {
	if in == nil {
		return nil
	}
	out := new(NatsCoStatus)
	in.DeepCopyInto(out)
	return out
}
