import { Context } from 'koa';
import { Service } from 'typedi';

import MemberService from '../../services/member';
import Member from '../../models/member';
import { Permission, PermissionGuard } from '../guard';

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

export interface UpdateMemberInput {
  MemberInput: {
    uuid: string;
    password?: string;
    email?: string;
    name?: string;
    phoneNumber?: string;
    department?: string;
    studentID?: string;
  };
}

export interface CreateAccessTokenInput {
  LoginInput: {
    loginID: string;
    password: string;
  };
}

export interface MemberUUIDInput {
  memberUUID: string;
}

@Service()
export default class MemberMutation {
  constructor(
    private readonly memberService: MemberService,
  ) {}

  public async createMember(args: CreateMemberInput, ctx?: Partial<Context>): Promise<Member> {
    const member = new Member();
    member.loginID = args.MemberInput.loginID;
    member.email = args.MemberInput.email;
    member.phoneNumber = args.MemberInput.phoneNumber;
    member.name = args.MemberInput.name;
    member.department = args.MemberInput.department;
    member.studentID = args.MemberInput.studentID;
    await member.setPassword(args.MemberInput.password);

    const dupError = await this.memberService.checkDuplication(member);
    if (dupError) {
      throw new Error(dupError);
    }

    return this.memberService.save(member);
  }

  public async updateMember(args: UpdateMemberInput, ctx: Partial<Context>): Promise<Member> {
    const member = await this.memberService.findByUuid(args.MemberInput.uuid);
    if (!member || ctx.state.member.uuid !== member.uuid) {
      return null;
    }

    if (args.MemberInput.password) {
      await member.setPassword(args.MemberInput.password);
    }
    member.name = args.MemberInput.name || member.name;
    member.email = args.MemberInput.email || member.email;
    member.phoneNumber = args.MemberInput.phoneNumber || member.phoneNumber;
    member.department = args.MemberInput.department || member.department;
    member.studentID = args.MemberInput.studentID || member.studentID;

    return this.memberService.save(member);
  }

  public async createAccessToken(args: CreateAccessTokenInput, ctx?: Partial<Context>): Promise<{key: string}> {
    const member = await this.memberService.findByLoginId(args.LoginInput.loginID);
    if (!(await member?.validatePassword(args.LoginInput.password))) {
      throw new Error('TKN000');
    } else if (!member.isActivated) {
      throw new Error('TKN002');
    }

    return { key: this.memberService.generateToken(member) };
  }

  @PermissionGuard(Permission.Admin)
  public async deleteMember(args: MemberUUIDInput, ctx?: Partial<Context>): Promise<Member> {
    const member = await this.memberService.findByUuid(args.memberUUID);
    if (member) {
      await this.memberService.delete(member);
    }
    return member;
  }

  @PermissionGuard(Permission.Admin)
  public async toggleMemberIsActivated(args: MemberUUIDInput, ctx?: Partial<Context>): Promise<Member> {
    const member = await this.memberService.findByUuid(args.memberUUID);
    if (member) {
      member.isActivated = !member.isActivated;
      await this.memberService.save(member);
    }
    return member;
  }

  @PermissionGuard(Permission.Admin)
  public async toggleMemberIsAdmin(args: MemberUUIDInput, ctx?: Partial<Context>): Promise<Member> {
    const member = await this.memberService.findByUuid(args.memberUUID);
    if (member) {
      member.isAdmin = !member.isAdmin;
      await this.memberService.save(member);
    }
    return member;
  }
}
