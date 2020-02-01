import chai from 'chai';
import chaiAsPromised from 'chai-as-promised';
import mochaPrepare from 'mocha-prepare';
import { createConnection } from 'typeorm';

import { defineFactories } from './factories';
import ormconfig from '../src/ormconfig';

mochaPrepare(
  async (done) => {
    await createConnection(ormconfig);
    defineFactories();

    chai.use(chaiAsPromised);

    done();
  },
  (done) => { done() },
);
