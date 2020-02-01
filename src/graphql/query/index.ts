import { Service } from 'typedi';
import { Context } from 'koa';

import MemberQuery from './member';

@Service()
export default class Query {
  public all: {[_: string]: (input: any, ctx: Context) => any};

  constructor(
    private readonly memberQuery: MemberQuery,
  ) {
    this.all = {
      me: (input: object, ctx: Context) => memberQuery.me(input, ctx),
    };
  }
}
