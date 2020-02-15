import { Service } from 'typedi';

import MemberMutation, {
  CreateAccessTokenInput, CreateMemberInput, MemberUUIDInput, UpdateMemberInput,
} from './member';
import BoardMutation, { BoardIDInput, CreateBoardInput, UpdateBoardInput } from './board';
import CommentMutation, { CreateCommentInput, DeleteCommentInput } from './comment';
import PostMutation, { CreatePostInput, DeletePostInput, UpdatePostInput } from './post';
import ProjectMutation, { CreateProjectInput, DeleteProjectInput, UpdateProjectInput } from './project';
import { QueryParams } from '../index';

@Service()
export default class Mutation {
  public all: {[_: string]: (...p: QueryParams) => any};

  constructor(
    private readonly memberMutation: MemberMutation,
    private readonly boardMutation: BoardMutation,
    private readonly postMutation: PostMutation,
    private readonly commentMutation: CommentMutation,
    private readonly projectMutation: ProjectMutation,
  ) {
    this.all = {
      createMember:
        (...params: QueryParams<CreateMemberInput>) => memberMutation.createMember(...params),
      updateMember:
        (...params: QueryParams<UpdateMemberInput>) => memberMutation.updateMember(...params),
      createAccessToken:
        (...params: QueryParams<CreateAccessTokenInput>) => memberMutation.createAccessToken(...params),
      deleteMember:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.deleteMember(...params),
      toggleMemberIsActivated:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.toggleMemberIsActivated(...params),
      toggleMemberIsAdmin:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.toggleMemberIsAdmin(...params),

      createBoard:
        (...params: QueryParams<CreateBoardInput>) => boardMutation.createBoard(...params),
      updateBoard:
        (...params: QueryParams<UpdateBoardInput>) => boardMutation.updateBoard(...params),
      deleteBoard:
        (...params: QueryParams<BoardIDInput>) => boardMutation.deleteBoard(...params),

      createPost:
        (...params: QueryParams<CreatePostInput>) => postMutation.createPost(...params),
      updatePost:
        (...params: QueryParams<UpdatePostInput>) => postMutation.updatePost(...params),
      deletePost:
        (...params: QueryParams<DeletePostInput>) => postMutation.deletePost(...params),

      createComment:
        (...params: QueryParams<CreateCommentInput>) => commentMutation.createComment(...params),
      deleteComment:
        (...params: QueryParams<DeleteCommentInput>) => commentMutation.deleteComment(...params),

      createProject:
        (...params: QueryParams<CreateProjectInput>) => projectMutation.createProject(...params),
      updateProject:
        (...params: QueryParams<UpdateProjectInput>) => projectMutation.updateProject(...params),
      deleteProject:
        (...params: QueryParams<DeleteProjectInput>) => projectMutation.deleteProject(...params),
    };
  }
}
