import faker from 'faker';
import { Factory } from 'rosie';
import uuidv4 from 'uuid/v4';

import Board, {BoardPermission} from '../src/models/board';
import Member from '../src/models/member';
import Post from '../src/models/post';

export function defineFactories() {
  Factory.define<Member>('Member')
    .attr('uuid', () => uuidv4())
    .attr('isActivated', () => false)
    .attr('isAdmin', () => false);

  Factory.define<Board>('Board')
    .sequence('id')
    .attr('name', () => faker.company.companyName())
    .attr('urlPath', () => faker.internet.domainWord())
    .attr('readPermission', () => BoardPermission.PUBLIC)
    .attr('writePermission', () => BoardPermission.PUBLIC)
    .attr('createdAt', () => new Date())
    .attr('updatedAt', () => new Date());

  Factory.define<Post>('Post')
    .sequence('id')
    .attr('board', () => Factory.build<Board>('Board'))
    .attr('author', () => Factory.build<Member>('Member'))
    .attr('title', () => faker.lorem.lines(1))
    .attr('body', () => faker.lorem.lines(10))
    .attr('createdAt', () => new Date())
    .attr('updatedAt', () => new Date());
}
