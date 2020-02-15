import { Service } from 'typedi';
import { Repository } from 'typeorm';
import { InjectRepository } from 'typeorm-typedi-extensions';

import Comment from '../models/comment';

@Service()
export default class CommentService {
  constructor(
    @InjectRepository(Comment) private readonly commentRepository: Repository<Comment>,
  ) {}

  public async findById(commentId: number, relations: string[]): Promise<Comment> {
    return this.commentRepository.findOne({ where: { id: commentId }, relations: relations || [] });
  }

  public async save(obj: Comment): Promise<Comment> {
    return this.commentRepository.save(obj);
  }

  public async delete(obj: Comment): Promise<void> {
    await this.commentRepository.delete(obj);
  }
}
