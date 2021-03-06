//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package resources

import (
	"fmt"

	researchv1alpha1 "github.com/IBM/integrity-enforcer/operator/pkg/apis/research/v1alpha1"
	admv1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

//service
func BuildServiceForCR(cr *researchv1alpha1.IntegrityEnforcer) *corev1.Service {
	var targetport intstr.IntOrString
	targetport.Type = intstr.String
	targetport.StrVal = "ac-api"
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.WebhookServiceName,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       443,
					TargetPort: targetport, //"ac-api"
				},
			},
			Selector: cr.Spec.SelectorLabels,
		},
	}
	return svc
}

//webhook configuration
func BuildMutatingWebhookConfigurationForIE(cr *researchv1alpha1.IntegrityEnforcer) *admv1.MutatingWebhookConfiguration {
	var scopetype *admv1.ScopeType
	namespaced := admv1.NamespacedScope
	scopetype = &namespaced

	var path *string
	mutate := "/mutate"
	path = &mutate

	var empty []byte

	sideEffect := admv1.SideEffectClassNone

	wc := &admv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.WebhookConfigName,
			Namespace: cr.Namespace,
		},
		Webhooks: []admv1.MutatingWebhook{
			{
				Name: fmt.Sprintf("ac-server.%s.svc", cr.Namespace),
				NamespaceSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      "integrity-enforced",
							Operator: "In",
							Values: []string{
								"true",
							},
						},
					},
				},
				ClientConfig: admv1.WebhookClientConfig{
					Service: &admv1.ServiceReference{
						Name:      cr.Spec.WebhookServiceName,
						Namespace: cr.Namespace,
						Path:      path, //"/mutate"
					},
					CABundle: empty,
				},
				Rules: []admv1.RuleWithOperations{
					{
						Operations: []admv1.OperationType{
							admv1.Create, admv1.Delete, admv1.Update,
						},
						Rule: admv1.Rule{
							APIGroups: []string{
								"*",
							},
							APIVersions: []string{
								"*",
							},
							Resources: []string{
								"*",
							},
							Scope: scopetype, //admv1.NamespacedScope
						},
					},
				},
				SideEffects: &sideEffect,
			},
		},
	}
	return wc
}
