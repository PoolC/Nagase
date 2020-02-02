import { expect } from 'chai';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';

import MemberQuery from '../../../src/graphql/query/member';
import Member from '../../../src/models/member';
import MemberService from '../../../src/services/member';

describe('MemberQuery', () => {
  const service = new MemberService(getRepository(Member));
  const mutation = new MemberQuery(service);

  describe('#me', () => {
    it('success', async () => {
      const member = Factory.build<Member>('Member');
      const result = await mutation.me({}, { state: { member } });
      expect(result.uuid).to.equal(member.uuid);
    });

    it('fail', () => {
      expect(function () { mutation.me({}, {}) }).to.throw(Error, /ERR401/);
    });
  });
});
