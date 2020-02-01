import { expect } from 'chai';
import { Factory } from 'rosie';

import MemberQuery from '../../../src/graphql/query/member';
import Member from '../../../src/models/member';

describe('MemberQuery', () => {
  const mutation = new MemberQuery();

  describe('#me', () => {
    it('success', async () => {
      const member = Factory.build<Member>('Member');
      const result = await mutation.me({}, { state: { member } });
      expect(result.uuid).to.equal(member.uuid);
    });

    it('fail', () => {
      return expect(mutation.me({}, {})).to.be.rejectedWith(Error, /ERR401/);
    });
  });
});
