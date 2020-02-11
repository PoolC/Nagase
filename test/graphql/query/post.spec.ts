import { expect } from 'chai';
import sinon from 'sinon';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';

import { memberCtx } from '../../fixtures/context';
import PostQuery from '../../../src/graphql/query/post';
import Board, { BoardPermission } from '../../../src/models/board';
import Post from '../../../src/models/post';
import BoardService from '../../../src/services/board';
import PostService from '../../../src/services/post';

describe('PostQuery', () => {
  const query = new PostQuery(
    new PostService(getRepository(Post)),
    new BoardService(getRepository(Board)),
  );

  describe('#post', () => {
    let post: Post;
    let board = new Board();
    let findPostSpy: sinon.SinonStub;
    beforeEach(() => {
      findPostSpy = sinon.stub(PostService.prototype, 'findById').resolves(post);
      if (post) {
        post.board = board;
      }
    });
    afterEach(() => findPostSpy.restore());

    context('post exists', () => {
      before(() => {
        post = Factory.build<Post>('Post');
        board.readPermission = BoardPermission.PUBLIC;
      });
      it('success', async () => {
        const result = await query.post({ postID: post.id }, memberCtx());
        expect(result.title).to.equal(post.title);
        expect(result.body).to.equal(post.body);
      });
    });

    context('post not exists', () => {
      before(() => { post = null; });
      it('raise error', () => expect(query.post({ postID: -1 }, memberCtx())).to.be.rejectedWith(Error, /ERR400/));
    });

    context('insufficient permission', () => {
      before(() => {
        post = Factory.build<Post>('Post');
        board.readPermission = BoardPermission.ADMIN;
      });
      it('raise error', () => expect(query.post({ postID: post.id }, memberCtx())).to.be.rejectedWith(Error, /ERR403/));
    });
  });

  describe('#postPage', () => {
    let board = new Board();
    let findBoardSpy: sinon.SinonStub;
    let findPostsSpy: sinon.SinonStub;

    beforeEach(() => {
      findBoardSpy = sinon.stub(BoardService.prototype, 'findById').resolves(board);
      findPostsSpy = sinon.stub(PostService.prototype, 'findPageByBoardId').resolves({
        items: Factory.buildList('Post', 10), pageInfo: { currentPage: 1, totalPage: 10 },
      });
    });
    afterEach(() => [findBoardSpy, findPostsSpy].forEach((each) => each.restore()));

    context('sufficient permission', () => {
      before(() => { board.readPermission = BoardPermission.MEMBER });
      it('success', async () => expect(await query.postPage({ boardID: board.id }, memberCtx())).not.to.be.null);
    });

    context('insufficient permission', () => {
      before(() => { board.readPermission = BoardPermission.ADMIN; });
      it('raise error', () => expect(query.postPage({ boardID: board.id }, memberCtx())).to.be.rejectedWith(Error, /ERR403/));
    });
  });
});
