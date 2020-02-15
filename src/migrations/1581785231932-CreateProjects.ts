import { MigrationInterface, QueryRunner, Table } from 'typeorm';

export class CreateProjects1581785231932 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.createTable(new Table({
      name: 'projects',
      columns: [
        { name: 'id', type: 'integer', isPrimary: true, isGenerated: true },
        { name: 'description', type: 'text' },
        { name: 'body', type: 'text' },
        { name: 'name', type: 'varchar', length: '255' },
        { name: 'genre', type: 'varchar', length: '255' },
        { name: 'participants', type: 'varchar', length: '255', isNullable: true },
        { name: 'duration', type: 'varchar', length: '255', isNullable: true },
        { name: 'thumbnail_url', type: 'varchar', length: '255', isNullable: true },
        { name: 'created_at', type: 'timestamp with time zone', default: 'now()' },
        { name: 'updated_at', type: 'timestamp with time zone', default: 'now()' },
      ],
    }));
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropTable('projects');
  }
}
