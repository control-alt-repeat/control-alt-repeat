// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	internalauthsmithy "github.com/aws/aws-sdk-go-v2/internal/auth/smithy"
	internalConfig "github.com/aws/aws-sdk-go-v2/internal/configsources"
	internalmiddleware "github.com/aws/aws-sdk-go-v2/internal/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/v4a"
	acceptencodingcust "github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding"
	internalChecksum "github.com/aws/aws-sdk-go-v2/service/internal/checksum"
	presignedurlcust "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url"
	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
	s3sharedconfig "github.com/aws/aws-sdk-go-v2/service/internal/s3shared/config"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	smithy "github.com/aws/smithy-go"
	smithyauth "github.com/aws/smithy-go/auth"
	smithydocument "github.com/aws/smithy-go/document"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/metrics"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/tracing"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const ServiceID = "S3"
const ServiceAPIVersion = "2006-03-01"

type operationMetrics struct {
	Duration                metrics.Float64Histogram
	SerializeDuration       metrics.Float64Histogram
	ResolveIdentityDuration metrics.Float64Histogram
	ResolveEndpointDuration metrics.Float64Histogram
	SignRequestDuration     metrics.Float64Histogram
	DeserializeDuration     metrics.Float64Histogram
}

func (m *operationMetrics) histogramFor(name string) metrics.Float64Histogram {
	switch name {
	case "client.call.duration":
		return m.Duration
	case "client.call.serialization_duration":
		return m.SerializeDuration
	case "client.call.resolve_identity_duration":
		return m.ResolveIdentityDuration
	case "client.call.resolve_endpoint_duration":
		return m.ResolveEndpointDuration
	case "client.call.signing_duration":
		return m.SignRequestDuration
	case "client.call.deserialization_duration":
		return m.DeserializeDuration
	default:
		panic("unrecognized operation metric")
	}
}

func timeOperationMetric[T any](
	ctx context.Context, metric string, fn func() (T, error),
	opts ...metrics.RecordMetricOption,
) (T, error) {
	instr := getOperationMetrics(ctx).histogramFor(metric)
	opts = append([]metrics.RecordMetricOption{withOperationMetadata(ctx)}, opts...)

	start := time.Now()
	v, err := fn()
	end := time.Now()

	elapsed := end.Sub(start)
	instr.Record(ctx, float64(elapsed)/1e9, opts...)
	return v, err
}

func startMetricTimer(ctx context.Context, metric string, opts ...metrics.RecordMetricOption) func() {
	instr := getOperationMetrics(ctx).histogramFor(metric)
	opts = append([]metrics.RecordMetricOption{withOperationMetadata(ctx)}, opts...)

	var ended bool
	start := time.Now()
	return func() {
		if ended {
			return
		}
		ended = true

		end := time.Now()

		elapsed := end.Sub(start)
		instr.Record(ctx, float64(elapsed)/1e9, opts...)
	}
}

func withOperationMetadata(ctx context.Context) metrics.RecordMetricOption {
	return func(o *metrics.RecordMetricOptions) {
		o.Properties.Set("rpc.service", middleware.GetServiceID(ctx))
		o.Properties.Set("rpc.method", middleware.GetOperationName(ctx))
	}
}

type operationMetricsKey struct{}

func withOperationMetrics(parent context.Context, mp metrics.MeterProvider) (context.Context, error) {
	meter := mp.Meter("github.com/aws/aws-sdk-go-v2/service/s3")
	om := &operationMetrics{}

	var err error

	om.Duration, err = operationMetricTimer(meter, "client.call.duration",
		"Overall call duration (including retries and time to send or receive request and response body)")
	if err != nil {
		return nil, err
	}
	om.SerializeDuration, err = operationMetricTimer(meter, "client.call.serialization_duration",
		"The time it takes to serialize a message body")
	if err != nil {
		return nil, err
	}
	om.ResolveIdentityDuration, err = operationMetricTimer(meter, "client.call.auth.resolve_identity_duration",
		"The time taken to acquire an identity (AWS credentials, bearer token, etc) from an Identity Provider")
	if err != nil {
		return nil, err
	}
	om.ResolveEndpointDuration, err = operationMetricTimer(meter, "client.call.resolve_endpoint_duration",
		"The time it takes to resolve an endpoint (endpoint resolver, not DNS) for the request")
	if err != nil {
		return nil, err
	}
	om.SignRequestDuration, err = operationMetricTimer(meter, "client.call.auth.signing_duration",
		"The time it takes to sign a request")
	if err != nil {
		return nil, err
	}
	om.DeserializeDuration, err = operationMetricTimer(meter, "client.call.deserialization_duration",
		"The time it takes to deserialize a message body")
	if err != nil {
		return nil, err
	}

	return context.WithValue(parent, operationMetricsKey{}, om), nil
}

func operationMetricTimer(m metrics.Meter, name, desc string) (metrics.Float64Histogram, error) {
	return m.Float64Histogram(name, func(o *metrics.InstrumentOptions) {
		o.UnitLabel = "s"
		o.Description = desc
	})
}

func getOperationMetrics(ctx context.Context) *operationMetrics {
	return ctx.Value(operationMetricsKey{}).(*operationMetrics)
}

