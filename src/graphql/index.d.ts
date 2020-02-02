import { Context } from 'koa';

export type QueryParams<T = object> = Parameters<(args?: T, ctx?: Context) => void>;
