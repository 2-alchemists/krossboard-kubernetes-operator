//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 2ALCHEMISTS SAS.

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
func (in *KbComponentInstance) DeepCopyInto(out *KbComponentInstance) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KbComponentInstance.
func (in *KbComponentInstance) DeepCopy() *KbComponentInstance {
	if in == nil {
		return nil
	}
	out := new(KbComponentInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KoaInstance) DeepCopyInto(out *KoaInstance) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KoaInstance.
func (in *KoaInstance) DeepCopy() *KoaInstance {
	if in == nil {
		return nil
	}
	out := new(KoaInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Krossboard) DeepCopyInto(out *Krossboard) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Krossboard.
func (in *Krossboard) DeepCopy() *Krossboard {
	if in == nil {
		return nil
	}
	out := new(Krossboard)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Krossboard) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KrossboardList) DeepCopyInto(out *KrossboardList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Krossboard, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KrossboardList.
func (in *KrossboardList) DeepCopy() *KrossboardList {
	if in == nil {
		return nil
	}
	out := new(KrossboardList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KrossboardList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KrossboardSpec) DeepCopyInto(out *KrossboardSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KrossboardSpec.
func (in *KrossboardSpec) DeepCopy() *KrossboardSpec {
	if in == nil {
		return nil
	}
	out := new(KrossboardSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KrossboardStatus) DeepCopyInto(out *KrossboardStatus) {
	*out = *in
	if in.KoaInstances != nil {
		in, out := &in.KoaInstances, &out.KoaInstances
		*out = make([]KoaInstance, len(*in))
		copy(*out, *in)
	}
	if in.KbComponentInstances != nil {
		in, out := &in.KbComponentInstances, &out.KbComponentInstances
		*out = make([]KbComponentInstance, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KrossboardStatus.
func (in *KrossboardStatus) DeepCopy() *KrossboardStatus {
	if in == nil {
		return nil
	}
	out := new(KrossboardStatus)
	in.DeepCopyInto(out)
	return out
}
