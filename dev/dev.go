// Package dev deploys the gospell app to GCP using a Cloud Run and Pulumi
package dev

import (
	cloudrun "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Go() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new Cloud Run Service
		gospellService, err := cloudrun.NewService(ctx, "gospell-service", &cloudrun.ServiceArgs{
			Location: pulumi.String("us-central1"),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							// TODO - move to config variable
							Image: pulumi.String("us-central1-docker.pkg.dev/main-composite-285417/cloud-run-source-deploy/gospell/gospell@sha256:00bd52d9a2bf274bea7bbb5ef229de046bfb8cbed673cead864ef632677d9f81"),
						},
					},
				},
			},
		})

		if err != nil {
			return err
		}

		// Create an IAM member to make the service publicly accessible.
		_, err = cloudrun.NewIamMember(ctx, "invoker", &cloudrun.IamMemberArgs{
			Service:  gospellService.Name,
			Location: gospellService.Location,
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
}
