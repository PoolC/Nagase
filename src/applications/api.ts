import { Server } from 'http';
import Koa from 'koa';
import mount from 'koa-mount';
import graphql from 'koa-graphql';
import { Service } from 'typedi';

import Mutation from '../graphql/mutation';
import Query from '../graphql/query';
import schema from '../graphql/schema';
import Member from '../models/member';
import MemberService from '../services/member';

@Service()
export default class ApiApplication {
  private app: Koa;

  constructor(
    private readonly memberService: MemberService,

    private readonly query: Query,
    private readonly mutation: Mutation,
  ) {
    this.initialize();
  }

  initialize(): void {
    const app = new Koa();

    app.use(async (ctx, next) => {
      const { authorization } = ctx.req.headers;
      if (authorization) {
        let member: Member;
        try {
          const jwtToken = authorization.replace('Bearer ', '');
          member = await this.memberService.findByUuid(this.memberService.validateToken(jwtToken));
        } catch {
          ctx.throw(401);
          return;
        }

        if (member?.isActivated) {
          ctx.state.member = member;
        } else {
          ctx.throw(401);
        }
      }

      await next();
    });

    app.use(mount('/graphql', graphql({
      schema,
      rootValue: { ...this.query.all, ...this.mutation.all },
      graphiql: true,
    })));
    this.app = app;
  }

  start(): void {
    this.app.listen(process.env.PORT || 8080);
  }
}
