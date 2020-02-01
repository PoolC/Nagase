import argon2 from 'argon2';
import crypto from 'crypto';
import {
  BeforeInsert, Column, CreateDateColumn, Entity, PrimaryColumn, UpdateDateColumn,
} from 'typeorm';
import uuidv4 from 'uuid/v4';

@Entity('members')
export default class Member {
  @PrimaryColumn()
  public uuid: string;

  @Column()
  public loginId: string;

  @Column({ name: 'password_hash', type: 'bytea' })
  private passwordHash: Buffer;

  @Column({ name: 'password_salt', type: 'bytea' })
  private passwordSalt: Buffer;

  @Column()
  public email: string;

  @Column()
  public phoneNumber: string;

  @Column()
  public name: string;

  @Column()
  public department: string;

  @Column()
  public studentId: string;

  @Column()
  public isActivated: boolean;

  @Column()
  public isAdmin: boolean;

  @Column()
  public passwordResetToken: string;

  @Column()
  public passwordResetTokenValidUntil: Date;

  @CreateDateColumn()
  public created_at: Date;

  @UpdateDateColumn()
  public updated_at: Date;

  @BeforeInsert()
  beforeInsert(): void {
    if (!this.uuid) {
      this.uuid = uuidv4();
    }
  }

  public async setPassword(value: string): Promise<void> {
    if (!this.passwordSalt) {
      this.passwordSalt = crypto.randomBytes(32);
    }

    const hashStr = await argon2.hash(value, {
      salt: this.passwordSalt,
      timeCost: 1,
      memoryCost: 8 * 1024,
      parallelism: 4,
      hashLength: 32,
    });
    this.passwordHash = Buffer.from(hashStr);
  }

  public async validatePassword(value: string): Promise<boolean> {
    return argon2.verify(this.passwordHash.toString(), value, { salt: this.passwordSalt });
  }
}
