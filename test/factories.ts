import { Factory } from 'rosie';
import uuidv4 from 'uuid/v4';

import Member from '../src/models/member';

export function defineFactories() {
  Factory.define<Member>('Member')
    .attr('uuid', () => uuidv4());
}
