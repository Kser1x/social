package db

import (
	"context"
	"fmt"
	"github.com/Kser1x/social/internal/store"
	"log"
	"math/rand"
)

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creation user", err)
			return
		}
	}

	posts := generatePost(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creation user", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creation user", err)
			return
		}
	}

	return
}

func generateUsers(num int) []*store.UserModel {
	users := make([]*store.UserModel, num)

	for i := 0; i < num; i++ {
		users[i] = &store.UserModel{
			Username: fmt.Sprintf("user%v", i),
			Email:    fmt.Sprintf("user%v", i) + "@example.com",
			Password: "password",
		}
	}
	return users
}
func generatePost(num int, users []*store.UserModel) []*store.PostModel {
	posts := make([]*store.PostModel, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.PostModel{
			UserId:  user.ID,
			Title:   fmt.Sprintf("title%v", i),
			Content: fmt.Sprintf("content%v", i),
			Tags:    []string{"tag1", "tag2"},
		}
	}

	return posts
}
func generateComments(count int, users []*store.UserModel, posts []*store.PostModel) []*store.Comment {
	comments := make([]*store.Comment, count)
	usersLen := len(users)
	postsLen := len(posts)

	for i := 0; i < count; i++ {
		comments[i] = &store.Comment{
			Content: fmt.Sprintf(" some content%v", i),
			PostID:  posts[rand.Intn(postsLen)].ID,
			UserID:  users[rand.Intn(usersLen)].ID,
		}
	}

	return comments
}
