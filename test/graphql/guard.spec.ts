import { expect } from 'chai';
import { Context } from 'koa';
import { Factory } from 'rosie';

import { Permission, PermissionGuard } from '../../src/graphql/guard';

class GuardTest {
  @PermissionGuard()
  static member(_1: object, _2: Partial<Context>) {}

  @PermissionGuard(Permission.Admin)
  static admin(_1: object, _2: Partial<Context>) {}
}

describe('PermissionGuard', () => {
  context('Permission.Member', () => {
    it('success for member', () => expect(() => { GuardTest.member({}, { state: { member: Factory.build('Member') } }) }).not.to.throw);
    it('success for admin', () => expect(() => { GuardTest.member({}, { state: { member: Factory.build('Member', { isAdmin: true }) } }) }).not.to.throw);
  });

  context('Permission.Admin', () => {
    it('fail for member', () => expect(() => { GuardTest.admin({}, { state: { member: Factory.build('Member') } }) }).to.throw(Error, /ERR401/));
    it('success for admin', () => expect(() => { GuardTest.admin({}, { state: { member: Factory.build('Member', { isAdmin: true }) } }) }).not.to.throw);
  });
});
