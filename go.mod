module my-controller

go 1.15

require (
    k8s.io/api v0.0.0-20201104162213-01c5338f427f
	k8s.io/apimachinery v0.0.0-20201104162036-79ef3cbd919a
	k8s.io/client-go v0.0.0-20201104162436-68bb4a9525d8
	k8s.io/code-generator v0.0.0-20201104161901-3609764c976f
	k8s.io/klog/v2 v2.4.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20201104162213-01c5338f427f
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20201104162036-79ef3cbd919a
	k8s.io/client-go => k8s.io/client-go v0.0.0-20201104162436-68bb4a9525d8
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20201104161901-3609764c976f
)
