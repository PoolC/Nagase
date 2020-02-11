import { Repository } from 'typeorm';

export interface PageOptions {
  page: number;
  count: number;
}

export interface PageInfo {
  currentPage: number;
  totalPage: number;
}

export interface PagedItems<T> {
  items: T[];
  pageInfo: PageInfo;
}

export async function findPageItems<T>(
  repository: Repository<T>, pageOpts: PageOptions, where: object, order: object, relations: string[],
): Promise<PagedItems<T>> {
  const items = await repository.find({
    skip: pageOpts.count * (pageOpts.page - 1),
    take: pageOpts.count,
    where,
    order,
    relations,
  });
  const totalCount = await repository.count();
  const totalPage = Math.floor((totalCount - 1) / pageOpts.count) + 1;

  return { items, pageInfo: { currentPage: pageOpts.page, totalPage } };
}
