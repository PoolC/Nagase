import { expect } from 'chai';
import sinon from 'sinon';
import { getRepository, Repository } from 'typeorm';

import { findPageItems } from '../../src/lib/pagination';
import Board from '../../src/models/board';

describe('#findPageItems', () => {
  const repository = getRepository(Board);

  let findSpy: sinon.SinonStub;
  let countSpy: sinon.SinonStub;
  beforeEach(() => { findSpy = sinon.stub(Repository.prototype, 'find').resolves([1, 2, 3, 4, 5]) });
  afterEach(() => [findSpy, countSpy].forEach((each) => each?.restore()));

  context('find', () => {
    [ // [count(==take), page, skip]
      [10, 1, 0], [10, 2, 10], [15, 3, 30],
    ].forEach((each) => {
      it(`[${each.join(',')}]`, async () => {
        await findPageItems(repository, { page: each[1], count: each[0] }, {}, {}, []);
        expect(findSpy.calledOnceWith({ skip: each[2], take: each[0], where: {}, order:{}, relations: [] })).to.be.true;
      });
    });
  });

  context('totalPage', () => {
    [
      [0, 0], [1, 1], [9, 1], [10, 1], [11, 2],
    ].forEach((each) => {
      it(`totalCount == ${each[0]}, totalPage == ${each[1]}`, async () => {
        countSpy = sinon.stub(Repository.prototype, 'count').resolves(each[0]);
        expect((await findPageItems(repository, { page: 1, count: 10 }, {}, {}, [])).pageInfo.totalPage).to.equal(each[1]);
      });
    });
  });
});
