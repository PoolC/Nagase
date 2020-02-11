import { MigrationInterface, QueryRunner, Table } from 'typeorm';

export class CreatePostStuffs1580998308003 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.createTable(new Table({
      name: 'boards',
      columns: [
        { name: 'id', type: 'integer', isPrimary: true, isGenerated: true },
        { name: 'name', type: 'varchar', length: '40', isUnique: true },
        { name: 'url_path', type: 'varchar', length: '40', isUnique: true },
        { name: 'read_permission', type: 'varchar', length: '10' },
        { name: 'write_permission', type: 'varchar', length: '10' },
        { name: 'created_at', type: 'timestamp with time zone', default: 'now()' },
        { name: 'updated_at', type: 'timestamp with time zone', default: 'now()' },
      ],
    }));

    await queryRunner.createTable(new Table({
      name: 'posts',
      columns: [
        { name: 'id', type: 'integer', isPrimary: true, isGenerated: true },
        { name: 'board_id', type: 'integer' },
        { name: 'author_uuid', type: 'varchar', length: '40' },
        { name: 'title', type: 'text' },
        { name: 'body', type: 'text' },
        { name: 'vote_id', type: 'integer', isNullable: true },
        { name: 'created_at', type: 'timestamp with time zone', default: 'now()' },
        { name: 'updated_at', type: 'timestamp with time zone', default: 'now()' },
      ],
      indices: [
        { columnNames: ['board_id'] },
      ]
    }));

    await queryRunner.createTable(new Table({
      name: 'comments',
      columns: [
        { name: 'id', type: 'integer', isPrimary: true, isGenerated: true },
        { name: 'post_id', type: 'integer' },
        { name: 'author_uuid', type: 'varchar', length: '40' },
        { name: 'body', type: 'text' },
        { name: 'created_at', type: 'timestamp with time zone', default: 'now()' },
        { name: 'updated_at', type: 'timestamp with time zone', default: 'now()' },
      ],
      indices: [
        { columnNames: ['post_id'] },
      ]
    }));
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropTable('comments');
    await queryRunner.dropTable('posts');
    await queryRunner.dropTable('boards');
  }
}
