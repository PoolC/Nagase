import { expect } from 'chai';
import { Factory } from 'rosie';
import timekeeper from 'timekeeper';
import { getRepository } from 'typeorm';

import Member from '../../src/models/member';
import MemberService from '../../src/services/member';

describe('MemberService', () => {
  const baseTime = new Date(1540309179000); // 2018-10-23T15:39:39.000Z
  const service = new MemberService(getRepository(Member));

  const sampleUuid = '00000000-0000-0000-0000-000000000000';
  const sampleToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfdXVpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlhdCI6MTU0MDMwOTE3OSwiZXhwIjoxNTQwOTEzOTc5LCJpc3MiOiJQb29sQy9OYWdhc2UifQ.J-H_bZMMqB-xnTmjJk7x2T-ch2UgGTMWoH8DAJ_0JQM';

  beforeEach(() => {
    timekeeper.freeze(baseTime);
  });

  afterEach(() => {
    timekeeper.reset();
  });

  describe('#generateToken', () => {
    it('success', async () => {
      const member = Factory.build<Member>('Member', { uuid: sampleUuid });
      expect(service.generateToken(member)).to.equal(sampleToken);
    });
  });

  describe('#validateToken', () => {
    it('success', () => {
      expect(service.validateToken(sampleToken)).to.equal(sampleUuid);
    });

    it('fail if token expired', () => {
      timekeeper.travel(new Date(baseTime.getTime() + (10 * 24 * 60 * 60 * 1000)));
      expect(() => { service.validateToken(sampleToken) }).to.throw(/expired/);
    });

    it('fail if token is invalid', () => {
      expect(() => { service.validateToken(`${sampleToken}_invalid`) }).to.throw(/invalid/);
    });
  });
});
