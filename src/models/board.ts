import {
  Column, CreateDateColumn, Entity, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';

export enum BoardPermission {
  ADMIN = 'ADMIN',
  MEMBER = 'MEMBER',
  PUBLIC = 'PUBLIC'
}

@Entity('boards')
export default class Board {
  @PrimaryColumn()
  public id: number;

  @Column()
  public name: string;

  @Column()
  public urlPath: string;

  @Column({ type: 'varchar' })
  public readPermission: BoardPermission;

  @Column({ type: 'varchar' })
  public writePermission: BoardPermission;

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;
}
