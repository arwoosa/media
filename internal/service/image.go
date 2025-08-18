package service

import (
	"context"
	"errors"
	"time"

	"github.com/arwoosa/media/internal/cloudflare"
	"github.com/arwoosa/media/internal/db"
	"github.com/arwoosa/media/internal/pb/image"
	"github.com/arwoosa/vulpes/db/cache"
	"github.com/arwoosa/vulpes/db/mgo"
	"github.com/arwoosa/vulpes/ezgrpc"

	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// imageServer 實作了 image.UnimplementedMediaServiceServer gRPC 服務。
type imageServer struct {
	image.UnimplementedImageServiceServer
}

var ()

func init() {
	// 將 imageServer 注入到 ezgrpc 中，以便 gRPC 伺服器可以註冊它。
	ezgrpc.InjectGrpcService(func(s grpc.ServiceRegistrar) {
		image.RegisterImageServiceServer(s, &imageServer{})
	})
	// 註冊 gRPC-Gateway 處理程序，將 HTTP 請求代理到 gRPC 服務。
	ezgrpc.RegisterHandlerFromEndpoint(image.RegisterImageServiceHandlerFromEndpoint)
}

// signedUrlSlice 是 []*image.SignedUrl 的輔助類型，用於簡化操作。
type signedUrlSlice []*image.SignedUrl

// GetImageIds 從 signedUrlSlice 中提取所有圖片的 ID。
func (s signedUrlSlice) GetImageIds() []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = v.ImageId
	}
	return result
}

// BatchUpload 處理批次圖片上傳請求。
// 它會為請求中的每張圖片生成一個預簽名的上傳 URL。
// 為了實現冪等性，它會將生成的 URL 存儲在會話中。
// 如果在同一個會話中再次調用，它將返回先前生成的 URL，而不是創建新的。
func (s *imageServer) BatchUpload(ctx context.Context, req *image.UploadRequest) (*image.UploadResponse, error) {
	// 1. 檢查會話中是否存在數據，如果存在，則直接返回數據以實現冪等性。
	data, err := ezgrpc.GetSessionData[signedUrlSlice](ctx)
	if err != nil && !errors.Is(err, ezgrpc.ErrSessionNotFound) {
		return nil, ezgrpc.ToStatus(err).Err()
	}
	if data != nil {
		return &image.UploadResponse{
			Images: data,
		}, nil
	}

	// 2. 為每張圖片生成預簽名的 URL。
	uploadImages := make(signedUrlSlice, len(req.Images))
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	for i := range req.Images {
		// 2.1. 獲取預簽名的 URL
		signedUrl, err := cloudflare.GetSignedUrl(ctx,
			cloudflare.ImageMetadataSize(req.Images[i].Size),
			cloudflare.ImageMetadataWidth(req.Images[i].Width),
			cloudflare.ImageMetadataHeight(req.Images[i].Height),
			cloudflare.ImageMetadataFormat(req.Images[i].ContentType.String()),
			cloudflare.ImageMetadataLatitude(req.Images[i].Latitude),
			cloudflare.ImageMetadataLongitude(req.Images[i].Longitude))
		if err != nil {
			return nil, cloudflare.ToStatus(err).Err()
		}
		uploadImages[i] = &image.SignedUrl{
			ImageId:   signedUrl.ID,
			SignedUrl: signedUrl.UploadURL,
		}
	}

	// 3. 將生成的 URL 數據設置到會話中。
	err = ezgrpc.SetSessionData(ctx, uploadImages)
	if err != nil {
		return nil, ezgrpc.ToStatus(err).Err()
	}

	// 4. 在響應中返回 URL。
	return &image.UploadResponse{
		Images: uploadImages,
	}, nil
}

// Complete 檢查圖片的上傳狀態。
// 它會從會話中檢索圖片 ID，然後向 Cloudflare 查詢這些圖片的詳細信息。
// 成功獲取信息後，它會刪除會話。
func (s *imageServer) Complete(ctx context.Context, req *image.StatusRequest) (*image.StatusResponse, error) {
	// 1. 從會話中獲取圖片數據。
	data, err := ezgrpc.GetSessionData[signedUrlSlice](ctx)
	if err != nil {
		return nil, ezgrpc.ToStatus(err).Err()
	}
	imageIds := data.GetImageIds()

	// 2. 查詢 Cloudflare 以獲取圖片的詳細信息。
	completeCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	images := cloudflare.GetImages(completeCtx, imageIds)
	result := make([]*image.ImageStatus, len(imageIds))
	bulk, err := mgo.NewBulkOperation(db.NewImage().C())
	if err != nil {
		return nil, mgo.ToStatus(err).Err()
	}

	for i, id := range imageIds {
		saveImage := images[id].Image
		myImage := db.NewImage(
			db.WithImageCloudflareID(id),
			db.WithImageFilename(saveImage.Filename),
			db.WithImageUploaded(saveImage.Uploaded),
			db.WithImageMeta(saveImage.Meta),
			db.WithImageVariants(saveImage.Variants),
			db.WithImageCount(0),
			db.WithSize(saveImage.GetSize()),
			db.WithLocation(saveImage.GetLongitude(), saveImage.GetLatitude()),
		)
		bulk.InsertOne(myImage)
		result[i] = &image.ImageStatus{
			ImageId: id,
			Metadata: &image.ImageMetadata{
				Width:      saveImage.GetWidth(),
				Height:     saveImage.GetHeight(),
				Format:     saveImage.GetImageFormat(),
				Size:       saveImage.GetSize(),
				UploadTime: saveImage.Uploaded.Format(time.RFC3339),
			},
			Variants: myImage.Variants,
		}
	}

	// 3. 存入資料庫

	_, err = bulk.Execute(ctx)
	if err != nil {
		return nil, mgo.ToStatus(err).Err()
	}

	// 4. 建立關係
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, ezgrpc.ToStatus(err).Err()
	}
	if user != nil {
		err = db.SaveImageUserOwner(ctx, user.ID, imageIds)
		if err != nil {
			return nil, db.ToStatus(err).Err()
		}
	}

	// 5. 刪除會話。
	err = ezgrpc.DeleteSession(ctx)
	if err != nil {
		return nil, ezgrpc.ToStatus(err).Err()
	}

	// 6. 返回包含圖片狀態和元數據的響應。
	return &image.StatusResponse{
		Images: result,
	}, nil
}

