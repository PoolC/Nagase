import { Service } from 'typedi';

import MemberQuery from './member';
import { QueryParams } from '..';
import BoardQuery, { BoardIDInput } from './board';

@Service()
export default class Query {
  public all: {[_: string]: (...p: QueryParams) => any};

  constructor(
    private readonly memberQuery: MemberQuery,
    private readonly boardQuery: BoardQuery,
  ) {
    this.all = {
      me: (...params: QueryParams) => memberQuery.me(...params),
      members: (...params: QueryParams) => memberQuery.members(...params),
      board: (...params: QueryParams<BoardIDInput>) => boardQuery.board(...params),
      boards: (...params: QueryParams) => boardQuery.boards(...params),
    };
  }
}
