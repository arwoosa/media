package db

import (
	"context"
	"fmt"

	"github.com/arwoosa/vulpes/relation"
)

const (
	nsImage = "Image"
)

func SaveImageUserOwner(ctx context.Context, userId string, imageIds []string) error {
	if len(imageIds) == 0 {
		return nil
	}
	for _, id := range imageIds {
		err := relation.AddUserResourceRole(ctx, userId, nsImage, id, relation.RoleOwner)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRelation, err)
		}
	}
	return nil
}

func DeleteImageUserRelation(ctx context.Context, imageIds ...string) error {
	for _, id := range imageIds {
		err := relation.DeleteObjectId(ctx, nsImage, id)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRelation, err)
		}
	}
	return nil
}
