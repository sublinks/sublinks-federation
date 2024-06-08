package service

import (
	"sublinks/sublinks-federation/internal/service/actors"
	"sublinks/sublinks-federation/internal/service/comments"
	"sublinks/sublinks-federation/internal/service/posts"
)

type ServiceManager struct {
	userService      *actors.UserService
	communityService *actors.CommunityService
	postService      *posts.PostService
	commentService   *comments.CommentService
}

func NewServiceManager(userService *actors.UserService, communityService *actors.CommunityService, postService *posts.PostService, commentService *comments.CommentService) *ServiceManager {
	return &ServiceManager{
		userService:      userService,
		communityService: communityService,
		postService:      postService,
		commentService:   commentService,
	}
}

func (sm *ServiceManager) GetCommunityService() *actors.CommunityService {
	return sm.communityService
}

func (sm *ServiceManager) GetPostService() *posts.PostService {
	return sm.postService
}

func (sm *ServiceManager) GetUserService() *actors.UserService {
	return sm.userService
}

func (sm *ServiceManager) GetCommentService() *comments.CommentService {
	return sm.commentService
}
