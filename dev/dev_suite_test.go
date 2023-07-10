package dev_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	cloudrun "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var err error
var expected []string
var imageURI string
var stackName string
var s auto.Stack
var resp *http.Response
var ctx context.Context

func TestDev(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dev Suite")
}

var _ = BeforeSuite(func() {
	expected = []string{"Hotel", "India"}
	stackName = "desteves/myProj/tdd"
	imageURI = "us-central1-docker.pkg.dev/main-composite-285417/cloud-run-source-deploy/gospell/gospell@sha256:00bd52d9a2bf274bea7bbb5ef229de046bfb8cbed673cead864ef632677d9f81"
	ctx = context.Background()
	s, err = auto.UpsertStackInlineSource(ctx, stackName, "myProj", func(ctx *pulumi.Context) error {
		gospellService, err := cloudrun.NewService(ctx, "gospell-test-service", &cloudrun.ServiceArgs{
			Project:  pulumi.String("main-composite-285417"),
			Location: pulumi.String("us-central1"),
			// Metadata: &cloudrun.ServiceMetadataArgs{
			// 	Annotations: pulumi.StringMap{
			// 		"run.googleapis.com/ingress":        pulumi.String("all"),
			// 		"run.googleapis.com/ingress-status": pulumi.String("all"),
			// 	},
			// },
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: pulumi.String(imageURI),
						},
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Create an IAM member to make the service publicly accessible.
		_, err = cloudrun.NewIamMember(ctx, "invoker", &cloudrun.IamMemberArgs{
			Service:  gospellService.Name,
			Location: gospellService.Location,
			Project:  gospellService.Project,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}
		// Exports
		ctx.Export("srvurl", gospellService.Statuses.Index(pulumi.Int(0)).Url())
		return nil
	})
	inProgress := true
	for inProgress {
		time.Sleep(time.Second * 2)
		upStatus, err := s.Info(ctx)
		Expect(err).NotTo(HaveOccurred())
		inProgress = upStatus.UpdateInProgress
	}
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// destroy stack and resources
	fmt.Println("Destroying stack...")
	inProgress := true
	for inProgress {
		time.Sleep(time.Second * 2)
		upStatus, err := s.Info(ctx)
		Expect(err).NotTo(HaveOccurred())
		inProgress = upStatus.UpdateInProgress
	}
	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))
	Expect(err).NotTo(HaveOccurred())
	inProgress = true
	for inProgress {
		time.Sleep(time.Second * 2)
		upStatus, err := s.Info(ctx)
		Expect(err).NotTo(HaveOccurred())
		inProgress = upStatus.UpdateInProgress
	}
	err = s.Workspace().RemoveStack(ctx, s.Name())
	Expect(err).NotTo(HaveOccurred())
})
