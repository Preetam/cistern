// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package codepipeline provides the client and types for making API
// requests to AWS CodePipeline.
//
// Overview
//
// This is the AWS CodePipeline API Reference. This guide provides descriptions
// of the actions and data types for AWS CodePipeline. Some functionality for
// your pipeline is only configurable through the API. For additional information,
// see the AWS CodePipeline User Guide (http://docs.aws.amazon.com/codepipeline/latest/userguide/welcome.html).
//
// You can use the AWS CodePipeline API to work with pipelines, stages, actions,
// gates, and transitions, as described below.
//
// Pipelines are models of automated release processes. Each pipeline is uniquely
// named, and consists of actions, gates, and stages.
//
// You can work with pipelines by calling:
//
//    * CreatePipeline, which creates a uniquely-named pipeline.
//
//    * DeletePipeline, which deletes the specified pipeline.
//
//    * GetPipeline, which returns information about a pipeline structure.
//
//    * GetPipelineExecution, which returns information about a specific execution
//    of a pipeline.
//
//    * GetPipelineState, which returns information about the current state
//    of the stages and actions of a pipeline.
//
//    * ListPipelines, which gets a summary of all of the pipelines associated
//    with your account.
//
//    * StartPipelineExecution, which runs the the most recent revision of an
//    artifact through the pipeline.
//
//    * UpdatePipeline, which updates a pipeline with edits or changes to the
//    structure of the pipeline.
//
// Pipelines include stages, which are logical groupings of gates and actions.
// Each stage contains one or more actions that must complete before the next
// stage begins. A stage will result in success or failure. If a stage fails,
// then the pipeline stops at that stage and will remain stopped until either
// a new version of an artifact appears in the source location, or a user takes
// action to re-run the most recent artifact through the pipeline. You can call
// GetPipelineState, which displays the status of a pipeline, including the
// status of stages in the pipeline, or GetPipeline, which returns the entire
// structure of the pipeline, including the stages of that pipeline. For more
// information about the structure of stages and actions, also refer to the
// AWS CodePipeline Pipeline Structure Reference (http://docs.aws.amazon.com/codepipeline/latest/userguide/pipeline-structure.html).
//
// Pipeline stages include actions, which are categorized into categories such
// as source or build actions performed within a stage of a pipeline. For example,
// you can use a source action to import artifacts into a pipeline from a source
// such as Amazon S3. Like stages, you do not work with actions directly in
// most cases, but you do define and interact with actions when working with
// pipeline operations such as CreatePipeline and GetPipelineState.
//
// Pipelines also include transitions, which allow the transition of artifacts
// from one stage to the next in a pipeline after the actions in one stage complete.
//
// You can work with transitions by calling:
//
//    * DisableStageTransition, which prevents artifacts from transitioning
//    to the next stage in a pipeline.
//
//    * EnableStageTransition, which enables transition of artifacts between
//    stages in a pipeline.
//
// Using the API to integrate with AWS CodePipeline
//
// For third-party integrators or developers who want to create their own integrations
// with AWS CodePipeline, the expected sequence varies from the standard API
// user. In order to integrate with AWS CodePipeline, developers will need to
// work with the following items:
//
// Jobs, which are instances of an action. For example, a job for a source action
// might import a revision of an artifact from a source.
//
// You can work with jobs by calling:
//
//    * AcknowledgeJob, which confirms whether a job worker has received the
//    specified job,
//
//    * GetJobDetails, which returns the details of a job,
//
//    * PollForJobs, which determines whether there are any jobs to act upon,
//
//
//    * PutJobFailureResult, which provides details of a job failure, and
//
//    * PutJobSuccessResult, which provides details of a job success.
//
// Third party jobs, which are instances of an action created by a partner action
// and integrated into AWS CodePipeline. Partner actions are created by members
// of the AWS Partner Network.
//
// You can work with third party jobs by calling:
//
//    * AcknowledgeThirdPartyJob, which confirms whether a job worker has received
//    the specified job,
//
//    * GetThirdPartyJobDetails, which requests the details of a job for a partner
//    action,
//
//    * PollForThirdPartyJobs, which determines whether there are any jobs to
//    act upon,
//
//    * PutThirdPartyJobFailureResult, which provides details of a job failure,
//    and
//
//    * PutThirdPartyJobSuccessResult, which provides details of a job success.
//
// See https://docs.aws.amazon.com/goto/WebAPI/codepipeline-2015-07-09 for more information on this service.
//
// See codepipeline package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/codepipeline/
//
// Using the Client
//
// To use the client for AWS CodePipeline you will first need
// to create a new instance of it.
//
// When creating a client for an AWS service you'll first need to have a Session
// already created. The Session provides configuration that can be shared
// between multiple service clients. Additional configuration can be applied to
// the Session and service's client when they are constructed. The aws package's
// Config type contains several fields such as Region for the AWS Region the
// client should make API requests too. The optional Config value can be provided
// as the variadic argument for Sessions and client creation.
//
// Once the service's client is created you can use it to make API requests the
// AWS service. These clients are safe to use concurrently.
//
//   // Create a session to share configuration, and load external configuration.
//   sess := session.Must(session.NewSession())
//
//   // Create the service's client with the session.
//   svc := codepipeline.New(sess)
//
// See the SDK's documentation for more information on how to use service clients.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws package's Config type for more information on configuration options.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the AWS CodePipeline client CodePipeline for more
// information on creating the service's client.
// https://docs.aws.amazon.com/sdk-for-go/api/service/codepipeline/#New
//
// Once the client is created you can make an API request to the service.
// Each API method takes a input parameter, and returns the service response
// and an error.
//
// The API method will document which error codes the service can be returned
// by the operation if the service models the API operation's errors. These
// errors will also be available as const strings prefixed with "ErrCode".
//
//   result, err := svc.AcknowledgeJob(params)
//   if err != nil {
//       // Cast err to awserr.Error to handle specific error codes.
//       aerr, ok := err.(awserr.Error)
//       if ok && aerr.Code() == <error code to check for> {
//           // Specific error code handling
//       }
//       return err
//   }
//
//   fmt.Println("AcknowledgeJob result:")
//   fmt.Println(result)
//
// Using the Client with Context
//
// The service's client also provides methods to make API requests with a Context
// value. This allows you to control the timeout, and cancellation of pending
// requests. These methods also take request Option as variadic parameter to apply
// additional configuration to the API request.
//
//   ctx := context.Background()
//
//   result, err := svc.AcknowledgeJobWithContext(ctx, params)
//
// See the request package documentation for more information on using Context pattern
// with the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/request/
package codepipeline
