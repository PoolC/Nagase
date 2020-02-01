import { Container } from 'typedi';
import { createConnection, useContainer } from 'typeorm';

import ormconfig from './ormconfig';

export default async function initialize(): Promise<void> {
  useContainer(Container);
  await createConnection(ormconfig);
}
