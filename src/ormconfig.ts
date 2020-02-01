import { ConnectionOptions } from 'typeorm';
import { SnakeNamingStrategy } from 'typeorm-naming-strategies';

export = {
  type: 'postgres',

  host: process.env.DB_HOST || '127.0.0.1',
  port: process.env.DB_PORT ? parseInt(process.env.DB_PORT, 10) : 5432,
  username: process.env.DB_USERNAME || 'postgres',
  password: process.env.DB_PASSWORD || 'rootpass',
  database: `${process.env.DB_NAME || 'poolc'}${process.env.NODE_ENV === 'test' ? '_test' : ''}`,

  entities: [`${__dirname}/models/**/*{.ts,.js}`],
  migrations: [`${__dirname}/migrations/**/*{.ts,.js}`],

  namingStrategy: new SnakeNamingStrategy(),
  synchronize: false,
  migrationsRun: true,
} as ConnectionOptions;
