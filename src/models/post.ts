import {
  AfterLoad, Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';

import Board from './board';
import Comment from './comment';
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

  @OneToMany((_) => Comment, (comment) => comment.post)
  public comments: Comment[];

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;

  @AfterLoad()
  private afterLoad() {
    if (this.comments?.length > 0) {
      this.comments = this.comments.sort(((a, b) => b.id - a.id));
    }
  }
}
