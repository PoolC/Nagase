import { Service } from 'typedi';

import MemberMutation, { CreateAccessTokenInput, CreateMemberInput } from './member';

@Service()
export default class Mutation {
  public all: {[_: string]: (input: any) => any};

  constructor(
    private readonly memberMutation: MemberMutation,
  ) {
    this.all = {
      createMember: (input: CreateMemberInput) => memberMutation.createMember(input),
      createAccessToken: (input: CreateAccessTokenInput) => memberMutation.createAccessToken(input),
    };
  }
}
