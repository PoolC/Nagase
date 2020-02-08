import { Service } from 'typedi';

import MemberMutation, {
  CreateAccessTokenInput, CreateMemberInput, MemberUUIDInput, UpdateMemberInput,
} from './member';
import BoardMutation, { BoardIDInput, CreateBoardInput, UpdateBoardInput } from './board';
import { QueryParams } from '../index';

@Service()
export default class Mutation {
  public all: {[_: string]: (...p: QueryParams) => any};

  constructor(
    private readonly memberMutation: MemberMutation,
    private readonly boardMutation: BoardMutation,
  ) {
    this.all = {
      createMember:
        (...params: QueryParams<CreateMemberInput>) => memberMutation.createMember(...params),
      updateMember:
        (...params: QueryParams<UpdateMemberInput>) => memberMutation.updateMember(...params),
      createAccessToken:
        (...params: QueryParams<CreateAccessTokenInput>) => memberMutation.createAccessToken(...params),
      deleteMember:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.deleteMember(...params),
      toggleMemberIsActivated:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.toggleMemberIsActivated(...params),
      toggleMemberIsAdmin:
        (...params: QueryParams<MemberUUIDInput>) => memberMutation.toggleMemberIsAdmin(...params),

      createBoard:
        (...params: QueryParams<CreateBoardInput>) => boardMutation.createBoard(...params),
      updateBoard:
        (...params: QueryParams<UpdateBoardInput>) => boardMutation.updateBoard(...params),
      deleteBoard:
        (...params: QueryParams<BoardIDInput>) => boardMutation.deleteBoard(...params),
    };
  }
}
