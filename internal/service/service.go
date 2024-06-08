package service

type ServiceManager struct {
	userService      *UserService
	communityService *CommunityService
	postService      *PostService
	commentService   *CommentService
}

func NewServiceManager(userService *UserService, communityService *CommunityService, postService *PostService, commentService *CommentService) *ServiceManager {
	return &ServiceManager{
		userService:      userService,
		communityService: communityService,
		postService:      postService,
		commentService:   commentService,
	}
}

func (sm *ServiceManager) GetCommunityService() *CommunityService {
	return sm.communityService
}

func (sm *ServiceManager) GetPostService() *PostService {
	return sm.postService
}

func (sm *ServiceManager) GetUserService() *UserService {
	return sm.userService
}

func (sm *ServiceManager) GetCommentService() *CommentService {
	return sm.commentService
}