func operationTracer(p tracing.TracerProvider) tracing.Tracer {
	return p.Tracer("github.com/aws/aws-sdk-go-v2/service/s3")
}

// Client provides the API client to make operations call for Amazon Simple
// Storage Service.
type Client struct {
	options Options

	// Difference between the time reported by the server and the client
	timeOffset *atomic.Int64
}

// New returns an initialized Client based on the functional options. Provide
// additional functional options to further configure the behavior of the client,
// such as changing the client's endpoint or adding custom middleware behavior.
func New(options Options, optFns ...func(*Options)) *Client {
	options = options.Copy()

	resolveDefaultLogger(&options)

	setResolvedDefaultsMode(&options)

	resolveRetryer(&options)

	resolveHTTPClient(&options)

	resolveHTTPSignerV4(&options)

	resolveEndpointResolverV2(&options)

	resolveHTTPSignerV4a(&options)

	resolveMeterProvider(&options)

	resolveTracerProvider(&options)

	resolveAuthSchemeResolver(&options)

	for _, fn := range optFns {
		fn(&options)
	}

	finalizeRetryMaxAttempts(&options)

	ignoreAnonymousAuth(&options)

	resolveExpressCredentials(&options)

	finalizeServiceEndpointAuthResolver(&options)

	resolveAuthSchemes(&options)

	client := &Client{
		options: options,
	}

	finalizeExpressCredentials(&options, client)

	initializeTimeOffsetResolver(client)

	return client
}

// Options returns a copy of the client configuration.
//
// Callers SHOULD NOT perform mutations on any inner structures within client
// config. Config overrides should instead be made on a per-operation basis through
// functional options.
func (c *Client) Options() Options {
	return c.options.Copy()
}

func (c *Client) invokeOperation(
	ctx context.Context, opID string, params interface{}, optFns []func(*Options), stackFns ...func(*middleware.Stack, Options) error,
) (
	result interface{}, metadata middleware.Metadata, err error,
) {
	ctx = middleware.ClearStackValues(ctx)
	ctx = middleware.WithServiceID(ctx, ServiceID)
	ctx = middleware.WithOperationName(ctx, opID)

	stack := middleware.NewStack(opID, smithyhttp.NewStackRequest)
	options := c.options.Copy()

	for _, fn := range optFns {
		fn(&options)
	}

	setSafeEventStreamClientLogMode(&options, opID)

	finalizeOperationRetryMaxAttempts(&options, *c)

	finalizeClientEndpointResolverOptions(&options)

	finalizeOperationExpressCredentials(&options, *c)

	finalizeOperationEndpointAuthResolver(&options)

	for _, fn := range stackFns {
		if err := fn(stack, options); err != nil {
			return nil, metadata, err
		}
	}

	for _, fn := range options.APIOptions {
		if err := fn(stack); err != nil {
			return nil, metadata, err
		}
	}

	ctx, err = withOperationMetrics(ctx, options.MeterProvider)
	if err != nil {
		return nil, metadata, err
	}

	tracer := operationTracer(options.TracerProvider)
	spanName := fmt.Sprintf("%s.%s", ServiceID, opID)

	ctx = tracing.WithOperationTracer(ctx, tracer)

	ctx, span := tracer.StartSpan(ctx, spanName, func(o *tracing.SpanOptions) {
		o.Kind = tracing.SpanKindClient
		o.Properties.Set("rpc.system", "aws-api")
		o.Properties.Set("rpc.method", opID)
		o.Properties.Set("rpc.service", ServiceID)
	})
	endTimer := startMetricTimer(ctx, "client.call.duration")
	defer endTimer()
	defer span.End()

	handler := smithyhttp.NewClientHandler(options.HTTPClient)
	decorated := middleware.DecorateHandler(handler, stack)
	result, metadata, err = decorated.Handle(ctx, params)
	if err != nil {
		span.SetProperty("exception.type", fmt.Sprintf("%T", err))
		span.SetProperty("exception.message", err.Error())

		var aerr smithy.APIError
		if errors.As(err, &aerr) {
			span.SetProperty("api.error_code", aerr.ErrorCode())
			span.SetProperty("api.error_message", aerr.ErrorMessage())
			span.SetProperty("api.error_fault", aerr.ErrorFault().String())
		}

		err = &smithy.OperationError{
			ServiceID:     ServiceID,
			OperationName: opID,
			Err:           err,
		}
	}

	span.SetProperty("error", err != nil)
	if err == nil {
		span.SetStatus(tracing.SpanStatusOK)
	} else {
		span.SetStatus(tracing.SpanStatusError)
	}

	return result, metadata, err
}

type operationInputKey struct{}

func setOperationInput(ctx context.Context, input interface{}) context.Context {
	return middleware.WithStackValue(ctx, operationInputKey{}, input)
}

func getOperationInput(ctx context.Context) interface{} {
	return middleware.GetStackValue(ctx, operationInputKey{})
}

type setOperationInputMiddleware struct {
}

func (*setOperationInputMiddleware) ID() string {
	return "setOperationInput"
}

func (m *setOperationInputMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	ctx = setOperationInput(ctx, in.Parameters)
	return next.HandleSerialize(ctx, in)
}

