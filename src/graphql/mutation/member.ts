import { Service } from 'typedi';

import MemberService from '../../services/member';
import Member from '../../models/member';

export interface CreateMemberInput {
  MemberInput: {
    loginID: string;
    password: string;
    email: string;
    name: string;
    phoneNumber: string;
    department: string;
    studentID: string;
  };
}

export interface CreateAccessTokenInput {
  LoginInput: {
    loginID: string;
    password: string;
  };
}

@Service()
export default class MemberMutation {
  constructor(
    private readonly memberService: MemberService,
  ) {}

  public async createMember(input: CreateMemberInput): Promise<Member> {
    const member = new Member();
    member.loginId = input.MemberInput.loginID;
    member.email = input.MemberInput.email;
    member.phoneNumber = input.MemberInput.phoneNumber;
    member.name = input.MemberInput.name;
    member.department = input.MemberInput.department;
    member.studentId = input.MemberInput.studentID;
    await member.setPassword(input.MemberInput.password);

    const dupError = await this.memberService.checkDuplication(member);
    if (dupError) {
      throw new Error(dupError);
    }

    return this.memberService.save(member);
  }

  public async createAccessToken(input: CreateAccessTokenInput): Promise<{key: string}> {
    const member = await this.memberService.findByLoginId(input.LoginInput.loginID);
    if (!(await member?.validatePassword(input.LoginInput.password))) {
      throw new Error('TKN000');
    } else if (!member.isActivated) {
      throw new Error('TKN002');
    }

    return { key: this.memberService.generateToken(member) };
  }
}
