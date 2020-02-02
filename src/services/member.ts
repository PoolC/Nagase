import jwt from 'jsonwebtoken';
import { Service } from 'typedi';
import { Repository } from 'typeorm';

import { InjectRepository } from 'typeorm-typedi-extensions';
import Member from '../models/member';

type MemberClaim = { member_uuid: string };

@Service()
export default class MemberService {
  public hmacSecret = process.env.NAGASE_SECRET_KEY;

  constructor(
    @InjectRepository(Member) private readonly memberRepository: Repository<Member>,
  ) {}

  public async findByUuid(uuid: string): Promise<Member> {
    return this.memberRepository.findOne({ uuid });
  }

  public async findByLoginId(loginId: string): Promise<Member> {
    return this.memberRepository.findOne({ loginID: loginId });
  }

  public async findAll(): Promise<Member[]> {
    return this.memberRepository.find();
  }

  public async checkDuplication(member: Member): Promise<string> {
    if (await this.memberRepository.count({ loginID: member.loginID }) !== 0) {
      return 'MEM000';
    } if (await this.memberRepository.count({ email: member.email }) !== 0) {
      return 'MEM001';
    }
    return null;
  }

  public async save(member: Member): Promise<Member> {
    return this.memberRepository.save(member);
  }

  public async delete(member: Member): Promise<void> {
    await this.memberRepository.delete(member);
  }

  public generateToken(member: Member): string {
    // Use snake case to keep the backward compatibility.
    // eslint-disable-next-line @typescript-eslint/camelcase
    const payload: MemberClaim = { member_uuid: member.uuid };
    return jwt.sign(payload, this.hmacSecret, { issuer: 'PoolC/Nagase', expiresIn: '7d' });
  }

  public validateToken(jwtToken: string): string {
    const claims = jwt.verify(jwtToken, this.hmacSecret) as MemberClaim;
    return claims.member_uuid;
  }
}
