import { Context } from 'koa';
import { Service } from 'typedi';

import Post from '../../models/post';
import PostService from '../../services/post';
import { PagedItems } from '../../lib/pagination';
import BoardService from '../../services/board';

export interface PostIDInput {
  postID: number;
}

export interface PostPageInput {
  boardID: number;
  page?: number;
  count?: number;
}

@Service()
export default class PostQuery {
  constructor(
    private readonly postService: PostService,
    private readonly boardService: BoardService,
  ) {}

  public async post(args: PostIDInput, ctx: Partial<Context>): Promise<Post> {
    const post = await this.postService.findById(args.postID, ['board', 'author']);
    if (!post) {
      throw new Error('ERR400');
    }

    if (!post.board.readPermittedTo(ctx.state?.member)) {
      throw new Error('ERR403');
    }
    return post;
  }

  public async postPage(args: PostPageInput, ctx: Partial<Context>): Promise<PagedItems<Post>> {
    const board = await this.boardService.findById(args.boardID);
    if (!board.readPermittedTo(ctx.state?.member)) {
      throw new Error('ERR403');
    }

    const pageOpts = { page: args.page || 1, count: args.count || 20 };
    return this.postService.findPageByBoardId(args.boardID, pageOpts, ['board', 'author']);
  }
}
