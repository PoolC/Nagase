import { Context } from 'koa';
import { QueryParams } from './index';

export enum Permission {
  Member,
  Admin,
}

export function PermissionGuard(permission?: Permission) {
  return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value;

    // eslint-disable-next-line no-param-reassign
    descriptor.value = function (...params: QueryParams) {
      if (params.length < 2 || !params[1]) {
        throw new Error('PermissionGuard requires context.');
      }

      const ctx = params[1] as Context;
      const member = ctx.state?.member;
      if (!member || (permission === Permission.Admin && !member.isAdmin)) {
        throw new Error('ERR401');
      }

      return originalMethod.apply(this, params);
    };

    return descriptor;
  };
}
