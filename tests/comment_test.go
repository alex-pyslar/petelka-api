package tests

import (
	"context"
	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/alex-pyslar/petelka-api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	commentRepo := repository.NewCommentRepository(db, redisClient)
	commentService := service.NewCommentService(commentRepo)

	comment := &models.Comment{
		ProductID: 1,
		UserID:    1,
		Text:      "Test Comment",
	}

	err := commentService.CreateComment(context.Background(), comment)
	assert.NoError(t, err)
	assert.NotZero(t, comment.ID)
}

func TestGetComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	commentRepo := repository.NewCommentRepository(db, redisClient)
	commentService := service.NewCommentService(commentRepo)

	comment := &models.Comment{
		ProductID: 1,
		UserID:    1,
		Text:      "Test Comment 2",
	}
	commentService.CreateComment(context.Background(), comment)

	fetchedComment, err := commentService.GetComment(context.Background(), comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.Text, fetchedComment.Text)
}

func TestListComments(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	commentRepo := repository.NewCommentRepository(db, redisClient)
	commentService := service.NewCommentService(commentRepo)

	comment1 := &models.Comment{ProductID: 1, UserID: 1, Text: "Comment 1"}
	comment2 := &models.Comment{ProductID: 1, UserID: 1, Text: "Comment 2"}
	commentService.CreateComment(context.Background(), comment1)
	commentService.CreateComment(context.Background(), comment2)

	comments, err := commentService.ListComments(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(comments), 2)
}

func TestUpdateComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	commentRepo := repository.NewCommentRepository(db, redisClient)
	commentService := service.NewCommentService(commentRepo)

	comment := &models.Comment{ProductID: 1, UserID: 1, Text: "Old Comment"}
	commentService.CreateComment(context.Background(), comment)

	comment.Text = "Updated Comment"
	err := commentService.UpdateComment(context.Background(), comment)
	assert.NoError(t, err)

	fetchedComment, err := commentService.GetComment(context.Background(), comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Comment", fetchedComment.Text)
}

func TestDeleteComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	commentRepo := repository.NewCommentRepository(db, redisClient)
	commentService := service.NewCommentService(commentRepo)

	comment := &models.Comment{ProductID: 1, UserID: 1, Text: "Delete Comment"}
	commentService.CreateComment(context.Background(), comment)

	err := commentService.DeleteComment(context.Background(), comment.ID)
	assert.NoError(t, err)

	_, err = commentService.GetComment(context.Background(), comment.ID)
	assert.Error(t, err)
}
