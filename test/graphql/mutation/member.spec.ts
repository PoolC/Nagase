import { expect } from 'chai';
import faker from 'faker';
import sinon from 'sinon';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';
import uuidv4 from 'uuid/v4';

import { adminCtx, memberCtx } from "../../fixtures/context";
import MemberMutation from '../../../src/graphql/mutation/member';
import Member from '../../../src/models/member';
import MemberService from '../../../src/services/member';

describe('MemberMutation', () => {
  const service = new MemberService(getRepository(Member));
  const mutation = new MemberMutation(service);

  let saveMemberSpy: sinon.SinonStub;
  beforeEach(() => { saveMemberSpy = sinon.stub(MemberService.prototype, 'save').resolvesArg(0); });
  afterEach(() => saveMemberSpy.restore());

  describe('#createMember', () => {
    const dummyInput = {
      loginID: faker.random.word(),
      password: faker.random.word(),
      email: faker.internet.email(),
      name: faker.name.firstName(),
      phoneNumber: faker.phone.phoneNumber(),
      department: faker.random.word(),
      studentID: faker.random.word(),
    };

    context('no duplicated entry', () => {
      it('success', async () => {
        const result = await mutation.createMember({ MemberInput: dummyInput });
        expect(saveMemberSpy.calledOnce).to.be.true;
        expect(result.loginID).to.equal(dummyInput.loginID);
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

  describe('#updateMember', () => {
    let member: Member;
    let findMemberSpy: sinon.SinonStub;

    context('member exists', () => {
      beforeEach(() => {
        member = new Member();
        member.uuid = uuidv4();
        member.department = faker.random.word();
        findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(member);
      });
      afterEach(() => findMemberSpy.restore());

      it('values changes', async () => {
        const input = { uuid: member.uuid, password: faker.random.word(), name: faker.random.word() };
        const result = await mutation.updateMember({ MemberInput: input }, memberCtx(member));
        expect(result.uuid).to.equal(member.uuid);
        expect(result.department).to.equal(member.department);
        expect(result.name).to.equal(input.name);
        expect(await result.validatePassword(input.password)).to.be.true;
        expect(saveMemberSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail for another member', async () => {
        const input = { uuid: member.uuid, name: faker.random.word() };
        await mutation.updateMember({ MemberInput: input }, memberCtx());
        expect(saveMemberSpy.notCalled).to.be.true;
      });
    });

    context('member not exists', () => {
      it('nothing happens', async () => {
        const input = { uuid: uuidv4(), name: faker.random.word() };
        await mutation.updateMember({ MemberInput: input }, memberCtx(member));
        expect(saveMemberSpy.notCalled).to.be.true;
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

  describe('#deleteMember', () => {
    let member: Member;
    let findMemberSpy: sinon.SinonStub;
    let deleteMemberSpy: sinon.SinonStub;
    afterEach(() => [findMemberSpy, deleteMemberSpy].forEach((each) => each.restore()));

    context('member exists', () => {
      beforeEach(() => {
        member = Factory.build('Member');
        findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(member);
        deleteMemberSpy = sinon.stub(MemberService.prototype, 'delete').resolves();
      });
      afterEach(() => [findMemberSpy, deleteMemberSpy].forEach((each) => each.restore()));

      it('success', async () => {
        expect(await mutation.deleteMember({ memberUUID: member.uuid }, adminCtx())).not.to.be.null;
        expect(findMemberSpy.calledOnceWith(member.uuid)).to.be.true;
        expect(deleteMemberSpy.calledOnceWith(member)).to.be.true;
      });
    });

    context('member not exists', () => {
      beforeEach(() => {
        findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(null);
        deleteMemberSpy = sinon.stub(MemberService.prototype, 'delete').resolves();
      });
      afterEach(() => [findMemberSpy, deleteMemberSpy].forEach((each) => each.restore()));

      it('success but nothing happens', async () => {
        expect(await mutation.deleteMember({ memberUUID: uuidv4() }, adminCtx())).to.be.null;
        expect(deleteMemberSpy.notCalled).to.be.true;
      });
    });
  });

  describe('#toggleMemberIsActivated', () => {
    let member: Member;
    let findMemberSpy: sinon.SinonStub;

    context('member exists', () => {
      beforeEach(() => {
        member = Factory.build<Member>('Member', { isActivated: true });
        findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(member);
      });
      afterEach(() => findMemberSpy.restore());

      it('success', async () => {
        const result = await mutation.toggleMemberIsActivated({ memberUUID: member.uuid }, adminCtx());
        expect(result).not.to.be.null;
        expect(result.isActivated).to.be.false;
        expect(findMemberSpy.calledOnceWith(member.uuid)).to.be.true;
        expect(saveMemberSpy.calledOnceWith(member)).to.be.true;
      });
    });

    context('member not exists', () => {
      beforeEach(() => findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(null));
      afterEach(() => findMemberSpy.restore());

      it('success but nothing happens', async () => {
        expect(await mutation.toggleMemberIsActivated({ memberUUID: uuidv4() }, adminCtx())).to.be.null;
        expect(saveMemberSpy.notCalled).to.be.true;
      });
    });
  });

  describe('#toggleMemberIsAdmin', () => {
    let member: Member;
    let findMemberSpy: sinon.SinonStub;
    afterEach(() => { findMemberSpy.restore() });

    context('member exists', () => {
      beforeEach(() => {
        member = Factory.build<Member>('Member', { isAdmin: true });
        findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(member);
      });
      afterEach(() => findMemberSpy.restore());

      it('success', async () => {
        const result = await mutation.toggleMemberIsAdmin({ memberUUID: member.uuid }, adminCtx());
        expect(result).not.to.be.null;
        expect(result.isAdmin).to.be.false;
        expect(findMemberSpy.calledOnceWith(member.uuid)).to.be.true;
        expect(saveMemberSpy.calledOnceWith(member)).to.be.true;
      });
    });

    context('member not exists', () => {
      beforeEach(() => findMemberSpy = sinon.stub(MemberService.prototype, 'findByUuid').resolves(null));
      afterEach(() => findMemberSpy.restore());

      it('success but nothing happens', async () => {
        expect(await mutation.toggleMemberIsAdmin({ memberUUID: uuidv4() }, adminCtx())).to.be.null;
        expect(saveMemberSpy.notCalled).to.be.true;
      });
    });
  });
});
