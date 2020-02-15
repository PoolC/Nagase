import {
  Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';
import Member from './member';
import Post from './post';

@Entity('comments')
export default class Comment {
  @PrimaryColumn()
  public id: number;

  @ManyToOne((_) => Post, (post) => post.comments)
  @JoinColumn({ name: 'post_id' })
  public post: Post;

  @ManyToOne((_) => Member)
  @JoinColumn({ name: 'member_uuid' })
  public author: Member;

  @Column()
  public body: string;

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;
}
