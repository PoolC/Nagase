import { expect } from 'chai';
import faker from 'faker';
import sinon from 'sinon';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';

import { adminCtx } from '../../fixtures/context';
import Board, { BoardPermission } from '../../../src/models/board';
import BoardMutation from '../../../src/graphql/mutation/board';
import BoardService from '../../../src/services/board';

describe('BoardMutation', () => {
  const service = new BoardService(getRepository(Board));
  const mutation = new BoardMutation(service);

  describe('#createBoard', () => {
    let spy: sinon.SinonStub;
    before(() => { spy = sinon.stub(BoardService.prototype, 'save').resolvesArg(0); });
    after(() => { spy.restore(); });

    it('success', async () => {
      const input = {
        name: faker.company.companyName(), urlPath: faker.internet.domainWord(),
        readPermission: 'MEMBER', writePermission: 'ADMIN',
      };

      const result = await mutation.createBoard({ BoardInput: input }, adminCtx());
      expect(result.name).to.equal(input.name);
      expect(result.urlPath).to.equal(input.urlPath);
      expect(result.readPermission).to.equal(BoardPermission.MEMBER);
      expect(result.writePermission).to.equal(BoardPermission.ADMIN);
    });
  });

  describe('#updateBoard', () => {
    let findBoardSpy: sinon.SinonStub;
    let saveSpy: sinon.SinonStub;
    before(() => {
      findBoardSpy = sinon.stub(BoardService.prototype, 'findById').resolves(Factory.build<Board>('Board'));
      saveSpy = sinon.stub(BoardService.prototype, 'save').resolvesArg(0);
    });
    after(() => [findBoardSpy, saveSpy].forEach((each) => { each.restore(); }));

    it('values changes', async () => {
      const input = {
        name: faker.company.companyName(), urlPath: faker.internet.domainWord(),
        readPermission: 'MEMBER', writePermission: 'ADMIN',
      };

      const result = await mutation.updateBoard({ boardID: faker.random.number(), BoardInput: input }, adminCtx());
      expect(result.name).to.equal(input.name);
      expect(result.urlPath).to.equal(input.urlPath);
      expect(result.readPermission).to.equal(BoardPermission.MEMBER);
      expect(result.writePermission).to.equal(BoardPermission.ADMIN);
    });
  });

  describe('#deleteBoard', () => {
    let findBoardSpy: sinon.SinonStub;
    let deleteBoardSpy: sinon.SinonStub;
    before(() => {
      findBoardSpy = sinon.stub(BoardService.prototype, 'findById').resolves(Factory.build<Board>('Board'));
      deleteBoardSpy = sinon.stub(BoardService.prototype, 'delete').resolves();
    });
    after(() => [findBoardSpy, deleteBoardSpy].forEach((each) => { each.restore(); }));

    it('success', async () => {
      expect(await mutation.deleteBoard({ boardID: faker.random.number() }, adminCtx())).not.to.be.null;
      expect(deleteBoardSpy.calledOnce).to.be.true;
    });
  });
});
