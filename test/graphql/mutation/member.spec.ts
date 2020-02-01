import { expect } from 'chai';
import faker from 'faker';
import sinon from 'sinon';
import { getRepository } from 'typeorm';

import MemberMutation from '../../../src/graphql/mutation/member';
import Member from '../../../src/models/member';
import MemberService from '../../../src/services/member';

describe('MemberMutation', () => {
  const service = new MemberService(getRepository(Member));
  const mutation = new MemberMutation(service);

  describe('#createMember', () => {
    const spy = sinon.stub(MemberService.prototype, 'save').resolvesArg(0);
    const dummyInput = {
      loginID: faker.random.word(),
      password: faker.random.word(),
      email: faker.internet.email(),
      name: faker.name.firstName(),
      phoneNumber: faker.phone.phoneNumber(),
      department: faker.random.word(),
      studentID: faker.random.word(),
    };

    after(() => spy.restore());

    context('no duplicated entry', () => {
      it('success', async () => {
        const result = await mutation.createMember({ MemberInput: dummyInput});
        expect(spy.calledOnce).to.be.true;
        expect(result.loginId).to.equal(dummyInput.loginID);
        expect(await result.validatePassword(dummyInput.password)).to.be.true;
      });
    });

    context('duplicated login id', () => {
      let loginIdSpy: sinon.SinonStub;
      before(() => { loginIdSpy = sinon.stub(MemberService.prototype, 'checkDuplication').resolves('MEM000'); });
      after(() => { loginIdSpy.restore() });
      it('fail', () => {
        return expect(mutation.createMember({ MemberInput: dummyInput })).to.be.rejectedWith(Error, /MEM000/);
      });
    });

    context('duplicated email', () => {
      let loginIdSpy: sinon.SinonStub;
      before(() => { loginIdSpy = sinon.stub(MemberService.prototype, 'checkDuplication').resolves('MEM001'); });
      after(() => { loginIdSpy.restore() });
      it('fail', () => {
        return expect(mutation.createMember({ MemberInput: dummyInput })).to.be.rejectedWith(Error, /MEM001/);
      });
    });
  });

  describe('#createAccessToken', () => {
    const dummyInput = { loginID: faker.random.word(), password: faker.random.word() };

    let member: Member;
    let memberSpy: sinon.SinonStub;
    beforeEach(() => { memberSpy = sinon.stub(MemberService.prototype, 'findByLoginId').resolves(member) });
    afterEach(() => { memberSpy.restore() });

    context('member not exists', () => {
      before(() => { member = null; });
      it('fail', async () => {
        return expect(mutation.createAccessToken({ LoginInput: dummyInput })).to.be.rejectedWith(Error, /TKN000/);
      });
    });

    context('member exists', () => {
      before(() => { member = new Member(); });

      it('fail if password not match', async () => {
        await member.setPassword(faker.random.words(2));
        return expect(mutation.createAccessToken({ LoginInput: dummyInput })).to.be.rejectedWith(Error, /TKN000/);
      });

      it('fail if not activated', async () => {
        await member.setPassword(dummyInput.password);
        member.isActivated = false;
        return expect(mutation.createAccessToken({ LoginInput: dummyInput })).to.be.rejectedWith(Error, /TKN002/);
      });

      it('success if activated', async () => {
        await member.setPassword(dummyInput.password);
        member.isActivated = true;
        return expect(await mutation.createAccessToken({ LoginInput: dummyInput })).not.to.be.empty;
      });
    });
  });
});
