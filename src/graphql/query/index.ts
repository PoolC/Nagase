import { Service } from 'typedi';

import MemberQuery from './member';
import { QueryParams } from '..';

@Service()
export default class Query {
  public all: {[_: string]: (...p: QueryParams) => any};

  constructor(
    private readonly memberQuery: MemberQuery,
  ) {
    this.all = {
      me: (...params: QueryParams) => memberQuery.me(...params),
      members: (...params: QueryParams) => memberQuery.members(...params),
    };
  }
}
