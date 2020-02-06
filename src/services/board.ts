import { Service } from 'typedi';
import { InjectRepository } from 'typeorm-typedi-extensions';
import { Repository } from 'typeorm';

import Board from '../models/board';

@Service()
export default class BoardService {
  constructor(
    @InjectRepository(Board) private readonly boardRepository: Repository<Board>,
  ) {}

  public async findById(id: number): Promise<Board> {
    return this.boardRepository.findOne({ id });
  }

  public async findAll(): Promise<Board[]> {
    return this.boardRepository.find({ order: { id: 'ASC' } });
  }
}
