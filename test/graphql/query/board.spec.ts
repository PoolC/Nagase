import { expect } from 'chai';
import sinon from 'sinon';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';

import BoardService from '../../../src/services/board';
import Board from '../../../src/models/board';
import BoardQuery from '../../../src/graphql/query/board';

describe('BoardQuery', () => {
  const service = new BoardService(getRepository(Board));
  const query = new BoardQuery(service);

  describe('#board', () => {
    let spy: sinon.SinonStub;
    let board: Board;
    beforeEach(() => { spy = sinon.stub(BoardService.prototype, 'findById').resolves(board); });
    afterEach(() => spy.restore());

    context('board exists', () => {
      before(() => { board = Factory.build<Board>('Board'); });
      it('existing id', async () => {
        const result = await query.board({ boardID: board.id });
        expect(result).not.to.be.null;
        expect(result.id).to.equal(board.id);
        expect(result.name).to.equal(board.name);
      });
    });

    context('board not exists', () => {
      before(() => { board = null; });
      it('return null', async () => expect(await query.board({ boardID: -1 })).to.be.null);
    });
  });

  describe('#boards', () => {
    let spy: sinon.SinonStub;
    let boards: Board[];
    beforeEach(() => {
      const sortedBoard = boards.sort((a, b) => a.id - b.id);
      spy = sinon.stub(BoardService.prototype, 'findAll').resolves(sortedBoard);
    });
    afterEach(() => spy.restore());

    context('boards exist', () => {
      before(() => { boards = Factory.buildList('Board', 5); });
      it('success', async () => {
        const result = await query.boards();
        expect(result).have.length(5);
        expect(result[0].id < result[4].id).to.be.true;
      });
    });
  });
});
