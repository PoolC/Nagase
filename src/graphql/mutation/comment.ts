import { Context } from 'koa';
import { Service } from 'typedi';

import { PermissionGuard } from '../guard';
import Comment from '../../models/comment';
import CommentService from '../../services/comment';
import PostService from '../../services/post';

export interface CreateCommentInput {
  postID: number;
  body: string;
}

export interface DeleteCommentInput {
  commentID: number;
}

@Service()
export default class CommentMutation {
  constructor(
    private readonly commentService: CommentService,
    private readonly postService: PostService,
  ) {}

  @PermissionGuard()
  public async createComment(args: CreateCommentInput, ctx: Partial<Context>) {
    const post = await this.postService.findById(args.postID, ['board']);
    if (!post || !post.board.readPermittedTo(ctx.state.member)) {
      throw new Error('ERR400');
    }

    const comment = new Comment();
    comment.author = ctx.state.member;
    comment.post = post;
    comment.body = args.body;

    return this.commentService.save(comment);
  }

  @PermissionGuard()
  public async deleteComment(args: DeleteCommentInput, ctx: Partial<Context>) {
    const comment = await this.commentService.findById(args.commentID, ['author']);
    if (!comment || (!ctx.state.member.isAdmin && comment.author.uuid !== ctx.state.member.uuid)) {
      throw new Error('ERR400');
    }

    await this.commentService.delete(comment);
    return comment;
  }
}
