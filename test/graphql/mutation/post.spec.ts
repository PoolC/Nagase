import { expect } from 'chai';
import faker from 'faker';
import sinon from 'sinon';
import { getRepository } from 'typeorm';
import { Factory } from 'rosie';

import {adminCtx, memberCtx} from '../../fixtures/context';
import Board, {BoardPermission} from '../../../src/models/board';
import Member from '../../../src/models/member';
import Post from '../../../src/models/post';
import PostMutation from '../../../src/graphql/mutation/post';
import BoardService from '../../../src/services/board';
import PostService from '../../../src/services/post';

describe('PostMutation', () => {
  const service = new PostService(getRepository(Post));
  const boardService = new BoardService(getRepository(Board));
  const mutation = new PostMutation(service, boardService);

  describe('#createPost', () => {
    const dummyInput = { title: faker.lorem.lines(1), body: faker.lorem.lines(10) };
    const board = new Board();

    let findBoardSpy: sinon.SinonStub;
    let savePostSpy: sinon.SinonStub;
    let member: Member;
    beforeEach(() => {
      findBoardSpy = sinon.stub(BoardService.prototype, 'findById').resolves(board);
      savePostSpy = sinon.stub(PostService.prototype, 'save').resolvesArg(0);
      member = Factory.build<Member>('Member', { isActivated: true });
    });
    afterEach(() => [findBoardSpy, savePostSpy].forEach((each) => each.restore()));

    context('permission sufficient', () => {
      before(() => { board.writePermission = BoardPermission.MEMBER; });
      it('success', async () => {
        const result = await mutation.createPost({ boardID: board.id, PostInput: dummyInput }, memberCtx(member));
        expect(result.title).to.equal(dummyInput.title);
        expect(result.body).to.equal(dummyInput.body);
        expect(result.board.id).to.equal(board.id);
        expect(result.author.uuid).to.equal(member.uuid);
        expect(findBoardSpy.calledOnceWith(board.id)).to.be.true;
        expect(savePostSpy.calledOnceWith(result)).to.be.true;
      });
    });

    context('permission insufficient', () => {
      before(() => { board.writePermission = BoardPermission.ADMIN; });
      it('failure', () => expect(mutation.createPost({ boardID: board.id, PostInput: dummyInput }, memberCtx())).to.be.rejectedWith(Error, /ERR400/));
    });
  });

  describe('#updatePost', () => {
    const dummyInput = { title: faker.lorem.lines(1), body: faker.lorem.lines(10) };

    let findPostSpy: sinon.SinonStub;
    let savePostSpy: sinon.SinonStub;
    let member: Member;
    let post: Post;
    beforeEach(() => {
      findPostSpy = sinon.stub(PostService.prototype, 'findById').resolves(post);
      savePostSpy = sinon.stub(PostService.prototype, 'save').resolvesArg(0);
    });
    afterEach(() => [findPostSpy, savePostSpy].forEach((each) => each.restore()));

    context('member and post are valid', () => {
      before(() => {
        member = Factory.build<Member>('Member', { isActivated: true });
        post = Factory.build<Post>('Post', { author: member });
      });

      it('success if permission sufficient', async () => {
        const result = await mutation.updatePost({ postID: post.id, PostInput: dummyInput }, memberCtx(member));
        expect(result.title).to.equal(dummyInput.title);
        expect(result.body).to.equal(dummyInput.body);
        expect(savePostSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail if permission insufficient', () => {
        expect(mutation.updatePost({ postID: post.id, PostInput: dummyInput }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(savePostSpy.notCalled).to.be.true;
      });
    });

    context('post not exists', () => {
      before(() => { post = null });
      it('fail', () => {
        expect(mutation.updatePost({ postID: 0, PostInput: dummyInput }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(savePostSpy.notCalled).to.be.true;
      });
    });
  });

  describe('#deletePost', () => {
    let findPostSpy: sinon.SinonStub;
    let deletePostSpy: sinon.SinonStub;
    let member: Member;
    let post: Post;
    beforeEach(() => {
      findPostSpy = sinon.stub(PostService.prototype, 'findById').resolves(post);
      deletePostSpy = sinon.stub(PostService.prototype, 'delete').resolvesArg(0);
    });
    afterEach(() => [findPostSpy, deletePostSpy].forEach((each) => each.restore()));

    context('member and post are valid', () => {
      before(() => {
        member = Factory.build<Member>('Member', { isActivated: true });
        post = Factory.build<Post>('Post', { author: member });
      });

      it('success if permission sufficient', async () => {
        const result = await mutation.deletePost({ postID: post.id }, memberCtx(member));
        expect(result).not.to.be.null;
        expect(deletePostSpy.calledOnceWith(result)).to.be.true;
      });

      it('success for admins', async () => {
        const result = await mutation.deletePost({ postID: post.id }, adminCtx());
        expect(result).not.to.be.null;
        expect(deletePostSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail if permission insufficient', () => {
        expect(mutation.deletePost({ postID: post.id }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(deletePostSpy.notCalled).to.be.true;
      });
    });

    context('post not exists', () => {
      before(() => { post = null });
      it('fail', () => {
        expect(mutation.deletePost({ postID: 0 }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(deletePostSpy.notCalled).to.be.true;
      });
    });
  });
});
