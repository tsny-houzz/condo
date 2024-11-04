package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/manifoldco/promptui"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Klient struct {
	kube  *kubernetes.Clientset
	istio *istio.Clientset
	cfg   *rest.Config
}

func newClient() *Klient {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home := os.Getenv("HOME")
		kubeconfig = fmt.Sprintf("%s/.kube/config", home)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}

	// Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	// Istio client
	istioClient, err := istio.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Istio clientset: %v", err)
	}

	return &Klient{kubeClient, istioClient, config}
}

func (k *Klient) validateCluster() error {
	b, err := exec.Command("kubectl", "config", "current-context").CombinedOutput()
	if err != nil {
		return err
	}
	if !strings.Contains(string(b), "stg-main-eks") {
		return fmt.Errorf("current-context cluster is not stg-main-eks; it is %v", string(b))
	}
	return nil
}

// ListNamespacesWithEmail lists namespaces that have the specified email address in an annotation
func (k *Klient) ListNamespacesWithEmail(email string) error {
	namespaces, err := k.kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %v", err)
	}

	fmt.Printf("Codespaces with owner email '%s':\n", email)
	for _, ns := range namespaces.Items {
		annotations := ns.Annotations
		if annotations != nil && annotations["owner"] == email {
			fmt.Println("-", ns.Name)
		}
	}
	return nil
}

// CreateNamespace creates a new namespace with specific annotations and labels
func (k *Klient) CreateNamespace(namespaceName, email string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("email does not contain '@'")
	}

	if ns, _ := k.getNamespaceByName(namespaceName); ns != nil {
		owner, ok := ns.Annotations["owner"]
		if !ok {
			owner = "no-owner-found"
		}
		return fmt.Errorf("namespace %v already exists: owned by %v", namespaceName, owner)
	}

	codespaceUser := strings.Split(email, "@")[0]
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
			Annotations: map[string]string{
				"owner": email,
			},
			Labels: map[string]string{
				"codespace-user":  codespaceUser,
				"istio-injection": "enabled",
			},
		},
	}

	_, err := k.kube.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

func (k *Klient) getNamespaceByName(namespaceName string) (*v1.Namespace, error) {
	namespace, err := k.kube.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace %s: %v", namespaceName, err)
	}
	return namespace, nil
}

func hasOwnerAnnotation(namespace *v1.Namespace, email string) bool {
	ownerEmail, ok := namespace.Annotations["owner"]
	return ok && ownerEmail == email
}

func (k *Klient) getResourceNames(namespace string) ([]string, error) {
	var names []string

	// Get all services
	services, err := k.kube.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}
	for _, svc := range services.Items {
		names = append(names, fmt.Sprintf("Service %s", svc.Name))
	}

	// Get all deployments
	deployments, err := k.kube.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}
	for _, deploy := range deployments.Items {
		names = append(names, fmt.Sprintf("Deployment %s", deploy.Name))
	}

	virtualServices, err := k.istio.NetworkingV1alpha3().VirtualServices(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list virtual services: %v", err)
	}
	for _, vs := range virtualServices.Items {
		names = append(names, fmt.Sprintf("VirtualService %s", vs.Name))
	}

	return names, nil
}

func (k *Klient) selectResource(namespace string) (string, error) {
	resources, err := k.getResourceNames(namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get resources: %v", err)
	}

	p := promptui.Select{
		Label: "Choose a resource",
		Items: resources,
	}
	_, resource, err := p.Run()
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(resource, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format; expected 'ResourceType name'")
	}
	resourceType := strings.TrimSpace(parts[0])
	resourceName := strings.TrimSpace(parts[1])

	b, err := exec.Command("kubectl", "get", resourceType, resourceName, "-o", "yaml").CombinedOutput()
	if err != nil {
		return "", err
	}

	fmt.Println("---")
	if err := quick.Highlight(os.Stdout, string(b), "yaml", "terminal", "monokai"); err != nil {
		return "", fmt.Errorf("failed to highlight YAML output: %v", err)
	}
	return resource, nil
}
