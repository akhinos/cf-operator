package statefulset

import (
	"context"
	"time"

	"gomodules.xyz/jsonpatch/v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"code.cloudfoundry.org/quarks-utils/pkg/ctxlog"
	"code.cloudfoundry.org/quarks-utils/pkg/pointers"

	cfcfg "code.cloudfoundry.org/quarks-utils/pkg/config"
	helper "code.cloudfoundry.org/quarks-utils/testing/testhelper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("When the muatating webhook handles a statefulset", func() {
	var (
		log     *zap.SugaredLogger
		ctx     context.Context
		decoder *admission.Decoder
		mutator admission.Handler
		old     v1beta2.StatefulSet
		new     v1beta2.StatefulSet
	)

	BeforeEach(func() {
		_, log = helper.NewTestLogger()
		ctx = ctxlog.NewParentContext(log)
		old = v1beta2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-statefulset",
				Namespace: "test",
				Annotations: map[string]string{
					AnnotationCanaryRolloutEnabled: "true",
				},
			},
			Spec: v1beta2.StatefulSetSpec{
				Replicas: pointers.Int32(2),
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name: "test-container",
						}},
					},
				},
			},
		}
	})

	JustBeforeEach(func() {
		scheme := runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).To(Succeed())

		decoder, _ = admission.NewDecoder(scheme)
		mutator = NewMutator(log, &cfcfg.Config{CtxTimeOut: 10 * time.Second})
		mutator.(admission.DecoderInjector).InjectDecoder(decoder)
	})

	Context("with no change in pod template", func() {
		BeforeEach(func() {
			new = old
		})
		It("no rollout is triggered", func() {

			oldRaw, _ := json.Marshal(old)
			newRaw, _ := json.Marshal(new)

			response := mutator.Handle(ctx, admission.Request{
				AdmissionRequest: admissionv1beta1.AdmissionRequest{
					OldObject: runtime.RawExtension{Raw: oldRaw},
					Object:    runtime.RawExtension{Raw: newRaw},
				},
			})
			Expect(response.AdmissionResponse.Allowed).To(BeTrue())
			Expect(response.Patches).To(BeEmpty())
		})
	})

	Context("when pod template changes", func() {
		BeforeEach(func() {
			old.DeepCopyInto(&new)
			new.Spec.Template.Spec.Containers[0].Name = "changed-name"
		})
		It("rollout is triggered", func() {

			oldRaw, _ := json.Marshal(old)
			newRaw, _ := json.Marshal(new)

			response := mutator.Handle(ctx, admission.Request{
				AdmissionRequest: admissionv1beta1.AdmissionRequest{
					OldObject: runtime.RawExtension{Raw: oldRaw},
					Object:    runtime.RawExtension{Raw: newRaw},
				},
			})

			Expect(response.Patches).To(ContainElement(
				jsonpatch.Operation{Operation: "add", Path: "/metadata/annotations/quarks.cloudfoundry.org~1canary-rollout", Value: "Pending"},
			))
			Expect(response.Patches).To(ContainElement(
				jsonpatch.Operation{Operation: "add", Path: "/spec/updateStrategy/type", Value: "RollingUpdate"},
			))

			// Does not work because no deepequal check (value is a map/reference)
			//Expect(response.Patches).To(ContainElement(
			//	jsonpatch.Operation{Operation: "add", Path: "/spec/updateStrategy/rollingUpdate", Value: map[string]interface{}{"partition": 0}},
			//))

			// Does not work because of unix timestamp
			//Expect(response.Patches).To(ContainElement(
			//	jsonpatch.Operation{Operation: "add", Path: "/metadata/annotations/quarks.cloudfoundry.org~1update-start-time", Value: "1574265011"},
			//))

			Expect(response.AdmissionResponse.Allowed).To(BeTrue())
		})
	})

})