func addProtocolFinalizerMiddlewares(stack *middleware.Stack, options Options, operation string) error {
	if err := stack.Finalize.Add(&resolveAuthSchemeMiddleware{operation: operation, options: options}, middleware.Before); err != nil {
		return fmt.Errorf("add ResolveAuthScheme: %w", err)
	}
	if err := stack.Finalize.Insert(&getIdentityMiddleware{options: options}, "ResolveAuthScheme", middleware.After); err != nil {
		return fmt.Errorf("add GetIdentity: %v", err)
	}
	if err := stack.Finalize.Insert(&resolveEndpointV2Middleware{options: options}, "GetIdentity", middleware.After); err != nil {
		return fmt.Errorf("add ResolveEndpointV2: %v", err)
	}
	if err := stack.Finalize.Insert(&signRequestMiddleware{options: options}, "ResolveEndpointV2", middleware.After); err != nil {
		return fmt.Errorf("add Signing: %w", err)
	}
	return nil
}
func resolveAuthSchemeResolver(options *Options) {
	if options.AuthSchemeResolver == nil {
		options.AuthSchemeResolver = &defaultAuthSchemeResolver{}
	}
}

func resolveAuthSchemes(options *Options) {
	if options.AuthSchemes == nil {
		options.AuthSchemes = []smithyhttp.AuthScheme{
			internalauth.NewHTTPAuthScheme("aws.auth#sigv4", &internalauthsmithy.V4SignerAdapter{
				Signer:     options.HTTPSignerV4,
				Logger:     options.Logger,
				LogSigning: options.ClientLogMode.IsSigning(),
			}),
			internalauth.NewHTTPAuthScheme("com.amazonaws.s3#sigv4express", &s3cust.ExpressSigner{
				Signer:     options.HTTPSignerV4,
				Logger:     options.Logger,
				LogSigning: options.ClientLogMode.IsSigning(),
			}),
			internalauth.NewHTTPAuthScheme("aws.auth#sigv4a", &v4a.SignerAdapter{
				Signer:     options.httpSignerV4a,
				Logger:     options.Logger,
				LogSigning: options.ClientLogMode.IsSigning(),
			}),
		}
	}
}

type noSmithyDocumentSerde = smithydocument.NoSerde

type legacyEndpointContextSetter struct {
	LegacyResolver EndpointResolver
}

func (*legacyEndpointContextSetter) ID() string {
	return "legacyEndpointContextSetter"
}

func (m *legacyEndpointContextSetter) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.LegacyResolver != nil {
		ctx = awsmiddleware.SetRequiresLegacyEndpoints(ctx, true)
	}

	return next.HandleInitialize(ctx, in)

}
func addlegacyEndpointContextSetter(stack *middleware.Stack, o Options) error {
	return stack.Initialize.Add(&legacyEndpointContextSetter{
		LegacyResolver: o.EndpointResolver,
	}, middleware.Before)
}

func resolveDefaultLogger(o *Options) {
	if o.Logger != nil {
		return
	}
	o.Logger = logging.Nop{}
}

func addSetLoggerMiddleware(stack *middleware.Stack, o Options) error {
	return middleware.AddSetLoggerMiddleware(stack, o.Logger)
}

func setResolvedDefaultsMode(o *Options) {
	if len(o.resolvedDefaultsMode) > 0 {
		return
	}

	var mode aws.DefaultsMode
	mode.SetFromString(string(o.DefaultsMode))

	if mode == aws.DefaultsModeAuto {
		mode = defaults.ResolveDefaultsModeAuto(o.Region, o.RuntimeEnvironment)
	}

	o.resolvedDefaultsMode = mode
}

// NewFromConfig returns a new client from the provided config.
func NewFromConfig(cfg aws.Config, optFns ...func(*Options)) *Client {
	opts := Options{
		Region:             cfg.Region,
		DefaultsMode:       cfg.DefaultsMode,
		RuntimeEnvironment: cfg.RuntimeEnvironment,
		HTTPClient:         cfg.HTTPClient,
		Credentials:        cfg.Credentials,
		APIOptions:         cfg.APIOptions,
		Logger:             cfg.Logger,
		ClientLogMode:      cfg.ClientLogMode,
		AppID:              cfg.AppID,
	}
	resolveAWSRetryerProvider(cfg, &opts)
	resolveAWSRetryMaxAttempts(cfg, &opts)
	resolveAWSRetryMode(cfg, &opts)
	resolveAWSEndpointResolver(cfg, &opts)
	resolveUseARNRegion(cfg, &opts)
	resolveDisableMultiRegionAccessPoints(cfg, &opts)
	resolveDisableExpressAuth(cfg, &opts)
	resolveUseDualStackEndpoint(cfg, &opts)
	resolveUseFIPSEndpoint(cfg, &opts)
	resolveBaseEndpoint(cfg, &opts)
	return New(opts, optFns...)
}

func resolveHTTPClient(o *Options) {
	var buildable *awshttp.BuildableClient

	if o.HTTPClient != nil {
		var ok bool
		buildable, ok = o.HTTPClient.(*awshttp.BuildableClient)
		if !ok {
			return
		}
	} else {
		buildable = awshttp.NewBuildableClient()
	}

	modeConfig, err := defaults.GetModeConfiguration(o.resolvedDefaultsMode)
	if err == nil {
		buildable = buildable.WithDialerOptions(func(dialer *net.Dialer) {
			if dialerTimeout, ok := modeConfig.GetConnectTimeout(); ok {
				dialer.Timeout = dialerTimeout
			}
		})

		buildable = buildable.WithTransportOptions(func(transport *http.Transport) {
			if tlsHandshakeTimeout, ok := modeConfig.GetTLSNegotiationTimeout(); ok {
				transport.TLSHandshakeTimeout = tlsHandshakeTimeout
			}
		})
	}

	o.HTTPClient = buildable
}

