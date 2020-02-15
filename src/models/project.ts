import {
  Column, CreateDateColumn, Entity, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';

@Entity('projects')
export default class Project {
  @PrimaryColumn()
  public id: number;

  @Column()
  public description: string;

  @Column()
  public body: string;

  @Column()
  public name: string;

  @Column()
  public genre: string;

  @Column()
  public participants: string;

  @Column()
  public duration: string;

  @Column()
  public thumbnailURL: string;

  @CreateDateColumn()
  public createdAt: Date;

  @UpdateDateColumn()
  public updatedAt: Date;
}
