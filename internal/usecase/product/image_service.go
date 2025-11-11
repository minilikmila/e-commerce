package product

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/minilik/ecommerce/internal/domain"
	"github.com/minilik/ecommerce/internal/domain/repository"
	"github.com/minilik/ecommerce/pkg/cloudinary"
)

type ImageService interface {
	UploadImages(ctx context.Context, productID uuid.UUID, files []*multipart.FileHeader) ([]domain.ProductImage, error)
	ListImages(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error)
}

type imageService struct {
	imagesRepo repository.ProductImageRepository
	uploader   *cloudinary.Client
	logger     *zap.Logger
	now        func() time.Time
}

func NewImageService(repo repository.ProductImageRepository, uploader *cloudinary.Client, logger *zap.Logger) ImageService {
	return &imageService{
		imagesRepo: repo,
		uploader:   uploader,
		logger:     logger,
		now:        time.Now,
	}
}

func (s *imageService) UploadImages(ctx context.Context, productID uuid.UUID, files []*multipart.FileHeader) ([]domain.ProductImage, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}
	if len(files) > 4 {
		return nil, fmt.Errorf("maximum 4 images allowed per request")
	}
	current, err := s.imagesRepo.CountByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	if current+int64(len(files)) > 4 {
		return nil, fmt.Errorf("upload would exceed limit of 4 images per product")
	}

	var uploaded []domain.ProductImage
	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			return nil, err
		}
		func() { // anonymous func here:
			defer src.Close()
			filename := safeFilename(fh.Filename)
			var url string
			var err error
			// Prefer signed upload when API key/secret are configured but unsigned / unauthenticated for worst case
			if s.uploader != nil && s.uploader.APIKey != "" && s.uploader.APISecret != "" {
				url, err = s.uploader.UploadSigned(ctx, src, filename, nil)
			} else {
				url, err = s.uploader.UploadUnsigned(ctx, src, filename)
			}
			if err != nil {
				s.logger.Error("cloudinary upload failed", zap.Error(err))
				return
			}
			uploaded = append(uploaded, domain.ProductImage{
				ID:        uuid.New(),
				ProductID: productID,
				URL:       url,
				CreatedAt: s.now(),
			})
		}()
	}
	if len(uploaded) == 0 {
		return nil, fmt.Errorf("no image uploaded")
	}
	if err := s.imagesRepo.AddMany(ctx, uploaded); err != nil {
		return nil, err
	}
	return uploaded, nil
}

func (s *imageService) ListImages(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error) {
	return s.imagesRepo.ListByProduct(ctx, productID)
}

func safeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
