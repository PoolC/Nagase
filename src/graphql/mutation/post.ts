import { Context } from 'koa';
import { Service } from 'typedi';

import Post from '../../models/post';
import BoardService from '../../services/board';
import PostService from '../../services/post';
import { PermissionGuard } from '../guard';

export interface CreatePostInput {
  boardID: number;
  PostInput: {
    title: string;
    body: string;
  };
}

export interface UpdatePostInput {
  postID: number;
  PostInput: {
    title?: string;
    body?: string;
  };
}

export interface DeletePostInput {
  postID: number;
}

@Service()
export default class PostMutation {
  constructor(
    private readonly postService: PostService,
    private readonly boardService: BoardService,
  ) {}

  @PermissionGuard()
  public async createPost(args: CreatePostInput, ctx: Partial<Context>): Promise<Post> {
    const board = await this.boardService.findById(args.boardID);
    if (!board || !board.writePermittedTo(ctx.state.member)) {
      throw new Error('ERR400');
    }

    const post = new Post();
    post.board = board;
    post.author = ctx.state.member;
    post.title = args.PostInput.title;
    post.body = args.PostInput.body;

    // TODO: create the vote
    // TODO: send push messages to subscribers.

    return this.postService.save(post);
  }

  @PermissionGuard()
  public async updatePost(args: UpdatePostInput, ctx: Partial<Context>): Promise<Post> {
    const post = await this.postService.findById(args.postID, ['author']);
    if (!post || post.author.uuid !== ctx.state.member.uuid) {
      throw new Error('ERR400');
    }

    post.title = args.PostInput.title || post.title;
    post.body = args.PostInput.body || post.body;

    return this.postService.save(post);
  }

  @PermissionGuard()
  public async deletePost(args: DeletePostInput, ctx: Partial<Context>): Promise<Post> {
    const post = await this.postService.findById(args.postID, ['author']);
    if (!post || (!ctx.state.member.isAdmin && (post.author.uuid !== ctx.state.member.uuid))) {
      throw new Error('ERR400');
    }

    await this.postService.delete(post);
    return post;
  }
}
