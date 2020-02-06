import { Service } from 'typedi';
import { Context } from 'koa';

import BoardService from '../../services/board';
import Board from '../../models/board';

export interface BoardIDInput {
  boardID: number;
}

@Service()
export default class BoardQuery {
  constructor(
    private readonly boardService: BoardService,
  ) {}

  async board(args: BoardIDInput, ctx?: Context): Promise<Board> {
    return this.boardService.findById(args.boardID);
  }

  async boards(_?: object, ctx?: Context): Promise<Board[]> {
    return this.boardService.findAll();
  }
}