func resolveRetryer(o *Options) {
	if o.Retryer != nil {
		return
	}

	if len(o.RetryMode) == 0 {
		modeConfig, err := defaults.GetModeConfiguration(o.resolvedDefaultsMode)
		if err == nil {
			o.RetryMode = modeConfig.RetryMode
		}
	}
	if len(o.RetryMode) == 0 {
		o.RetryMode = aws.RetryModeStandard
	}

	var standardOptions []func(*retry.StandardOptions)
	if v := o.RetryMaxAttempts; v != 0 {
		standardOptions = append(standardOptions, func(so *retry.StandardOptions) {
			so.MaxAttempts = v
		})
	}

	switch o.RetryMode {
	case aws.RetryModeAdaptive:
		var adaptiveOptions []func(*retry.AdaptiveModeOptions)
		if len(standardOptions) != 0 {
			adaptiveOptions = append(adaptiveOptions, func(ao *retry.AdaptiveModeOptions) {
				ao.StandardOptions = append(ao.StandardOptions, standardOptions...)
			})
		}
		o.Retryer = retry.NewAdaptiveMode(adaptiveOptions...)

	default:
		o.Retryer = retry.NewStandard(standardOptions...)
	}
}

func resolveAWSRetryerProvider(cfg aws.Config, o *Options) {
	if cfg.Retryer == nil {
		return
	}
	o.Retryer = cfg.Retryer()
}

func resolveAWSRetryMode(cfg aws.Config, o *Options) {
	if len(cfg.RetryMode) == 0 {
		return
	}
	o.RetryMode = cfg.RetryMode
}
func resolveAWSRetryMaxAttempts(cfg aws.Config, o *Options) {
	if cfg.RetryMaxAttempts == 0 {
		return
	}
	o.RetryMaxAttempts = cfg.RetryMaxAttempts
}

func finalizeRetryMaxAttempts(o *Options) {
	if o.RetryMaxAttempts == 0 {
		return
	}

	o.Retryer = retry.AddWithMaxAttempts(o.Retryer, o.RetryMaxAttempts)
}

func finalizeOperationRetryMaxAttempts(o *Options, client Client) {
	if v := o.RetryMaxAttempts; v == 0 || v == client.options.RetryMaxAttempts {
		return
	}

	o.Retryer = retry.AddWithMaxAttempts(o.Retryer, o.RetryMaxAttempts)
}

func resolveAWSEndpointResolver(cfg aws.Config, o *Options) {
	if cfg.EndpointResolver == nil && cfg.EndpointResolverWithOptions == nil {
		return
	}
	o.EndpointResolver = withEndpointResolver(cfg.EndpointResolver, cfg.EndpointResolverWithOptions)
}

func addClientUserAgent(stack *middleware.Stack, options Options) error {
	ua, err := getOrAddRequestUserAgent(stack)
	if err != nil {
		return err
	}

	ua.AddSDKAgentKeyValue(awsmiddleware.APIMetadata, "s3", goModuleVersion)
	if len(options.AppID) > 0 {
		ua.AddSDKAgentKey(awsmiddleware.ApplicationIdentifier, options.AppID)
	}

	return nil
}

func getOrAddRequestUserAgent(stack *middleware.Stack) (*awsmiddleware.RequestUserAgent, error) {
	id := (*awsmiddleware.RequestUserAgent)(nil).ID()
	mw, ok := stack.Build.Get(id)
	if !ok {
		mw = awsmiddleware.NewRequestUserAgent()
		if err := stack.Build.Add(mw, middleware.After); err != nil {
			return nil, err
		}
	}

	ua, ok := mw.(*awsmiddleware.RequestUserAgent)
	if !ok {
		return nil, fmt.Errorf("%T for %s middleware did not match expected type", mw, id)
	}

	return ua, nil
}

type HTTPSignerV4 interface {
	SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, signingTime time.Time, optFns ...func(*v4.SignerOptions)) error
}

func resolveHTTPSignerV4(o *Options) {
	if o.HTTPSignerV4 != nil {
		return
	}
	o.HTTPSignerV4 = newDefaultV4Signer(*o)
}

func newDefaultV4Signer(o Options) *v4.Signer {
	return v4.NewSigner(func(so *v4.SignerOptions) {
		so.Logger = o.Logger
		so.LogSigning = o.ClientLogMode.IsSigning()
		so.DisableURIPathEscaping = true
	})
}

func addClientRequestID(stack *middleware.Stack) error {
	return stack.Build.Add(&awsmiddleware.ClientRequestID{}, middleware.After)
}

func addComputeContentLength(stack *middleware.Stack) error {
	return stack.Build.Add(&smithyhttp.ComputeContentLength{}, middleware.After)
}

