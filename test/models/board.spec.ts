import { expect } from 'chai';

import Board, { BoardPermission } from '../../src/models/board';

describe('Board', () => {
  describe('.permissionFromString', () => {
    it('ADMIN', () => { expect(Board.permissionFromString('ADMIN')).to.equal(BoardPermission.ADMIN) });
    it('MEMBER', () => { expect(Board.permissionFromString('MEMBER')).to.equal(BoardPermission.MEMBER) });
    it('PUBLIC', () => { expect(Board.permissionFromString('PUBLIC')).to.equal(BoardPermission.PUBLIC) });
    it('Invalid text', () => {
      expect(() => { Board.permissionFromString('INVALID') }).to.throw(Error, /Invalid/);
    });
  });
});
