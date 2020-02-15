import faker from 'faker';
import { Factory } from 'rosie';
import uuidv4 from 'uuid/v4';

import Board, { BoardPermission } from '../src/models/board';
import Comment from '../src/models/comment';
import Member from '../src/models/member';
import Post from '../src/models/post';
import Project from "../src/models/project";

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
    .attr('comments', (): Comment[] => [])  // prevent circular dependency
    .attr('title', () => faker.lorem.lines(1))
    .attr('body', () => faker.lorem.lines(10))
    .attr('createdAt', () => new Date())
    .attr('updatedAt', () => new Date());

  Factory.define<Comment>('Comment')
    .sequence('id')
    .attr('post', () => Factory.build<Post>('Post'))
    .attr('author', () => Factory.build<Member>('Member'))
    .attr('body', () => faker.lorem.lines(3))
    .attr('createdAt', () => new Date())
    .attr('updatedAt', () => new Date());

  Factory.define<Project>('Project')
    .sequence('id')
    .attr('description', () => faker.lorem.lines(3))
    .attr('body', () => faker.lorem.lines(10))
    .attr('name', () => faker.random.words(3))
    .attr('genre', () => faker.random.words(3))
    .attr('participants', () => faker.random.words(3))
    .attr('duration', () => faker.random.words(3))
    .attr('thumbnailURL', () => faker.random.words(3))
    .attr('createdAt', () => new Date())
    .attr('updatedAt', () => new Date());
}
