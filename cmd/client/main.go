// This application is to demonstrate that clients with different versions of
// an object API can continue to use the stored object API which may be in a
// different version.
//
// Example usage:
// 		cmd createv1 foo
// 		cmd getv2 foo
// 		cmd deletev2 foo
//
package main

import (
	"context"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	batchv1 "github.com/darkowlzz/crd-api-conversion-demo/api/v1"
	batchv2 "github.com/darkowlzz/crd-api-conversion-demo/api/v2"
)

var (
	scheme = runtime.NewScheme()
	log    = ctrl.Log.WithName("cronjob-crd-client")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(batchv1.AddToScheme(scheme))
	utilruntime.Must(batchv2.AddToScheme(scheme))
}

func main() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if len(os.Args) < 3 {
		log.Info("help: cmd <subcommand> <resource-name>")
		os.Exit(1)
	}
	subCmd := os.Args[1]
	resName := os.Args[2]

	cl, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "failed to create a client")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	v1Cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: "default",
		},
		Spec: batchv1.CronJobSpec{
			Foo: "lala",
		},
	}
	gvk1, err := apiutil.GVKForObject(v1Cronjob, scheme)
	if err != nil {
		log.Error(err, "failed to get object GVK")
	}

	v2Cronjob := &batchv2.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: "default",
		},
		Spec: batchv2.CronJobSpec{
			Foo: "wawa",
		},
	}
	gvk2, err := apiutil.GVKForObject(v2Cronjob, scheme)
	if err != nil {
		log.Error(err, "failed to get object GVK")
	}

	switch subCmd {
	case "createv1":
		if err := createObj(ctx, cl, v1Cronjob); err != nil {
			log.Error(err, "failed to create v1 object")
			os.Exit(1)
		}
		log.Info("created successfully", "name", v1Cronjob.GetName(), "GVK", gvk1)
	case "deletev1":
		if err := deleteObj(ctx, cl, v1Cronjob); err != nil {
			log.Error(err, "failed to delete v1 object")
			os.Exit(1)
		}
		log.Info("deleted successfully", "name", v1Cronjob.GetName(), "GVK", gvk1)
	case "getv1":
		obj := &batchv1.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      v1Cronjob.Name,
				Namespace: v1Cronjob.Namespace,
			},
		}
		if err := getObj(ctx, cl, obj); err != nil {
			log.Error(err, "failed to get v1 object")
			os.Exit(1)
		}
		log.Info("got successfully using batchv1 API", "obj", obj)
	case "createv2":
		if err := createObj(ctx, cl, v2Cronjob); err != nil {
			log.Error(err, "failed to create v2 object")
			os.Exit(1)
		}
		log.Info("created successfully", "name", v2Cronjob.GetName(), "GVK", gvk2)
	case "deletev2":
		if err := deleteObj(ctx, cl, v2Cronjob); err != nil {
			log.Error(err, "failed to delete v2 object")
			os.Exit(1)
		}
		log.Info("deleted successfully", "name", v2Cronjob.GetName(), "GVK", gvk2)
	case "getv2":
		obj := &batchv2.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      v2Cronjob.Name,
				Namespace: v2Cronjob.Namespace,
			},
		}
		if err := getObj(ctx, cl, obj); err != nil {
			log.Error(err, "failed to get v2 object")
			os.Exit(1)
		}
		log.Info("got successfully using batchv2 API", "obj", obj)
	default:
		log.Info("Use arguments createv1, deletev1, createv2, deletev2")
	}
}

func createObj(ctx context.Context, cl client.Client, obj client.Object) error {
	return cl.Create(ctx, obj)
}

func deleteObj(ctx context.Context, cl client.Client, obj client.Object) error {
	return cl.Delete(ctx, obj)
}

func getObj(ctx context.Context, cl client.Client, obj client.Object) error {
	return cl.Get(ctx, client.ObjectKeyFromObject(obj), obj)
}
