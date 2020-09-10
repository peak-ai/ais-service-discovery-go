package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/peak-ai/ais-service-discovery-go/pkg/automate"
	"github.com/peak-ai/ais-service-discovery-go/pkg/backends"
	"github.com/peak-ai/ais-service-discovery-go/pkg/function"
	"github.com/peak-ai/ais-service-discovery-go/pkg/locator"
	"github.com/peak-ai/ais-service-discovery-go/pkg/logger"
	"github.com/peak-ai/ais-service-discovery-go/pkg/pubsub"
	"github.com/peak-ai/ais-service-discovery-go/pkg/queue"
	"github.com/peak-ai/ais-service-discovery-go/pkg/tracer"
)

// WithAWSBackend initializes the Discovery object with default AWS services.
// Override these services by using the Set methods.
//
// This is duplicated in `discovery.go`, this is known, this is for consistency and
// backwards compatibility. New users will find 'backends' in this sub-package,
// so that they're all in one place. But original users, should expect to be able
// to use the existing functionality without disruption.
func Factory() backends.Option {
	return func(args *backends.Options) {
		sess := session.Must(session.NewSession())
		args.QueueAdapter = queue.NewSQSAdapter(sqs.New(sess))
		args.FunctionAdapter = function.NewLambdaAdapter(lambda.New(sess))
		args.AutomateAdapter = automate.NewSSMAdapter(ssm.New(sess))
		args.PubsubAdapter = pubsub.NewSNSAdapter(sns.New(sess))
		args.Locator = locator.NewCloudmapLocator(servicediscovery.New(sess))
		args.LogAdapter = logger.NewSTDOutAdapter()
		args.TraceAdapter = tracer.NewXrayAdapter()
	}
}
