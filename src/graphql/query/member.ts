import { Context } from 'koa';

import Member from '../../models/member';

export default class MemberQuery {
  // eslint-disable-next-line class-methods-use-this
  async me(_: object, ctx: Partial<Context>): Promise<Member> {
    const member = ctx.state?.member;
    if (!member) {
      throw new Error('ERR401');
    }
    return member;
  }
}
