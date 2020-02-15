import { expect } from 'chai';
import faker from 'faker';
import sinon from 'sinon';
import { getRepository } from 'typeorm';
import { Factory } from 'rosie';

import { adminCtx, memberCtx } from '../../fixtures/context';
import Board, { BoardPermission } from '../../../src/models/board';
import Comment from '../../../src/models/comment';
import Member from '../../../src/models/member';
import Post from '../../../src/models/post';
import CommentMutation from '../../../src/graphql/mutation/comment';
import CommentService from '../../../src/services/comment';
import PostService from '../../../src/services/post';

describe('CommentMutation', () => {
  const service = new CommentService(getRepository(Comment));
  const boardService = new PostService(getRepository(Post));
  const mutation = new CommentMutation(service, boardService);

  describe('#createComment', () => {
    const dummyBody = faker.lorem.lines(1);
    const post = Factory.build<Post>('Post');
    const board = new Board();

    let findPostSpy: sinon.SinonStub;
    let saveCommentSpy: sinon.SinonStub;
    let member: Member;
    beforeEach(() => {
      post.board = board;
      member = Factory.build<Member>('Member', { isActivated: true });
      findPostSpy = sinon.stub(PostService.prototype, 'findById').resolves(post);
      saveCommentSpy = sinon.stub(CommentService.prototype, 'save').resolvesArg(0);
    });
    afterEach(() => [findPostSpy, saveCommentSpy].forEach((each) => each.restore()));

    context('permission sufficient', () => {
      before(() => { board.readPermission = BoardPermission.MEMBER; });
      it('success', async () => {
        const result = await mutation.createComment({ postID: post.id, body: dummyBody }, memberCtx(member));
        expect(result.author.uuid).to.equal(member.uuid);
        expect(result.post.id).to.equal(post.id);
        expect(result.body).to.equal(dummyBody);
        expect(findPostSpy.calledOnceWith(post.id)).to.be.true;
        expect(saveCommentSpy.calledOnceWith(result)).to.be.true;
      });
    });

    context('permission insufficient', () => {
      before(() => { board.readPermission = BoardPermission.ADMIN; });
      it('failure', () => expect(mutation.createComment({ postID: post.id, body: dummyBody }, memberCtx())).to.be.rejectedWith(Error, /ERR400/));
    });
  });

  describe('#deleteComment', () => {
    let findCommentSpy: sinon.SinonStub;
    let deleteCommentSpy: sinon.SinonStub;
    let member: Member;
    let comment: Comment;
    beforeEach(() => {
      findCommentSpy = sinon.stub(CommentService.prototype, 'findById').resolves(comment);
      deleteCommentSpy = sinon.stub(CommentService.prototype, 'delete').resolvesArg(0);
    });
    afterEach(() => [findCommentSpy, deleteCommentSpy].forEach((each) => each.restore()));

    context('member and comment are valid', () => {
      before(() => {
        member = Factory.build<Member>('Member', { isActivated: true });
        comment = Factory.build<Comment>('Comment', { author: member });
      });

      it('success if permission sufficient', async () => {
        const result = await mutation.deleteComment({ commentID: comment.id }, memberCtx(member));
        expect(result).not.to.be.null;
        expect(deleteCommentSpy.calledOnceWith(result)).to.be.true;
      });

      it('success for admins', async () => {
        const result = await mutation.deleteComment({ commentID: comment.id }, adminCtx());
        expect(result).not.to.be.null;
        expect(deleteCommentSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail if permission insufficient', () => {
        expect(mutation.deleteComment({ commentID: comment.id }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(deleteCommentSpy.notCalled).to.be.true;
      });
    });

    context('comment not exists', () => {
      before(() => { comment = null });
      it('fail', () => {
        expect(mutation.deleteComment({ commentID: 0 }, memberCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(deleteCommentSpy.notCalled).to.be.true;
      });
    });
  });
});