func addRawResponseToMetadata(stack *middleware.Stack) error {
	return stack.Deserialize.Add(&awsmiddleware.AddRawResponse{}, middleware.Before)
}

func addRecordResponseTiming(stack *middleware.Stack) error {
	return stack.Deserialize.Add(&awsmiddleware.RecordResponseTiming{}, middleware.After)
}

func addSpanRetryLoop(stack *middleware.Stack, options Options) error {
	return stack.Finalize.Insert(&spanRetryLoop{options: options}, "Retry", middleware.Before)
}

type spanRetryLoop struct {
	options Options
}

func (*spanRetryLoop) ID() string {
	return "spanRetryLoop"
}

func (m *spanRetryLoop) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	middleware.FinalizeOutput, middleware.Metadata, error,
) {
	tracer := operationTracer(m.options.TracerProvider)
	ctx, span := tracer.StartSpan(ctx, "RetryLoop")
	defer span.End()

	return next.HandleFinalize(ctx, in)
}
func addStreamingEventsPayload(stack *middleware.Stack) error {
	return stack.Finalize.Add(&v4.StreamingEventsPayload{}, middleware.Before)
}

func addUnsignedPayload(stack *middleware.Stack) error {
	return stack.Finalize.Insert(&v4.UnsignedPayload{}, "ResolveEndpointV2", middleware.After)
}

func addComputePayloadSHA256(stack *middleware.Stack) error {
	return stack.Finalize.Insert(&v4.ComputePayloadSHA256{}, "ResolveEndpointV2", middleware.After)
}

func addContentSHA256Header(stack *middleware.Stack) error {
	return stack.Finalize.Insert(&v4.ContentSHA256Header{}, (*v4.ComputePayloadSHA256)(nil).ID(), middleware.After)
}

func addIsWaiterUserAgent(o *Options) {
	o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
		ua, err := getOrAddRequestUserAgent(stack)
		if err != nil {
			return err
		}

		ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureWaiter)
		return nil
	})
}

func addIsPaginatorUserAgent(o *Options) {
	o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
		ua, err := getOrAddRequestUserAgent(stack)
		if err != nil {
			return err
		}

		ua.AddUserAgentFeature(awsmiddleware.UserAgentFeaturePaginator)
		return nil
	})
}

func addRetry(stack *middleware.Stack, o Options) error {
	attempt := retry.NewAttemptMiddleware(o.Retryer, smithyhttp.RequestCloner, func(m *retry.Attempt) {
		m.LogAttempts = o.ClientLogMode.IsRetries()
		m.OperationMeter = o.MeterProvider.Meter("github.com/aws/aws-sdk-go-v2/service/s3")
	})
	if err := stack.Finalize.Insert(attempt, "Signing", middleware.Before); err != nil {
		return err
	}
	if err := stack.Finalize.Insert(&retry.MetricsHeader{}, attempt.ID(), middleware.After); err != nil {
		return err
	}
	return nil
}

// resolves UseARNRegion S3 configuration
func resolveUseARNRegion(cfg aws.Config, o *Options) error {
	if len(cfg.ConfigSources) == 0 {
		return nil
	}
	value, found, err := s3sharedconfig.ResolveUseARNRegion(context.Background(), cfg.ConfigSources)
	if err != nil {
		return err
	}
	if found {
		o.UseARNRegion = value
	}
	return nil
}

// resolves DisableMultiRegionAccessPoints S3 configuration
func resolveDisableMultiRegionAccessPoints(cfg aws.Config, o *Options) error {
	if len(cfg.ConfigSources) == 0 {
		return nil
	}
	value, found, err := s3sharedconfig.ResolveDisableMultiRegionAccessPoints(context.Background(), cfg.ConfigSources)
	if err != nil {
		return err
	}
	if found {
		o.DisableMultiRegionAccessPoints = value
	}
	return nil
}

// resolves dual-stack endpoint configuration
func resolveUseDualStackEndpoint(cfg aws.Config, o *Options) error {
	if len(cfg.ConfigSources) == 0 {
		return nil
	}
	value, found, err := internalConfig.ResolveUseDualStackEndpoint(context.Background(), cfg.ConfigSources)
	if err != nil {
		return err
	}
	if found {
		o.EndpointOptions.UseDualStackEndpoint = value
	}
	return nil
}

// resolves FIPS endpoint configuration
func resolveUseFIPSEndpoint(cfg aws.Config, o *Options) error {
	if len(cfg.ConfigSources) == 0 {
		return nil
	}
	value, found, err := internalConfig.ResolveUseFIPSEndpoint(context.Background(), cfg.ConfigSources)
	if err != nil {
		return err
	}
	if found {
		o.EndpointOptions.UseFIPSEndpoint = value
	}
	return nil
}

func resolveAccountID(identity smithyauth.Identity, mode aws.AccountIDEndpointMode) *string {
	if mode == aws.AccountIDEndpointModeDisabled {
		return nil
	}

	if ca, ok := identity.(*internalauthsmithy.CredentialsAdapter); ok && ca.Credentials.AccountID != "" {
		return aws.String(ca.Credentials.AccountID)
	}

	return nil
}

type httpSignerV4a interface {
	SignHTTP(ctx context.Context, credentials v4a.Credentials, r *http.Request, payloadHash,
		service string, regionSet []string, signingTime time.Time,
		optFns ...func(*v4a.SignerOptions)) error
}

