import { MigrationInterface, QueryRunner, Table } from 'typeorm';

export class CreateMembers1580309580074 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.createTable(new Table({
      name: 'members',
      columns: [
        { name: 'uuid', type: 'varchar', length: '40', isPrimary: true },
        { name: 'login_id', type: 'varchar', length: '40', isUnique: true },
        { name: 'password_hash', type: 'bytea' },
        { name: 'password_salt', type: 'bytea' },
        { name: 'email', type: 'varchar', length: '255', isUnique: true },
        { name: 'phone_number', type: 'varchar', length: '20' },
        { name: 'name', type: 'varchar', length: '40' },
        { name: 'department', type: 'varchar', length: '40' },
        { name: 'student_id', type: 'varchar', length: '40', isUnique: true },
        { name: 'is_activated', type: 'boolean', default: 'false' },
        { name: 'is_admin', type: 'boolean', default: 'false' },
        { name: 'password_reset_token', type: 'varchar', length: '255', isUnique: true, isNullable: true },
        { name: 'password_reset_token_valid_until', type: 'timestamp with time zone', isNullable: true },
        { name: 'created_at', type: 'timestamp with time zone', default: 'now()' },
        { name: 'updated_at', type: 'timestamp with time zone', default: 'now()' },
      ],
    }));
  }

  public async down(queryRunner: QueryRunner): Promise<any> {
    await queryRunner.dropTable('members');
  }
}
