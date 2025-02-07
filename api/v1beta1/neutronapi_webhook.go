/*
Copyright 2022.

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

//
// Generated by:
//
// operator-sdk create webhook --group neutron --version v1beta1 --kind NeutronAPI --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"

	"github.com/openstack-k8s-operators/lib-common/modules/common/service"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
)

// NeutronAPIDefaults -
type NeutronAPIDefaults struct {
	ContainerImageURL string
	APITimeout        int
}

var neutronAPIDefaults NeutronAPIDefaults

// log is for logging in this package.
var neutronapilog = logf.Log.WithName("neutronapi-resource")

// SetupNeutronAPIDefaults - initialize NeutronAPI spec defaults for use with either internal or external webhooks
func SetupNeutronAPIDefaults(defaults NeutronAPIDefaults) {
	neutronAPIDefaults = defaults
	neutronapilog.Info("NeutronAPI defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *NeutronAPI) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-neutron-openstack-org-v1beta1-neutronapi,mutating=true,failurePolicy=fail,sideEffects=None,groups=neutron.openstack.org,resources=neutronapis,verbs=create;update,versions=v1beta1,name=mneutronapi.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &NeutronAPI{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NeutronAPI) Default() {
	neutronapilog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this NeutronAPI spec
func (spec *NeutronAPISpec) Default() {
	// only container image validations go here
	if spec.ContainerImage == "" {
		spec.ContainerImage = neutronAPIDefaults.ContainerImageURL
	}
	spec.NeutronAPISpecCore.Default()
}

// Default - set defaults for this NeutronAPI spec core. This version gets used by OpenStackControlplane
func (spec *NeutronAPISpecCore) Default() {
	if spec.APITimeout == 0 {
		spec.APITimeout = neutronAPIDefaults.APITimeout
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-neutron-openstack-org-v1beta1-neutronapi,mutating=false,failurePolicy=fail,sideEffects=None,groups=neutron.openstack.org,resources=neutronapis,verbs=create;update,versions=v1beta1,name=vneutronapi.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &NeutronAPI{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NeutronAPI) ValidateCreate() (admission.Warnings, error) {
	neutronapilog.Info("validate create", "name", r.Name)

	allErrs := field.ErrorList{}
	basePath := field.NewPath("spec")

	// When a TopologyRef CR is referenced, fail if a different Namespace is
	// referenced because is not supported
	if r.Spec.TopologyRef != nil {
		if err := topologyv1.ValidateTopologyNamespace(r.Spec.TopologyRef.Namespace, *basePath, r.Namespace); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if err := r.Spec.ValidateCreate(basePath); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("NeutronAPI").GroupKind(), r.Name, allErrs)
	}

	return nil, nil
}

// ValidateCreate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an NeutronAPI spec.
func (r *NeutronAPISpec) ValidateCreate(basePath *field.Path) field.ErrorList {
	return r.NeutronAPISpecCore.ValidateCreate(basePath)
}

func (r *NeutronAPISpecCore) ValidateCreate(basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(basePath.Child("override").Child("service"), r.Override.Service)...)

	allErrs = append(allErrs, ValidateDefaultConfigOverwrite(basePath, r.DefaultConfigOverwrite)...)

	return allErrs
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NeutronAPI) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	neutronapilog.Info("validate update", "name", r.Name)

	oldNeutronAPI, ok := old.(*NeutronAPI)
	if !ok || oldNeutronAPI == nil {
		return nil, apierrors.NewInternalError(fmt.Errorf("unable to convert existing object"))
	}

	allErrs := field.ErrorList{}
	basePath := field.NewPath("spec")

	// When a TopologyRef CR is referenced, fail if a different Namespace is
	// referenced because is not supported
	if r.Spec.TopologyRef != nil {
		if err := topologyv1.ValidateTopologyNamespace(r.Spec.TopologyRef.Namespace, *basePath, r.Namespace); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if err := r.Spec.ValidateUpdate(oldNeutronAPI.Spec, basePath); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("NeutronAPI").GroupKind(), r.Name, allErrs)
	}

	return nil, nil
}

// ValidateUpdate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an neutron spec.
func (spec *NeutronAPISpec) ValidateUpdate(old NeutronAPISpec, basePath *field.Path) field.ErrorList {
	return spec.NeutronAPISpecCore.ValidateUpdate(old.NeutronAPISpecCore, basePath)
}

func (spec *NeutronAPISpecCore) ValidateUpdate(old NeutronAPISpecCore, basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(basePath.Child("override").Child("service"), spec.Override.Service)...)
	// validate the defaultConfigOverwrite is valid
	allErrs = append(allErrs, ValidateDefaultConfigOverwrite(basePath, spec.DefaultConfigOverwrite)...)

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NeutronAPI) ValidateDelete() (admission.Warnings, error) {
	neutronapilog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (spec *NeutronAPISpec) GetDefaultRouteAnnotations() (annotations map[string]string) {
	return spec.NeutronAPISpecCore.GetDefaultRouteAnnotations()
}

func (spec *NeutronAPISpecCore) GetDefaultRouteAnnotations() (annotations map[string]string) {
	return map[string]string{
		"haproxy.router.openshift.io/timeout": fmt.Sprintf("%ds", neutronAPIDefaults.APITimeout),
	}
}

func ValidateDefaultConfigOverwrite(
	basePath *field.Path,
	validateConfigOverwrite map[string]string,
) field.ErrorList {
	var errors field.ErrorList
	for requested := range validateConfigOverwrite {
		if requested != "policy.yaml" {
			errors = append(
				errors,
				field.Invalid(
					basePath.Child("defaultConfigOverwrite"),
					requested,
					"Only the following keys are valid: policy.yaml",
				),
			)
		}
	}
	return errors
}

// SetDefaultRouteAnnotations sets HAProxy timeout values of the route
func (spec *NeutronAPISpecCore) SetDefaultRouteAnnotations(annotations map[string]string) {
	const haProxyAnno = "haproxy.router.openshift.io/timeout"
	// Use a custom annotation to flag when the operator has set the default HAProxy timeout
	// With the annotation func determines when to overwrite existing HAProxy timeout with the APITimeout
	const neutronAnno = "api.neutron.openstack.org/timeout"
	valNeutronAPI, okNeutronAPI := annotations[neutronAnno]
	valHAProxy, okHAProxy := annotations[haProxyAnno]
	// Human operator set the HAProxy timeout manually
	if (!okNeutronAPI && okHAProxy) {
		return
	}
	// Human operator modified the HAProxy timeout manually without removing the NeutronAPI flag
	if (okNeutronAPI && okHAProxy && valNeutronAPI != valHAProxy) {
		delete(annotations, neutronAnno)
		return
	}
	timeout := fmt.Sprintf("%ds", spec.APITimeout)
	annotations[neutronAnno] = timeout
	annotations[haProxyAnno] = timeout
}