func resolveHTTPSignerV4a(o *Options) {
	if o.httpSignerV4a != nil {
		return
	}
	o.httpSignerV4a = newDefaultV4aSigner(*o)
}

func newDefaultV4aSigner(o Options) *v4a.Signer {
	return v4a.NewSigner(func(so *v4a.SignerOptions) {
		so.Logger = o.Logger
		so.LogSigning = o.ClientLogMode.IsSigning()
	})
}

func addTimeOffsetBuild(stack *middleware.Stack, c *Client) error {
	mw := internalmiddleware.AddTimeOffsetMiddleware{Offset: c.timeOffset}
	if err := stack.Build.Add(&mw, middleware.After); err != nil {
		return err
	}
	return stack.Deserialize.Insert(&mw, "RecordResponseTiming", middleware.Before)
}
func initializeTimeOffsetResolver(c *Client) {
	c.timeOffset = new(atomic.Int64)
}

func addUserAgentRetryMode(stack *middleware.Stack, options Options) error {
	ua, err := getOrAddRequestUserAgent(stack)
	if err != nil {
		return err
	}

	switch options.Retryer.(type) {
	case *retry.Standard:
		ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureRetryModeStandard)
	case *retry.AdaptiveMode:
		ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureRetryModeAdaptive)
	}
	return nil
}

func resolveTracerProvider(options *Options) {
	if options.TracerProvider == nil {
		options.TracerProvider = &tracing.NopTracerProvider{}
	}
}

func resolveMeterProvider(options *Options) {
	if options.MeterProvider == nil {
		options.MeterProvider = metrics.NopMeterProvider{}
	}
}

func addMetadataRetrieverMiddleware(stack *middleware.Stack) error {
	return s3shared.AddMetadataRetrieverMiddleware(stack)
}

func add100Continue(stack *middleware.Stack, options Options) error {
	return s3shared.Add100Continue(stack, options.ContinueHeaderThresholdBytes)
}

func addRecursionDetection(stack *middleware.Stack) error {
	return stack.Build.Add(&awsmiddleware.RecursionDetection{}, middleware.After)
}

// ComputedInputChecksumsMetadata provides information about the algorithms used
// to compute the checksum(s) of the input payload.
type ComputedInputChecksumsMetadata struct {
	// ComputedChecksums is a map of algorithm name to checksum value of the computed
	// input payload's checksums.
	ComputedChecksums map[string]string
}

// GetComputedInputChecksumsMetadata retrieves from the result metadata the map of
// algorithms and input payload checksums values.
func GetComputedInputChecksumsMetadata(m middleware.Metadata) (ComputedInputChecksumsMetadata, bool) {
	values, ok := internalChecksum.GetComputedInputChecksums(m)
	if !ok {
		return ComputedInputChecksumsMetadata{}, false
	}
	return ComputedInputChecksumsMetadata{
		ComputedChecksums: values,
	}, true

}

// ChecksumValidationMetadata contains metadata such as the checksum algorithm
// used for data integrity validation.
type ChecksumValidationMetadata struct {
	// AlgorithmsUsed is the set of the checksum algorithms used to validate the
	// response payload. The response payload must be completely read in order for the
	// checksum validation to be performed. An error is returned by the operation
	// output's response io.ReadCloser if the computed checksums are invalid.
	AlgorithmsUsed []string
}

// GetChecksumValidationMetadata returns the set of algorithms that will be used
// to validate the response payload with. The response payload must be completely
// read in order for the checksum validation to be performed. An error is returned
// by the operation output's response io.ReadCloser if the computed checksums are
// invalid. Returns false if no checksum algorithm used metadata was found.
func GetChecksumValidationMetadata(m middleware.Metadata) (ChecksumValidationMetadata, bool) {
	values, ok := internalChecksum.GetOutputValidationAlgorithmsUsed(m)
	if !ok {
		return ChecksumValidationMetadata{}, false
	}
	return ChecksumValidationMetadata{
		AlgorithmsUsed: append(make([]string, 0, len(values)), values...),
	}, true

}

// nopGetBucketAccessor is no-op accessor for operation that don't support bucket
// member as input
func nopGetBucketAccessor(input interface{}) (*string, bool) {
	return nil, false
}

func addResponseErrorMiddleware(stack *middleware.Stack) error {
	return s3shared.AddResponseErrorMiddleware(stack)
}

func disableAcceptEncodingGzip(stack *middleware.Stack) error {
	return acceptencodingcust.AddAcceptEncodingGzip(stack, acceptencodingcust.AddAcceptEncodingGzipOptions{})
}

// ResponseError provides the HTTP centric error type wrapping the underlying
// error with the HTTP response value and the deserialized RequestID.
type ResponseError interface {
	error

	ServiceHostID() string
	ServiceRequestID() string
}

var _ ResponseError = (*s3shared.ResponseError)(nil)

// GetHostIDMetadata retrieves the host id from middleware metadata returns host
// id as string along with a boolean indicating presence of hostId on middleware
// metadata.
func GetHostIDMetadata(metadata middleware.Metadata) (string, bool) {
	return s3shared.GetHostIDMetadata(metadata)
}

