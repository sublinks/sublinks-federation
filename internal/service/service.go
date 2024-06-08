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

func (sm *ServiceManager) CommunityService() *CommunityService {
	return sm.communityService
}

func (sm *ServiceManager) PostService() *PostService {
	return sm.postService
}

func (sm *ServiceManager) UserService() *UserService {
	return sm.userService
}

func (sm *ServiceManager) CommentService() *CommentService {
	return sm.commentService
}
