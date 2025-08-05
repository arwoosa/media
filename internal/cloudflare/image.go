package cloudflare

import (
	"context"
	"fmt"
	"time"

	"github.com/arwoosa/media/internal/cloudflare/dao"
	cloudflare "github.com/cloudflare/cloudflare-go/v4"
	images "github.com/cloudflare/cloudflare-go/v4/images"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/spf13/viper"
)

var accountID string
var apiToken string
var expiryDuration time.Duration

func initialByViper() {
	accountID = viper.Get("cloudflare.account_id").(string)
	apiToken = viper.Get("cloudflare.api_token").(string)
	expiryDuration = viper.GetDuration("cloudflare.expiry_duration")
}

func checkConfig() error {
	if accountID == "" || apiToken == "" || expiryDuration == 0 {
		initialByViper()
	}
	if accountID == "" || apiToken == "" || expiryDuration == 0 {
		return fmt.Errorf("%w: check env variables [cloudflare.account_id, cloudflare.api_token, cloudflare.expiry_duration]", ErrCloudflareConfigNotInitialized)
	}
	return nil
}

type imageMetadataOption func(*dao.ImageMetadata)

func ImageMetadataSize(size uint64) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Size = size
	}
}

func ImageMetadataWidth(width uint32) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Width = width
	}
}

func ImageMetadataHeight(height uint32) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Height = height
	}
}

func ImageMetadataFormat(format string) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Format = format
	}
}

func ImageMetadataLatitude(latitude *float64) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Latitude = latitude
	}
}

func ImageMetadataLongitude(longitude *float64) imageMetadataOption {
	return func(m *dao.ImageMetadata) {
		m.Longitude = longitude
	}
}

func GetSignedUrl(ctx context.Context, opts ...imageMetadataOption) (*struct {
	UploadURL string
	ID        string
}, error) {
	if err := checkConfig(); err != nil {
		return nil, err
	}

	metadata := &dao.ImageMetadata{}
	for _, opt := range opts {
		opt(metadata)
	}
	service := images.NewV2DirectUploadService(
		option.WithAPIToken(apiToken),
		option.WithEnvironmentProduction())
	resp, err := service.New(ctx, images.V2DirectUploadNewParams{
		AccountID:         cloudflare.F(accountID),
		RequireSignedURLs: cloudflare.F(false),
		Expiry:            cloudflare.F(time.Now().Add(expiryDuration)),
		Metadata:          cloudflare.F(metadata.ToCoudflareFieldMetadata()),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
	}
	return &struct {
		UploadURL string
		ID        string
	}{
		UploadURL: resp.UploadURL,
		ID:        resp.ID,
	}, nil
}

func GetImageDetail(ctx context.Context, id string) (*dao.Image, error) {
	if err := checkConfig(); err != nil {
		return nil, err
	}
	service := images.NewV1Service(
		option.WithAPIToken(apiToken),
		option.WithEnvironmentProduction(),
	)
	resp, err := service.Get(ctx, id, images.V1GetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
	}
	if resp.Meta == nil {
		resp.Meta = map[string]string{}
	}
	metadata := resp.Meta.(map[string]any)
	meta := map[string]string{}
	for k, v := range metadata {
		meta[k] = v.(string)
	}
	return &dao.Image{
		ID:       id,
		Filename: resp.Filename,
		Uploaded: resp.Uploaded,
		Meta:     meta,
		Variants: resp.Variants,
	}, nil
}

func GetImages(ctx context.Context, ids []string) map[string]struct {
	Image *dao.Image
	Err   error
} {
	if err := checkConfig(); err != nil {
		return nil
	}
	result := make(map[string]struct {
		Image *dao.Image
		Err   error
	})
	for _, id := range ids {
		image, err := GetImageDetail(ctx, id)
		result[id] = struct {
			Image *dao.Image
			Err   error
		}{
			Image: image,
			Err:   err,
		}
	}
	return result
}

func DeleteImages(ctx context.Context, id ...string) error {
	if err := checkConfig(); err != nil {
		return err
	}
	service := images.NewV1Service(
		option.WithAPIToken(apiToken),
		option.WithEnvironmentProduction(),
	)
	for _, id := range id {
		_, err := service.Delete(ctx, id, images.V1DeleteParams{
			AccountID: cloudflare.F(accountID),
		})
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
		}
	}
	return nil
}