// HTTPPresignerV4 represents presigner interface used by presign url client
type HTTPPresignerV4 interface {
	PresignHTTP(
		ctx context.Context, credentials aws.Credentials, r *http.Request,
		payloadHash string, service string, region string, signingTime time.Time,
		optFns ...func(*v4.SignerOptions),
	) (url string, signedHeader http.Header, err error)
}

// httpPresignerV4a represents sigv4a presigner interface used by presign url
// client
type httpPresignerV4a interface {
	PresignHTTP(
		ctx context.Context, credentials v4a.Credentials, r *http.Request,
		payloadHash string, service string, regionSet []string, signingTime time.Time,
		optFns ...func(*v4a.SignerOptions),
	) (url string, signedHeader http.Header, err error)
}

// PresignOptions represents the presign client options
type PresignOptions struct {

	// ClientOptions are list of functional options to mutate client options used by
	// the presign client.
	ClientOptions []func(*Options)

	// Presigner is the presigner used by the presign url client
	Presigner HTTPPresignerV4

	// Expires sets the expiration duration for the generated presign url. This should
	// be the duration in seconds the presigned URL should be considered valid for. If
	// not set or set to zero, presign url would default to expire after 900 seconds.
	Expires time.Duration

	// presignerV4a is the presigner used by the presign url client
	presignerV4a httpPresignerV4a
}

func (o PresignOptions) copy() PresignOptions {
	clientOptions := make([]func(*Options), len(o.ClientOptions))
	copy(clientOptions, o.ClientOptions)
	o.ClientOptions = clientOptions
	return o
}

// WithPresignClientFromClientOptions is a helper utility to retrieve a function
// that takes PresignOption as input
func WithPresignClientFromClientOptions(optFns ...func(*Options)) func(*PresignOptions) {
	return withPresignClientFromClientOptions(optFns).options
}

type withPresignClientFromClientOptions []func(*Options)

func (w withPresignClientFromClientOptions) options(o *PresignOptions) {
	o.ClientOptions = append(o.ClientOptions, w...)
}

// WithPresignExpires is a helper utility to append Expires value on presign
// options optional function
func WithPresignExpires(dur time.Duration) func(*PresignOptions) {
	return withPresignExpires(dur).options
}

type withPresignExpires time.Duration

func (w withPresignExpires) options(o *PresignOptions) {
	o.Expires = time.Duration(w)
}

// PresignClient represents the presign url client
type PresignClient struct {
	client  *Client
	options PresignOptions
}

// NewPresignClient generates a presign client using provided API Client and
// presign options
func NewPresignClient(c *Client, optFns ...func(*PresignOptions)) *PresignClient {
	var options PresignOptions
	for _, fn := range optFns {
		fn(&options)
	}
	if len(options.ClientOptions) != 0 {
		c = New(c.options, options.ClientOptions...)
	}

	if options.Presigner == nil {
		options.Presigner = newDefaultV4Signer(c.options)
	}

	if options.presignerV4a == nil {
		options.presignerV4a = newDefaultV4aSigner(c.options)
	}

	return &PresignClient{
		client:  c,
		options: options,
	}
}

func withNopHTTPClientAPIOption(o *Options) {
	o.HTTPClient = smithyhttp.NopClient{}
}

type presignContextPolyfillMiddleware struct {
}

func (*presignContextPolyfillMiddleware) ID() string {
	return "presignContextPolyfill"
}

func (m *presignContextPolyfillMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	rscheme := getResolvedAuthScheme(ctx)
	if rscheme == nil {
		return out, metadata, fmt.Errorf("no resolved auth scheme")
	}

	schemeID := rscheme.Scheme.SchemeID()
	ctx = s3cust.SetSignerVersion(ctx, schemeID)
	if schemeID == "aws.auth#sigv4" || schemeID == "com.amazonaws.s3#sigv4express" {
		if sn, ok := smithyhttp.GetSigV4SigningName(&rscheme.SignerProperties); ok {
			ctx = awsmiddleware.SetSigningName(ctx, sn)
		}
		if sr, ok := smithyhttp.GetSigV4SigningRegion(&rscheme.SignerProperties); ok {
			ctx = awsmiddleware.SetSigningRegion(ctx, sr)
		}
	} else if schemeID == "aws.auth#sigv4a" {
		if sn, ok := smithyhttp.GetSigV4ASigningName(&rscheme.SignerProperties); ok {
			ctx = awsmiddleware.SetSigningName(ctx, sn)
		}
		if sr, ok := smithyhttp.GetSigV4ASigningRegions(&rscheme.SignerProperties); ok {
			ctx = awsmiddleware.SetSigningRegion(ctx, sr[0])
		}
	}

	return next.HandleFinalize(ctx, in)
}

type presignConverter PresignOptions

