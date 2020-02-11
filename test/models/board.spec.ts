import { expect } from 'chai';
import { Factory } from 'rosie';

import Board, { BoardPermission } from '../../src/models/board';
import Member from '../../src/models/member';

describe('Board', () => {
  describe('.permissionFromString', () => {
    it('ADMIN', () => { expect(Board.permissionFromString('ADMIN')).to.equal(BoardPermission.ADMIN) });
    it('MEMBER', () => { expect(Board.permissionFromString('MEMBER')).to.equal(BoardPermission.MEMBER) });
    it('PUBLIC', () => { expect(Board.permissionFromString('PUBLIC')).to.equal(BoardPermission.PUBLIC) });
    it('Invalid text', () => {
      expect(() => { Board.permissionFromString('INVALID') }).to.throw(Error, /Invalid/);
    });
  });

  describe('#readPermittedTo, #writePermittedTo', () => {
    const board = new Board();
    let adminMember: Member;
    let normalMember: Member;
    let deactivatedMember: Member;
    before(() => {
      adminMember = Factory.build<Member>('Member', { isAdmin: true, isActivated: true });
      normalMember = Factory.build<Member>('Member', { isAdmin: false, isActivated: true });
      deactivatedMember = Factory.build<Member>('Member', { isAdmin: false, isActivated: false });
    });

    context('read ADMIN', () => {
      before(() => { board.readPermission = BoardPermission.ADMIN });
      it('admin member', () => { expect(board.readPermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.readPermittedTo(normalMember)).to.be.false; });
      it('deactivated member', () => { expect(board.readPermittedTo(deactivatedMember)).to.be.false; });
      it('non member', () => { expect(board.readPermittedTo(null)).to.be.false; });
    });

    context('read MEMBER', () => {
      before(() => { board.readPermission = BoardPermission.MEMBER });
      it('admin member', () => { expect(board.readPermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.readPermittedTo(normalMember)).to.be.true; });
      it('deactivated member', () => { expect(board.readPermittedTo(deactivatedMember)).to.be.false; });
      it('non member', () => { expect(board.readPermittedTo(null)).to.be.false; });
    });

    context('read PUBLIC', () => {
      before(() => { board.readPermission = BoardPermission.PUBLIC });
      it('admin member', () => { expect(board.readPermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.readPermittedTo(normalMember)).to.be.true; });
      it('deactivated member', () => { expect(board.readPermittedTo(deactivatedMember)).to.be.true; });
      it('non member', () => { expect(board.readPermittedTo(null)).to.be.true; });
    });

    context('write ADMIN', () => {
      before(() => { board.writePermission = BoardPermission.ADMIN });
      it('admin member', () => { expect(board.writePermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.writePermittedTo(normalMember)).to.be.false; });
      it('deactivated member', () => { expect(board.writePermittedTo(deactivatedMember)).to.be.false; });
      it('non member', () => { expect(board.writePermittedTo(null)).to.be.false; });
    });

    context('write MEMBER', () => {
      before(() => { board.writePermission = BoardPermission.MEMBER });
      it('admin member', () => { expect(board.writePermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.writePermittedTo(normalMember)).to.be.true; });
      it('deactivated member', () => { expect(board.writePermittedTo(deactivatedMember)).to.be.false; });
      it('non member', () => { expect(board.writePermittedTo(null)).to.be.false; });
    });

    context('write PUBLIC', () => {
      before(() => { board.writePermission = BoardPermission.PUBLIC });
      it('admin member', () => { expect(board.writePermittedTo(adminMember)).to.be.true; });
      it('normal member', () => { expect(board.writePermittedTo(normalMember)).to.be.true; });
      it('deactivated member', () => { expect(board.writePermittedTo(deactivatedMember)).to.be.true; });
      it('non member', () => { expect(board.writePermittedTo(null)).to.be.true; });
    });
  });
});
