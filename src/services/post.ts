import { Service } from 'typedi';
import { Repository } from 'typeorm';
import { InjectRepository } from 'typeorm-typedi-extensions';

import Post from '../models/post';
import { findPageItems, PagedItems, PageOptions } from '../lib/pagination';

@Service()
export default class PostService {
  constructor(
    @InjectRepository(Post) private readonly postRepository: Repository<Post>,
  ) {}

  public async findById(id: number, relations?: string[]): Promise<Post> {
    return this.postRepository.findOne({ where: { id }, relations: relations || [] });
  }

  public async findPageByBoardId(
    boardId: number, pageOpts: PageOptions, relations?: string[],
  ): Promise<PagedItems<Post>> {
    // eslint-disable-next-line @typescript-eslint/camelcase
    return findPageItems<Post>(this.postRepository, pageOpts, { board_id: boardId }, { id: 'DESC' }, relations);
  }
}
