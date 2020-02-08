import { Service } from 'typedi';

import { Context } from 'koa';
import BoardService from '../../services/board';
import { Permission, PermissionGuard } from '../guard';
import Board from '../../models/board';

interface BoardInput {
  name: string;
  urlPath: string;
  readPermission: string;
  writePermission: string;
}

export interface CreateBoardInput {
  BoardInput: BoardInput;
}

export interface UpdateBoardInput {
  boardID: number;
  BoardInput: BoardInput;
}

export interface BoardIDInput {
  boardID: number;
}

@Service()
export default class BoardMutation {
  constructor(
    private readonly boardService: BoardService,
  ) {}

  @PermissionGuard(Permission.Admin)
  public async createBoard(args: CreateBoardInput, ctx: Partial<Context>): Promise<Board> {
    const board = new Board();
    board.name = args.BoardInput.name;
    board.urlPath = args.BoardInput.urlPath;
    board.readPermission = Board.permissionFromString(args.BoardInput.readPermission);
    board.writePermission = Board.permissionFromString(args.BoardInput.writePermission);

    return this.boardService.save(board);
  }

  @PermissionGuard(Permission.Admin)
  public async updateBoard(args: UpdateBoardInput, ctx: Partial<Context>): Promise<Board> {
    const board = await this.boardService.findById(args.boardID);
    if (!board) {
      return null;
    }

    board.name = args.BoardInput.name || board.name;
    board.urlPath = args.BoardInput.urlPath || board.urlPath;
    if (args.BoardInput.readPermission) {
      board.readPermission = Board.permissionFromString(args.BoardInput.readPermission);
    }
    if (args.BoardInput.writePermission) {
      board.writePermission = Board.permissionFromString(args.BoardInput.writePermission);
    }

    return this.boardService.save(board);
  }

  @PermissionGuard(Permission.Admin)
  public async deleteBoard(args: BoardIDInput, ctx: Partial<Context>): Promise<Board> {
    const board = await this.boardService.findById(args.boardID);
    if (board) {
      await this.boardService.delete(board);
    }
    return board;
  }
}
