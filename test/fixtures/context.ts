import { Factory } from 'rosie';

import Member from '../../src/models/member';

export function adminCtx(): object {
  const member = Factory.build('Member', { isActivated: true, isAdmin: true });
  return { state: { member } };
}

export function memberCtx(member?: Member): object {
  const ctxMember = member ? member : Factory.build('Member', { isActivated: true, isAdmin: false });
  return { state: { member: ctxMember } };
}
