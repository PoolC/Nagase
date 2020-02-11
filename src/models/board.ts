import {
  Column, CreateDateColumn, Entity, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';

import Member from './member';

export enum BoardPermission {
  ADMIN = 'ADMIN',
  MEMBER = 'MEMBER',
  PUBLIC = 'PUBLIC',
}

@Entity('boards')
export default class Board {
  @PrimaryColumn()
  public id: number;

  @Column()
  public name: string;

  @Column()
  public urlPath: string;

  @Column({ name: 'read_permission' })
  private readPermissionRecord: string;

  public set readPermission(value: BoardPermission) {
    this.readPermissionRecord = value;
  }

  public get readPermission(): BoardPermission {
    return Board.permissionFromString(this.readPermissionRecord);
  }

  @Column({ name: 'write_permission' })
  private writePermissionRecord: string;

  public set writePermission(value: BoardPermission) {
    this.writePermissionRecord = value;
  }

  public get writePermission(): BoardPermission {
    return Board.permissionFromString(this.writePermissionRecord);
  }

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;

  public readPermittedTo(member: Member): boolean {
    switch (this.readPermission) {
      case BoardPermission.ADMIN: return member?.isAdmin || false;
      case BoardPermission.MEMBER: return member?.isActivated || false;
      case BoardPermission.PUBLIC: return true;
      default: return false;
    }
  }

  public writePermittedTo(member: Member): boolean {
    switch (this.writePermission) {
      case BoardPermission.ADMIN: return member?.isAdmin || false;
      case BoardPermission.MEMBER: return member?.isActivated || false;
      case BoardPermission.PUBLIC: return true;
      default: return false;
    }
  }

  public static permissionFromString(value: string): BoardPermission {
    switch (value) {
      case BoardPermission.ADMIN: return BoardPermission.ADMIN;
      case BoardPermission.MEMBER: return BoardPermission.MEMBER;
      case BoardPermission.PUBLIC: return BoardPermission.PUBLIC;
      default: throw new Error(`Invalid permission: ${value}`);
    }
  }
}
