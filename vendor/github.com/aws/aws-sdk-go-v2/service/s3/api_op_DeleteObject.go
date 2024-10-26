// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"fmt"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Removes an object from a bucket. The behavior depends on the bucket's
// versioning state:
//
//   - If bucket versioning is not enabled, the operation permanently deletes the
//     object.
//
//   - If bucket versioning is enabled, the operation inserts a delete marker,
//     which becomes the current version of the object. To permanently delete an object
//     in a versioned bucket, you must include the object’s versionId in the request.
//     For more information about versioning-enabled buckets, see [Deleting object versions from a versioning-enabled bucket].
//
//   - If bucket versioning is suspended, the operation removes the object that
//     has a null versionId , if there is one, and inserts a delete marker that
//     becomes the current version of the object. If there isn't an object with a null
//     versionId , and all versions of the object have a versionId , Amazon S3 does
//     not remove the object and only inserts a delete marker. To permanently delete an
//     object that has a versionId , you must include the object’s versionId in the
//     request. For more information about versioning-suspended buckets, see [Deleting objects from versioning-suspended buckets].
//
//   - Directory buckets - S3 Versioning isn't enabled and supported for directory
//     buckets. For this API operation, only the null value of the version ID is
//     supported by directory buckets. You can only specify null to the versionId
//     query parameter in the request.
//
//   - Directory buckets - For directory buckets, you must make requests for this
//     API operation to the Zonal endpoint. These endpoints support
//     virtual-hosted-style requests in the format
//     https://bucket_name.s3express-az_id.region.amazonaws.com/key-name .
//     Path-style requests are not supported. For more information, see [Regional and Zonal endpoints]in the
//     Amazon S3 User Guide.
//
// To remove a specific version, you must use the versionId query parameter. Using
// this query parameter permanently deletes the version. If the object deleted is a
// delete marker, Amazon S3 sets the response header x-amz-delete-marker to true.
//
// If the object you want to delete is in a bucket where the bucket versioning
// configuration is MFA Delete enabled, you must include the x-amz-mfa request
// header in the DELETE versionId request. Requests that include x-amz-mfa must
// use HTTPS. For more information about MFA Delete, see [Using MFA Delete]in the Amazon S3 User
// Guide. To see sample requests that use versioning, see [Sample Request].
//
// Directory buckets - MFA delete is not supported by directory buckets.
//
// You can delete objects by explicitly calling DELETE Object or calling ([PutBucketLifecycle] ) to
// enable Amazon S3 to remove them for you. If you want to block users or accounts
// from removing or deleting objects from your bucket, you must deny them the
// s3:DeleteObject , s3:DeleteObjectVersion , and s3:PutLifeCycleConfiguration
// actions.
//
// Directory buckets - S3 Lifecycle is not supported by directory buckets.
//
// Permissions
//
//   - General purpose bucket permissions - The following permissions are required
//     in your policies when your DeleteObjects request includes specific headers.
//
//   - s3:DeleteObject - To delete an object from a bucket, you must always have
//     the s3:DeleteObject permission.
//
//   - s3:DeleteObjectVersion - To delete a specific version of an object from a
//     versioning-enabled bucket, you must have the s3:DeleteObjectVersion permission.
//
//   - Directory bucket permissions - To grant access to this API operation on a
//     directory bucket, we recommend that you use the [CreateSession]CreateSession API operation
//     for session-based authorization. Specifically, you grant the
//     s3express:CreateSession permission to the directory bucket in a bucket policy
//     or an IAM identity-based policy. Then, you make the CreateSession API call on
//     the bucket to obtain a session token. With the session token in your request
//     header, you can make API requests to this operation. After the session token
//     expires, you make another CreateSession API call to generate a new session
//     token for use. Amazon Web Services CLI or SDKs create session and refresh the
//     session token automatically to avoid service interruptions when a session
//     expires. For more information about authorization, see [CreateSession]CreateSession .
//
// HTTP Host header syntax  Directory buckets - The HTTP Host header syntax is
// Bucket_name.s3express-az_id.region.amazonaws.com .
//
// The following action is related to DeleteObject :
//
// [PutObject]
//
// [Sample Request]: https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectDELETE.html#ExampleVersionObjectDelete
// [Regional and Zonal endpoints]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-express-Regions-and-Zones.html
// [Deleting objects from versioning-suspended buckets]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/DeletingObjectsfromVersioningSuspendedBuckets.html
// [PutObject]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
// [PutBucketLifecycle]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketLifecycle.html
// [CreateSession]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateSession.html
// [Deleting object versions from a versioning-enabled bucket]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/DeletingObjectVersions.html
// [Using MFA Delete]: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMFADelete.html
func (c *Client) DeleteObject(ctx context.Context, params *DeleteObjectInput, optFns ...func(*Options)) (*DeleteObjectOutput, error) {
	if params == nil {
		params = &DeleteObjectInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DeleteObject", params, optFns, c.addOperationDeleteObjectMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DeleteObjectOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DeleteObjectInput struct {

	// The bucket name of the bucket containing the object.
	//
	// Directory buckets - When you use this operation with a directory bucket, you
	// must use virtual-hosted-style requests in the format
	// Bucket_name.s3express-az_id.region.amazonaws.com . Path-style requests are not
	// supported. Directory bucket names must be unique in the chosen Availability
	// Zone. Bucket names must follow the format bucket_base_name--az-id--x-s3 (for
	// example, DOC-EXAMPLE-BUCKET--usw2-az1--x-s3 ). For information about bucket
	// naming restrictions, see [Directory bucket naming rules]in the Amazon S3 User Guide.
	//
	// Access points - When you use this action with an access point, you must provide
	// the alias of the access point in place of the bucket name or specify the access
	// point ARN. When using the access point ARN, you must direct requests to the
	// access point hostname. The access point hostname takes the form
	// AccessPointName-AccountId.s3-accesspoint.Region.amazonaws.com. When using this
	// action with an access point through the Amazon Web Services SDKs, you provide
	// the access point ARN in place of the bucket name. For more information about
	// access point ARNs, see [Using access points]in the Amazon S3 User Guide.
	//
	// Access points and Object Lambda access points are not supported by directory
	// buckets.
	//
	// S3 on Outposts - When you use this action with Amazon S3 on Outposts, you must
	// direct requests to the S3 on Outposts hostname. The S3 on Outposts hostname
	// takes the form
	// AccessPointName-AccountId.outpostID.s3-outposts.Region.amazonaws.com . When you
	// use this action with S3 on Outposts through the Amazon Web Services SDKs, you
	// provide the Outposts access point ARN in place of the bucket name. For more
	// information about S3 on Outposts ARNs, see [What is S3 on Outposts?]in the Amazon S3 User Guide.
	//
	// [Directory bucket naming rules]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/directory-bucket-naming-rules.html
	// [What is S3 on Outposts?]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html
	// [Using access points]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-access-points.html
	//
	// This member is required.
	Bucket *string

	// Key name of the object to delete.
	//
	// This member is required.
	Key *string

	// Indicates whether S3 Object Lock should bypass Governance-mode restrictions to
	// process this operation. To use this header, you must have the
	// s3:BypassGovernanceRetention permission.
	//
	// This functionality is not supported for directory buckets.
	BypassGovernanceRetention *bool

	// The account ID of the expected bucket owner. If the account ID that you provide
	// does not match the actual owner of the bucket, the request fails with the HTTP
	// status code 403 Forbidden (access denied).
	ExpectedBucketOwner *string

	// The concatenation of the authentication device's serial number, a space, and
	// the value that is displayed on your authentication device. Required to
	// permanently delete a versioned object if versioning is configured with MFA
	// delete enabled.
	//
	// This functionality is not supported for directory buckets.
	MFA *string

	// Confirms that the requester knows that they will be charged for the request.
	// Bucket owners need not specify this parameter in their requests. If either the
	// source or destination S3 bucket has Requester Pays enabled, the requester will
	// pay for corresponding charges to copy the object. For information about
	// downloading objects from Requester Pays buckets, see [Downloading Objects in Requester Pays Buckets]in the Amazon S3 User
	// Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Downloading Objects in Requester Pays Buckets]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
	RequestPayer types.RequestPayer

	// Version ID used to reference a specific version of the object.
	//
	// For directory buckets in this API operation, only the null value of the version
	// ID is supported.
	VersionId *string

	noSmithyDocumentSerde
}

func (in *DeleteObjectInput) bindEndpointParams(p *EndpointParameters) {

	p.Bucket = in.Bucket
	p.Key = in.Key

}

type DeleteObjectOutput struct {

	// Indicates whether the specified object version that was permanently deleted was
	// (true) or was not (false) a delete marker before deletion. In a simple DELETE,
	// this header indicates whether (true) or not (false) the current version of the
	// object is a delete marker.
	//
	// This functionality is not supported for directory buckets.
	DeleteMarker *bool

	// If present, indicates that the requester was successfully charged for the
	// request.
	//
	// This functionality is not supported for directory buckets.
	RequestCharged types.RequestCharged

	// Returns the version ID of the delete marker created as a result of the DELETE
	// operation.
	//
	// This functionality is not supported for directory buckets.
	VersionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDeleteObjectMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpDeleteObject{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpDeleteObject{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DeleteObject"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addPutBucketContextMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIsExpressUserAgent(stack); err != nil {
		return err
	}
	if err = addOpDeleteObjectValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDeleteObject(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addDeleteObjectUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func (v *DeleteObjectInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opDeleteObject(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DeleteObject",
	}
}

// getDeleteObjectBucketMember returns a pointer to string denoting a provided
// bucket member valueand a boolean indicating if the input has a modeled bucket
// name,
func getDeleteObjectBucketMember(input interface{}) (*string, bool) {
	in := input.(*DeleteObjectInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addDeleteObjectUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getDeleteObjectBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}

// PresignDeleteObject is used to generate a presigned HTTP Request which contains
// presigned URL, signed headers and HTTP method used.
func (c *PresignClient) PresignDeleteObject(ctx context.Context, params *DeleteObjectInput, optFns ...func(*PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	if params == nil {
		params = &DeleteObjectInput{}
	}
	options := c.options.copy()
	for _, fn := range optFns {
		fn(&options)
	}
	clientOptFns := append(options.ClientOptions, withNopHTTPClientAPIOption)

	result, _, err := c.client.invokeOperation(ctx, "DeleteObject", params, clientOptFns,
		c.client.addOperationDeleteObjectMiddlewares,
		presignConverter(options).convertToPresignMiddleware,
		addDeleteObjectPayloadAsUnsigned,
	)
	if err != nil {
		return nil, err
	}

	out := result.(*v4.PresignedHTTPRequest)
	return out, nil
}

func addDeleteObjectPayloadAsUnsigned(stack *middleware.Stack, options Options) error {
	v4.RemoveContentSHA256HeaderMiddleware(stack)
	v4.RemoveComputePayloadSHA256Middleware(stack)
	return v4.AddUnsignedPayloadMiddleware(stack)
}