func (c presignConverter) convertToPresignMiddleware(stack *middleware.Stack, options Options) (err error) {
	if _, ok := stack.Finalize.Get((*acceptencodingcust.DisableGzip)(nil).ID()); ok {
		stack.Finalize.Remove((*acceptencodingcust.DisableGzip)(nil).ID())
	}
	if _, ok := stack.Finalize.Get((*retry.Attempt)(nil).ID()); ok {
		stack.Finalize.Remove((*retry.Attempt)(nil).ID())
	}
	if _, ok := stack.Finalize.Get((*retry.MetricsHeader)(nil).ID()); ok {
		stack.Finalize.Remove((*retry.MetricsHeader)(nil).ID())
	}
	stack.Deserialize.Clear()
	stack.Build.Remove((*awsmiddleware.ClientRequestID)(nil).ID())
	stack.Build.Remove("UserAgent")
	if err := stack.Finalize.Insert(&presignContextPolyfillMiddleware{}, "Signing", middleware.Before); err != nil {
		return err
	}

	pmw := v4.NewPresignHTTPRequestMiddleware(v4.PresignHTTPRequestMiddlewareOptions{
		CredentialsProvider: options.Credentials,
		Presigner:           c.Presigner,
		LogSigning:          options.ClientLogMode.IsSigning(),
	})
	if _, err := stack.Finalize.Swap("Signing", pmw); err != nil {
		return err
	}
	if err = smithyhttp.AddNoPayloadDefaultContentTypeRemover(stack); err != nil {
		return err
	}

	// extended s3 presigning
	signermv := s3cust.NewPresignHTTPRequestMiddleware(s3cust.PresignHTTPRequestMiddlewareOptions{
		CredentialsProvider: options.Credentials,
		ExpressCredentials:  options.ExpressCredentials,
		V4Presigner:         c.Presigner,
		V4aPresigner:        c.presignerV4a,
		LogSigning:          options.ClientLogMode.IsSigning(),
	})
	err = s3cust.RegisterPreSigningMiddleware(stack, signermv)
	if err != nil {
		return err
	}

	if c.Expires < 0 {
		return fmt.Errorf("presign URL duration must be 0 or greater, %v", c.Expires)
	}
	// add middleware to set expiration for s3 presigned url, if expiration is set to
	// 0, this middleware sets a default expiration of 900 seconds
	err = stack.Build.Add(&s3cust.AddExpiresOnPresignedURL{Expires: c.Expires}, middleware.After)
	if err != nil {
		return err
	}
	err = presignedurlcust.AddAsIsPresigningMiddleware(stack)
	if err != nil {
		return err
	}
	return nil
}

func addRequestResponseLogging(stack *middleware.Stack, o Options) error {
	return stack.Deserialize.Add(&smithyhttp.RequestResponseLogger{
		LogRequest:          o.ClientLogMode.IsRequest(),
		LogRequestWithBody:  o.ClientLogMode.IsRequestWithBody(),
		LogResponse:         o.ClientLogMode.IsResponse(),
		LogResponseWithBody: o.ClientLogMode.IsResponseWithBody(),
	}, middleware.After)
}

type disableHTTPSMiddleware struct {
	DisableHTTPS bool
}

func (*disableHTTPSMiddleware) ID() string {
	return "disableHTTPS"
}

func (m *disableHTTPSMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.DisableHTTPS && !smithyhttp.GetHostnameImmutable(ctx) {
		req.URL.Scheme = "http"
	}

	return next.HandleFinalize(ctx, in)
}

func addDisableHTTPSMiddleware(stack *middleware.Stack, o Options) error {
	return stack.Finalize.Insert(&disableHTTPSMiddleware{
		DisableHTTPS: o.EndpointOptions.DisableHTTPS,
	}, "ResolveEndpointV2", middleware.After)
}

type spanInitializeStart struct {
}

func (*spanInitializeStart) ID() string {
	return "spanInitializeStart"
}

func (m *spanInitializeStart) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	ctx, _ = tracing.StartSpan(ctx, "Initialize")

	return next.HandleInitialize(ctx, in)
}

type spanInitializeEnd struct {
}

func (*spanInitializeEnd) ID() string {
	return "spanInitializeEnd"
}

func (m *spanInitializeEnd) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	middleware.InitializeOutput, middleware.Metadata, error,
) {
	ctx, span := tracing.PopSpan(ctx)
	span.End()

	return next.HandleInitialize(ctx, in)
}

type spanBuildRequestStart struct {
}

func (*spanBuildRequestStart) ID() string {
	return "spanBuildRequestStart"
}

func (m *spanBuildRequestStart) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	middleware.SerializeOutput, middleware.Metadata, error,
) {
	ctx, _ = tracing.StartSpan(ctx, "BuildRequest")

	return next.HandleSerialize(ctx, in)
}

type spanBuildRequestEnd struct {
}

func (*spanBuildRequestEnd) ID() string {
	return "spanBuildRequestEnd"
}

func (m *spanBuildRequestEnd) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	middleware.BuildOutput, middleware.Metadata, error,
) {
	ctx, span := tracing.PopSpan(ctx)
	span.End()

	return next.HandleBuild(ctx, in)
}

func addSpanInitializeStart(stack *middleware.Stack) error {
	return stack.Initialize.Add(&spanInitializeStart{}, middleware.Before)
}

func addSpanInitializeEnd(stack *middleware.Stack) error {
	return stack.Initialize.Add(&spanInitializeEnd{}, middleware.After)
}

func addSpanBuildRequestStart(stack *middleware.Stack) error {
	return stack.Serialize.Add(&spanBuildRequestStart{}, middleware.Before)
}

func addSpanBuildRequestEnd(stack *middleware.Stack) error {
	return stack.Build.Add(&spanBuildRequestEnd{}, middleware.After)
}
