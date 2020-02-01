import fs from 'fs';
import path from 'path';
import { buildSchema } from 'graphql';

const schema = fs.readFileSync(path.join(__dirname, 'schema.graphql'));

export default buildSchema(schema.toString('utf-8'));
