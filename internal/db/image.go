package db

import (
	"net/url"
	"strings"
	"time"

	"github.com/arwoosa/vulpes/db/mgo"
	"github.com/arwoosa/vulpes/db/mgo/types"
	"github.com/arwoosa/vulpes/validate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {
	mgo.RegisterIndex(imageCollection)
}

const ImageCollectionName = "images"

var (
	imageCollection = mgo.NewCollectDef(ImageCollectionName, func() []mongo.IndexModel {
		optionsBuilder := &options.IndexOptionsBuilder{}
		optionsBuilder.SetUnique(true)
		return []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "cloudflare_id", Value: 1}},
				Options: optionsBuilder,
			},
			{
				Keys: bson.D{{Key: "location", Value: "2dsphere"}},
			},
		}
	})
)

type imageOption func(*image)

func WithImageCloudflareID(id string) imageOption {
	return func(i *image) {
		i.CloudflareID = id
	}
}

func WithImageFilename(filename string) imageOption {
	return func(i *image) {
		i.Filename = filename
	}
}

func WithImageUploaded(uploaded time.Time) imageOption {
	return func(i *image) {
		i.Uploaded = uploaded
	}
}

func WithImageMeta(meta map[string]string) imageOption {
	return func(i *image) {
		i.Meta = meta
	}
}

func WithImageVariants(variants []string) imageOption {
	return func(i *image) {
		result := map[string]string{}
		for _, v := range variants {
			parsed, err := url.Parse(v)
			if err != nil {
				continue // Skip invalid URLs
			}

			p := parsed.Path
			if p == "" {
				continue // Skip URLs with no path
			}

			// The path from Cloudflare is expected to be in the format:
			// /<account_hash>/<image_id>/<variant_name>
			// We want to transform this into a relative path for our CDN:
			// /cdn-images/<image_id>/<variant_name>

			// Find the start of the image_id part of the path, which is after the first path segment (account_hash).
			// We find the first '/' after the initial one.
			firstSlashIndex := strings.Index(p[1:], "/")
			if firstSlashIndex == -1 {
				// This means the path is something like "/segment", not "/segment1/segment2/..."
				// which doesn't match the expected format.
				continue
			}

			// The index is relative to p[1:], so add 1 to get the absolute index in p.
			pathSuffix := p[firstSlashIndex+1:]
			value := "/cdn-images" + pathSuffix

			// The key for our map is the variant name, which is the last part of the path.
			// We trim slashes to correctly get the last segment.
			cleanPath := strings.Trim(p, "/")
			pathSegments := strings.Split(cleanPath, "/")
			if len(pathSegments) == 0 {
				continue
			}
			key := pathSegments[len(pathSegments)-1]

			result[key] = value
		}
		i.Variants = result
	}
}

func WithImageCount(count int) imageOption {
	return func(i *image) {
		i.Count = count
	}
}

func WithLocation(longitude, latitude *float64) imageOption {
	return func(i *image) {
		if longitude == nil || latitude == nil {
			return
		}
		i.Location = types.NewLocationPoint(*longitude, *latitude)
	}
}

func WithSize(size uint64) imageOption {
	return func(i *image) {
		i.Size = size
	}
}

type image struct {
	mgo.Index    `bson:"-"`
	ID           bson.ObjectID   `bson:"_id,omitempty" validate:"required"`
	CloudflareID string          `bson:"cloudflare_id,omitempty" validate:"required"`
	Filename     string          `bson:"filename,omitempty" validate:"required"`
	Uploaded     time.Time       `bson:"uploaded,omitempty" validate:"required"`
	Size         uint64          `bson:"size,omitempty" validate:"required"`
	Location     *types.Location `bson:"location,omitempty"`

	Meta     map[string]string `bson:"meta,omitempty"`
	Variants map[string]string `bson:"variants,omitempty" validate:"required"`
	Count    int               `bson:"count,omitempty"`
}

func (i *image) Validate() error {
	return validate.Struct(i)
}

func (i *image) GetId() any {
	return i.ID
}

func (i *image) SetId(id any) {
	if oid, ok := id.(bson.ObjectID); ok {
		i.ID = oid
		return
	}
}

func NewImage(opts ...imageOption) *image {
	i := &image{
		Index: imageCollection,
		ID:    bson.NewObjectID(),
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}
