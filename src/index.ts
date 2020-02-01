import 'reflect-metadata';
import { Container } from 'typedi';

import ApiApplication from './applications/api';
import initialize from './initialize';

(async (): Promise<void> => {
  await initialize();

  const apiApp = Container.get(ApiApplication);
  apiApp.start();
})();
