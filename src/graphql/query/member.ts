import { Context } from 'koa';
import { Service } from 'typedi';

import { Permission, PermissionGuard } from '../guard';
import Member from '../../models/member';
import MemberService from '../../services/member';

@Service()
export default class MemberQuery {
  constructor(
    private readonly memberService: MemberService,
  ) {}

  @PermissionGuard()
  // eslint-disable-next-line class-methods-use-this
  async me(_: object, ctx: Partial<Context>): Promise<Member> {
    return ctx.state.member;
  }

  @PermissionGuard(Permission.Admin)
  async members(_: object, ctx: Partial<Context>): Promise<Member[]> {
    return this.memberService.findAll();
  }
}
