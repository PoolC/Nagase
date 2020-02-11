import {
  Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';

import Board from './board';
import Member from './member';

@Entity('posts')
export default class Post {
  @PrimaryColumn()
  public id: number;

  @ManyToOne((_) => Board)
  @JoinColumn({ name: 'board_id' })
  public board: Board;

  @ManyToOne((_) => Member)
  @JoinColumn({ name: 'author_uuid' })
  public author: Member;

  @Column()
  public title: string;

  @Column()
  public body: string;

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;
}