// Clear 清除預簽名 URL 的緩存。
func (s *imageServer) Clear(ctx context.Context, req *image.ClearRequest) (*image.ClearResponse, error) {
	// 1. 清除給定命名空間的預簽名 URL 緩存。
	err := ezgrpc.DeleteSession(ctx)
	if err != nil {
		return nil, ezgrpc.ToStatus(err).Err()
	}
	// 2. 返回成功響應。
	return &image.ClearResponse{
		Message: "Cache cleared successfully",
	}, nil
}

// Delete 刪除單張圖片。
func (s *imageServer) Delete(ctx context.Context, req *image.DeleteRequest) (*image.DeleteResponse, error) {
	// 1. 刪除圖片
	err := cloudflare.DeleteImages(ctx, req.GetImageId())
	if err != nil {
		return nil, cloudflare.ToStatus(err).Err()
	}
	// 2. 刪除資料庫中的圖片
	_, err = mgo.DeleteMany(ctx, db.NewImage(), bson.D{{Key: "cloudflare_id", Value: req.GetImageId()}})
	if err != nil {
		return nil, mgo.ToStatus(err).Err()
	}
	// 3. 刪除資料庫中的圖片關係
	err = db.DeleteImageUserRelation(ctx, req.ImageId)
	if err != nil {
		return nil, db.ToStatus(err).Err()
	}
	// 4. 返回成功響應。
	return &image.DeleteResponse{
		Message: "Image deleted successfully",
	}, nil
}

// BatchDelete 刪除多張圖片。
func (s *imageServer) BatchDelete(ctx context.Context, req *image.BatchDeleteRequest) (*image.BatchDeleteResponse, error) {
	// 1. 刪除圖片
	err := cloudflare.DeleteImages(ctx, req.GetImageIds()...)
	if err != nil {
		return nil, cloudflare.ToStatus(err).Err()
	}
	// 2. 刪除資料庫中的圖片
	_, err = mgo.DeleteMany(ctx, db.NewImage(), bson.D{{Key: "cloudflare_id", Value: bson.M{"$in": req.GetImageIds()}}})
	if err != nil {
		return nil, mgo.ToStatus(err).Err()
	}
	// 3. 刪除資料庫中的圖片關係
	err = db.DeleteImageUserRelation(ctx, req.ImageIds...)
	if err != nil {
		return nil, db.ToStatus(err).Err()
	}
	// 4. 返回成功響應。
	return &image.BatchDeleteResponse{
		Message: "Images deleted successfully",
	}, nil
}

func (s *imageServer) GetImageURI(ctx context.Context, req *image.ImageRequest) (*image.ImageResponse, error) {
	// 1. 從資料庫中取得圖片的 Cloudflare ID
	queryImg := db.NewImage()
	queryCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	err := mgo.FindOne(queryCtx, queryImg, bson.M{"cloudflare_id": req.GetId()})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, "Image not found")
		}
		return nil, mgo.ToStatus(err).Err()
	}
	// 2. 取得圖片的變體
	url, ok := queryImg.Variants[req.GetVariant()]
	if !ok {
		return nil, status.Error(codes.NotFound, "Variant not found")
	}
	// 3. 設置重定向 URL
	ezgrpc.SetRedirectUrl(ctx, url)
	// 4. 增加計數器
	_, err = cache.Incr(ctx, queryImg.CloudflareID)
	if err != nil {
		return nil, cache.ToStatus(err).Err()
	}
	// 5. 返回成功響應。
	return &image.ImageResponse{
		Uri: url,
	}, nil
}

func (s *imageServer) SyncImageCount(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	// random error
	time.Sleep(time.Second * 2)

	// if rand.New(rand.NewSource(time.Now().UnixNano())).Int()%2 == 0 {
	// 	return nil, status.Error(codes.Internal, "Internal error")
	// }
	return &empty.Empty{}, nil
	// 1. get all key from cache
	err := cache.DeleteAfterScanExecuteInt(ctx, "*", func(key string, val int) error {
		updateCtx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		fields := bson.D{{Key: "count", Value: val}}
		_, err := mgo.UpdateOne(
			updateCtx, db.NewImage(),
			bson.D{{Key: "cloudflare_id", Value: key}},
			bson.D{{Key: "$inc", Value: fields}},
		)
		return err
	})

	if err != nil {
		switch {
		case errors.Is(err, mgo.ErrWriteFailed):
			return nil, mgo.ToStatus(err).Err()
		default:
			return nil, cache.ToStatus(err).Err()
		}
	}
	return &empty.Empty{}, nil
}
